package exchange

import (
	"encoding/xml"
	"log"
)
/**
   xml文件解析
 */

type HzbankRequest struct {
	XMLName xml.Name `xml:"Hzbank"`
	Header Header `xml:"Header"`
	Body   RequestBody   `xml:"RequestBody"`
}

type HzbankResponse struct {
	XMLName xml.Name `xml:"Hzbank"`
	Header Header `xml:"Header"`
	Body   ResponseBody   `xml:"ResponseBody"`
}

/**
  交易头部信息
 */
type Header struct {
	// 交易流水号
	SerialNo string `xml:"SerialNo"`
	// 交易名称
	Name string `xml:"Name"`
	// 交易码
	TransCode string `xml:"TransCode"`
}


/**
  请求参数体
 */
type RequestBody struct {
	// 客户号
	ClientNo string `xml:"ClientNo"`
}


/**
  响应参数体
 */
type ResponseBody struct {
	Status Status
}

/**
  处理状态
 */
type Status struct {
	// 交易状态，0成功，否则失败
	ErrorCode int
	// 非零状态的错误描述
	ErrorMsg  string
}

// 交互报文
type hzbankParameter struct {
	request HzbankRequest
	response HzbankResponse
}

type exchangeXmlHandler struct {
	name string
}

/**
  返回业务处理器名称
 */
func (xmlHandler exchangeXmlHandler) getName() string {
	return xmlHandler.name
}
/**
  解析文本xml为结构对象
 */
func (xmlHandler exchangeXmlHandler) inboundHandle(context *context) (int, string) {

	message := context.parameter["recvMessage"]
	var hzbankResponse HzbankResponse
	err := xml.Unmarshal([]byte(message), &hzbankResponse)
	if err != nil {
		log.Println("字符串不能解析为对应类：", message)
		log.Println("解析错误信息：", err)
		return -1, err.Error()
	}
    // 设置交易码
	context.transCode = hzbankResponse.Header.TransCode
	requestParameter := hzbankParameter{
		request:  HzbankRequest{},
		response: hzbankResponse,
	}
	context.message = requestParameter
	log.Println("交易状态为：", hzbankResponse.Body.Status.ErrorMsg )

    return 0, ""
}

/**
  将结构数据转化为xml文本数据
 */
func (xmlHandler exchangeXmlHandler) outboundHandle(context *context) (int, string) {

	request := context.message.request
	requestBytes, err := xml.Marshal(request)
    if err != nil {
    	log.Printf("报文编码失败：", request)
    	log.Println("报文编码失败原因：", err)
    	return -1, err.Error()
	}

	log.Println("发送的报文内容为：\n", string(requestBytes))
	context.sendBytes = requestBytes
	return 0, ""
}