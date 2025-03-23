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
	conn        net.Conn
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

	vodMessage := value.(ReverseCommandMessage)

	if vodMessage.conn != nil {
		vodMessage.conn.Close()
	}

	return
}

func SendReverseCommand(secret string, rcm ReverseCommandMessage, paramStr string) (err *common.Rtmp2FlvCustomError) {
	_, ok := reverseCommandMessageMap.Load(rcm.MessageId)
	if ok {
		logs.Error("MessageId: %s exists", rcm.MessageId)
		errStr := fmt.Sprintf("MessageId: %s exists", rcm.MessageId)
		err = common.CustomError(errStr)
		return
	}
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
		vodMessage := value.(ReverseCommandMessage)
		vodMessage.conn = conn

		readMessage(clientInfo.Secret, conn, vodMessage)
	} else {
		value, ok := reverseCommandMessageMap.Load(registerInfo.MessageId)
		if !ok {
			logs.Error("messageId: %s channel not found", registerInfo.MessageId)
			return
		}
		vodMessage := value.(ReverseCommandMessage)
		vodMessage.conn = conn

		readRes(clientInfo.Secret, conn, vodMessage)
	}
}

func readMessage(secret string, conn net.Conn, vodMessage ReverseCommandMessage) {
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

func readRes(secret string, conn net.Conn, vodMessage ReverseCommandMessage) {
	defer func() {
		close(vodMessage.MessageChan)
	}()
	readOneMessage(secret, conn, vodMessage)
}

func readOneMessage(secret string, conn net.Conn, vodMessage ReverseCommandMessage) bool {
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

	vodMessage.MessageChan <- &resMessage
	return false
}
