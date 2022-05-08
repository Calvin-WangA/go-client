package exchange

import (
	"encoding/xml"
	"log"
)

type nodeXmlHandler struct {
	name string
}

func (nodeXmlHandler nodeXmlHandler) getName() string {
	return nodeXmlHandler.name
}

func (nodeXmlHandler nodeXmlHandler) inboundHandle(context *context) (int, string) {

	var hzbankResponse HzbankResponse
	message := context.parameter["recvMessage"]
	err := xml.Unmarshal([]byte(message), &hzbankResponse)
	if err != nil {
		log.Println("字符串不能解析为对应类：", message)
		log.Println("解析错误信息：", err)
		return -1, err.Error()
	}

	return hzbankResponse.Body.Status.ErrorCode, hzbankResponse.Body.Status.ErrorMsg
}