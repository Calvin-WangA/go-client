package exchange

import (
	"encoding/xml"
	"log"
)


type nodeXmlHandler struct {
	name string
}

/**
  返回业务处理器名称
*/
func (xmlHandler nodeXmlHandler) getName() string {
	return xmlHandler.name
}

/**
  解析文本xml为结构对象
*/
func (xmlHandler nodeXmlHandler) inboundHandle(context *context) (int, string) {

	message := context.parameter["recvMessage"]
	var requestNode Node
	err := xml.Unmarshal([]byte(message), &requestNode)
	if err != nil {
		log.Println("字符串不能解析为对应类：", message)
		return -1, err.Error()
	}


	context.node = requestNode
	log.Println("当前请求节点：", requestNode.Name)

	return 0, ""
}

/**
  将结构数据转化为xml文本数据
*/
func (xmlHandler nodeXmlHandler) outboundHandle(context *context) (int, string) {

	node := context.node
	responseBytes, err := xml.Marshal(node)
	if err != nil {
		log.Println("报文编码失败：", node)
		return -1, err.Error()
	}

	log.Println("发送的报文内容为：\n", string(responseBytes))
	context.nodeBytes = responseBytes

	return 0, ""
}
