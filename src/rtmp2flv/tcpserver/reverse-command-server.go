package tcpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/tcpserver/reportcamerastatus"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/utils"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/common"
	base_service "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/service/base"
)

type ResMessage struct {
	MessageId string
	Data      *[]byte
}

type ReverseCommandMessage struct {
	ClientCode string
	// "cameraAq" "historyVideoPage" "flvFileMediaInfo" "flvPlay" "flvFetchMoreData" "startPushRtmp" "stopPushRtmp" "getLiveMediaInfo"
	MessageType string
	MessageId   string
	Created     time.Time
	MessageChan chan<- *ResMessage
	// HTTP 请求超时或断开时关闭 Done，用来通知 TCP 读协程退出。
	Done <-chan struct{}
	// 每个命令对应的回包连接，挂在共享对象上，便于清理时关闭。
	conn   net.Conn
	mutex  sync.Mutex
	closed bool
}

var reverseCommandClientMap sync.Map
var reverseCommandMessageMap sync.Map

func ClearReverseCommand(messageId string) (err error) {
	value, ok := reverseCommandMessageMap.LoadAndDelete(messageId)
	if !ok {
		logs.Error("MessageId: %s not exists", messageId)
		err = fmt.Errorf("MessageId: %s not exists", messageId)
		return
	}

	// 删除命令后关闭回包连接，让正在读 TCP 的协程尽快退出。
	vodMessage := value.(*ReverseCommandMessage)
	vodMessage.CloseConn()

	return
}

// SetConn 在 HTTP 侧已经清理命令时返回 false。
func (rcm *ReverseCommandMessage) SetConn(conn net.Conn) bool {
	rcm.mutex.Lock()
	defer rcm.mutex.Unlock()
	if rcm.closed {
		return false
	}
	rcm.conn = conn
	return true
}

// CloseConn 需要幂等，超时、取消和 TCP 退出可能同时发生。
func (rcm *ReverseCommandMessage) CloseConn() {
	rcm.mutex.Lock()
	if rcm.closed {
		rcm.mutex.Unlock()
		return
	}
	rcm.closed = true
	conn := rcm.conn
	rcm.mutex.Unlock()

	if conn != nil {
		err := conn.Close()
		if err != nil {
			logs.Error("close reverse command conn error: %v", err)
		}
	}
}

// Send 在 HTTP 请求已经结束时直接退出，避免阻塞发送。
func (rcm *ReverseCommandMessage) Send(resMessage *ResMessage) bool {
	if rcm.Done == nil {
		rcm.MessageChan <- resMessage
		return true
	}
	select {
	case rcm.MessageChan <- resMessage:
		return true
	case <-rcm.Done:
		return false
	}
}

func SendReverseCommand(secret string, rcm *ReverseCommandMessage, paramStr string) (err *common.Rtmp2FlvCustomError) {
	_, ok := reverseCommandMessageMap.Load(rcm.MessageId)
	if ok {
		logs.Error("MessageId: %s exists", rcm.MessageId)
		errStr := fmt.Sprintf("MessageId: %s exists", rcm.MessageId)
		err = common.CustomError(errStr)
		return
	}
	// 保存原始指针，让 handleConn 和 ClearReverseCommand 共享 conn/closed 状态。
	reverseCommandMessageMap.Store(rcm.MessageId, rcm)

	value, ok := reverseCommandClientMap.Load(rcm.ClientCode)
	if !ok {
		logs.Error("reverse conn not exists")
		errStr := fmt.Sprintf("reverse conn: %s not exists", rcm.ClientCode)
		err = common.CustomError(errStr)
		return
	}
	conn := value.(net.Conn)
	cm := CommandMessage{
		MessageType: rcm.MessageType,
		Param:       paramStr,
		MessageId:   rcm.MessageId,
	}

	_, writeErr := writeCommandMessage(secret, cm, conn)
	if writeErr != nil {
		logs.Error("writeCommandMessage error: %v", writeErr)
		err = common.InternalError(writeErr)
		return
	}

	return
}

func ReverseCommandServer() {
	port, err := config.Int("server.reverse.command.port")
	if err != nil {
		logs.Error("get reverse command port error: %v", err)
		return
	}
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logs.Error("ReverseCommandServer listen error: %v", err)
	}
	defer func() {
		listen.Close()
	}()
	for {
		conn, err := listen.Accept()
		if err != nil {
			logs.Error("ReverseCommandServer Accept error: %v", err)
			break
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	registerInfo := RegisterInfo{}
	defer func() {
		if registerInfo.ConnType != "keepChannel" {
			conn.Close()
		}
	}()
	// read first message
	dataLenBytes := make([]byte, 4)
	_, err := conn.Read(dataLenBytes)
	if err != nil {
		logs.Error("conn read message len error: %v", err)
		return
	}
	dataLen := utils.BigEndianToUint32(dataLenBytes)

	registerMaxLen := uint32(64 * 1024)
	if dataLen > registerMaxLen {
		logs.Error("register message len too long: %d, max len: %d", dataLen, registerMaxLen)
		return
	}
	dataBodyBytes := make([]byte, 0)
	for {
		buffer := make([]byte, dataLen-uint32(len(dataBodyBytes)))
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				logs.Error("conn read message body error: %v", err)
				return
			}
			break
		}

		// 处理读取到的数据，n是实际读取的字节数
		dataBodyBytes = append(dataBodyBytes, buffer[:n]...)
		if uint32(len(dataBodyBytes)) == dataLen {
			break
		}
	}

	err = json.Unmarshal(dataBodyBytes, &registerInfo)
	if err != nil {
		logs.Error("Unmarshal RegisterInfo error: %v", err)
		return
	}
	// validate sign
	condition := common.GetEqualCondition("clientCode", registerInfo.ClientCode)
	clientInfo, err := base_service.ClientInfoFindOneByCondition(condition)
	if err != nil {
		logs.Error("ClientInfoFindOneByCondition error: %v", err)
		return
	}
	planText := fmt.Sprintf("clientCode=%s&dateStr=%s&signSecret=%s", registerInfo.ClientCode, registerInfo.DateStr, clientInfo.SignSecret)
	signStr := utils.Md5(planText)
	if signStr != registerInfo.Sign {
		logs.Error("sign: %s error", registerInfo.Sign)
		return
	}

	// 配置了大于0的时间才做有效期验证
	clientRangSeconds, err := config.Int64("server.client-rang-seconds")
	if err != nil {
		logs.Error("get server client-rang-seconds error: %v. \n ", err)
		return
	}
	if clientRangSeconds > 0 {
		registerDate, err := time.Parse(time.RFC3339, registerInfo.DateStr)
		if err != nil {
			logs.Error("parse register dateStr: %s error: %v", registerInfo.DateStr, err)
			return
		}

		fgExpires := math.Abs(float64(time.Since(registerDate))) > float64(time.Duration(clientRangSeconds)*time.Second)
		if fgExpires {
			logs.Error("dateStr: %s expires", registerInfo.DateStr)
			return
		}
	}

	if registerInfo.ConnType == "keepChannel" {
		value, ok := reverseCommandClientMap.Load(registerInfo.ClientCode)
		if ok {
			oldConn := value.(net.Conn)
			err := oldConn.Close()
			if err != nil {
				logs.Error("close old conn error: %v", err)
			}
			reverseCommandClientMap.Delete(registerInfo.ClientCode)
		}
		logs.Info("keepChannel: %s connect", registerInfo.ClientCode)
		reverseCommandClientMap.Store(registerInfo.ClientCode, conn)
	} else if registerInfo.ConnType == "cameraOnline" {
		reportcamerastatus.OnlineStatus(registerInfo.CameraCode)
	} else if registerInfo.ConnType == "cameraOffline" {
		reportcamerastatus.OfflineStatus(registerInfo.CameraCode)
	} else if registerInfo.ConnType == "flvPlay" {
		value, ok := reverseCommandMessageMap.Load(registerInfo.MessageId)
		if !ok {
			logs.Error("messageId: %s channel not found", registerInfo.MessageId)
			return
		}
		vodMessage := value.(*ReverseCommandMessage)
		// 将客户端回包连接绑定到命令，后续由 HTTP 侧统一清理。
		if !vodMessage.SetConn(conn) {
			logs.Error("messageId: %s is closed", registerInfo.MessageId)
			return
		}

		readMessage(clientInfo.Secret, conn, vodMessage)
	} else {
		value, ok := reverseCommandMessageMap.Load(registerInfo.MessageId)
		if !ok {
			logs.Error("messageId: %s channel not found", registerInfo.MessageId)
			return
		}
		vodMessage := value.(*ReverseCommandMessage)
		// 将一次性回包连接绑定到命令，超时或取消时可以关闭连接。
		if !vodMessage.SetConn(conn) {
			logs.Error("messageId: %s is closed", registerInfo.MessageId)
			return
		}

		readRes(clientInfo.Secret, conn, vodMessage)
	}
}

func readMessage(secret string, conn net.Conn, vodMessage *ReverseCommandMessage) {
	defer func() {
		close(vodMessage.MessageChan)
	}()
	for {
		shouldReturn := readOneMessage(secret, conn, vodMessage)
		if shouldReturn {
			return
		}
	}
}

func readRes(secret string, conn net.Conn, vodMessage *ReverseCommandMessage) {
	defer func() {
		close(vodMessage.MessageChan)
	}()
	readOneMessage(secret, conn, vodMessage)
}

func readOneMessage(secret string, conn net.Conn, vodMessage *ReverseCommandMessage) bool {
	dataLenBytes := make([]byte, 4)
	_, err := conn.Read(dataLenBytes)
	if err != nil {
		logs.Error("conn read message len error: %v", err)
		return true
	}

	dataLen := utils.BigEndianToUint32(dataLenBytes)

	dataBodyBytes := make([]byte, 0)
	for {
		buffer := make([]byte, dataLen-uint32(len(dataBodyBytes)))
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				logs.Error("conn read message body error: %v", err)
				return true
			}
			break
		}
		// 处理读取到的数据，n是实际读取的字节数
		dataBodyBytes = append(dataBodyBytes, buffer[:n]...)
		if uint32(len(dataBodyBytes)) == dataLen {
			break
		}
	}

	plainBytes, err := utils.DecryptAES([]byte(secret), dataBodyBytes)
	if err != nil {
		logs.Error("DecryptAES message body error: %v", err)
		return true
	}

	resMessage := ResMessage{
		MessageId: vodMessage.MessageId,
		Data:      &plainBytes,
	}

	// 如果 HTTP 请求已结束，不要让 TCP 读协程卡在发送结果上。
	return !vodMessage.Send(&resMessage)
}
