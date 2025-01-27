package tcpserver

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/utils"
)

type ReverseCommandMessage struct {
	ClientCode string
	// "cameraAq" "historyVideoPage" "flvFileMediaInfo" "flvPlay" "flvFetchMoreData" "startPushRtmp" "stopPushRtmp"
	MessageType string
	MessageId   string
	Created     time.Time
	MessageChan chan<- []byte
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

func SendReverseCommand(secret string, rcm ReverseCommandMessage, paramStr string) (err error) {
	_, ok := reverseCommandMessageMap.Load(rcm.MessageId)
	if ok {
		logs.Error("MessageId: %s exists", rcm.MessageId)
		err = fmt.Errorf("MessageId: %s exists", rcm.MessageId)
		return
	}
	reverseCommandMessageMap.Store(rcm.MessageId, rcm)

	value, ok := reverseCommandClientMap.Load(rcm.ClientCode)
	if !ok {
		logs.Error("reverse conn not exists")
		err = fmt.Errorf("reverse conn: %s not exists", rcm.ClientCode)
		return
	}
	conn := value.(net.Conn)
	cm := CommandMessage{
		MessageType: rcm.MessageType,
		Param:       paramStr,
		MessageId:   rcm.MessageId,
	}

	_, err = writeCommandMessage(secret, cm, conn)
	if err != nil {
		logs.Error("writeCommandMessage error: %v", err)
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
	// read first message
	dataLenBytes := make([]byte, 4)
	_, err := conn.Read(dataLenBytes)
	if err != nil {
		logs.Error("conn read message len error: %v", err)
		return
	}
	dataLen := utils.BigEndianToUint32(dataLenBytes)
	logs.Info("dataLen: %d", dataLen)
	registerMaxLen := uint32(64 * 1024)
	if dataLen > registerMaxLen {
		logs.Error("register message len too long: %d, max len: %d", dataLen, registerMaxLen)
		return
	}
	dataBodyBytes := make([]byte, dataLen)
	_, err = conn.Read(dataBodyBytes)
	if err != nil {
		logs.Error("conn read message body error: %v", err)
		return
	}

	registerInfo := RegisterInfo{}
	err = json.Unmarshal(dataBodyBytes, &registerInfo)
	if err != nil {
		logs.Error("Unmarshal RegisterInfo error: %v", err)
		return
	}
	// validate sign

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
		reverseCommandClientMap.Store(registerInfo.ClientCode, conn)
	} else if registerInfo.ConnType == "flvPlay" {
		value, ok := reverseCommandMessageMap.Load(registerInfo.MessageId)
		if !ok {
			logs.Error("messageId: %s channel not found", registerInfo.MessageId)
			return
		}
		vodMessage := value.(ReverseCommandMessage)
		vodMessage.conn = conn

		readMessage(conn, vodMessage)
	} else {
		value, ok := reverseCommandMessageMap.Load(registerInfo.MessageId)
		if !ok {
			logs.Error("messageId: %s channel not found", registerInfo.MessageId)
			return
		}
		vodMessage := value.(ReverseCommandMessage)
		vodMessage.conn = conn

		readRes(conn, vodMessage)
	}
}

func readMessage(conn net.Conn, vodMessage ReverseCommandMessage) {
	defer func() {
		conn.Close()
		close(vodMessage.MessageChan)
	}()
	for {
		shouldReturn := readOneMessage(conn, vodMessage)
		if shouldReturn {
			return
		}
	}
}

func readRes(conn net.Conn, vodMessage ReverseCommandMessage) {
	defer func() {
		conn.Close()
		close(vodMessage.MessageChan)
	}()
	readOneMessage(conn, vodMessage)
}

func readOneMessage(conn net.Conn, vodMessage ReverseCommandMessage) bool {
	dataLenBytes := make([]byte, 4)
	_, err := conn.Read(dataLenBytes)
	if err != nil {
		logs.Error("conn read message len error: %v", err)
		return true
	}
	dataLen := utils.BigEndianToUint32(dataLenBytes)

	dataBodyBytes := make([]byte, dataLen)
	_, err = conn.Read(dataBodyBytes)
	if err != nil {
		logs.Error("conn read message body error: %v", err)
		return true
	}

	//TODO get secret
	secret := "A012345678901234"
	resultStr, err := utils.DecryptAES([]byte(secret), string(dataBodyBytes))
	if err != nil {
		logs.Error("DecryptAES message body error: %v", err)
		return true
	}

	vodMessage.MessageChan <- []byte(resultStr)
	return false
}
