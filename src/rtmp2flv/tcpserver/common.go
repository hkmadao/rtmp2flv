package tcpserver

import (
	"encoding/json"
	"io"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/utils"
)

type CommandMessage struct {
	// "cameraAq" "historyVideoPage" "flvFileMediaInfo" "flvPlay" "flvFetchMoreData" "startPushRtmp" "stopPushRtmp"
	MessageType string `json:"messageType"`
	Param       string `json:"param"`
	MessageId   string `json:"messageId"`
}

// when connect to server, first send register packet to server
type RegisterInfo struct {
	ClientCode string `json:"clientCode"`
	DateStr    string `json:"dateStr"`
	Sign       string `json:"sign"`
	// "keepChannel" "cameraAq" "historyVideoPage" "flvFileMediaInfo" "flvPlay" "flvFetchMoreData" "startPushRtmp" "stopPushRtmp"
	ConnType  string `json:"connType"`
	MessageId string `json:"messageId"`
}

func writeCommandMessage(secretStr string, commandMessage CommandMessage, writer io.Writer) (n int, err error) {
	messageBytes, err := json.Marshal(commandMessage)
	if err != nil {
		logs.Error(err)
		return
	}

	// messageLenBytes := utils.Int32ToByteBigEndian(int32(len(messageBytes)))
	// fullMessageBytes := append(messageLenBytes, messageBytes...)
	// n, err = writer.Write(fullMessageBytes)
	// if err != nil {
	// 	logs.Error("register error: %v", err)
	// 	return
	// }
	// return

	encryptMessageBytes, err := utils.EncryptAES([]byte(secretStr), messageBytes)
	if err != nil {
		logs.Error("EncryptAES error: %v", err)
		return
	}

	encryptMessageLen := len(encryptMessageBytes)
	encryptMessageLenBytes := utils.Int32ToByteBigEndian(int32(encryptMessageLen))
	fullMessageBytes := append(encryptMessageLenBytes, encryptMessageBytes...)
	n, err = writer.Write(fullMessageBytes)
	if err != nil {
		logs.Error("register error: %v", err)
		return
	}
	return
}
