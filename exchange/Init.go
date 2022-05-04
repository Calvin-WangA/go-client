package exchange

/**
  初始化交互接收和发送需要调用的处理器
 */
func initExcchange() ([]InboundHandler, []OutboundHandler) {

	// slice只能实际长度跟内容匹配，否则会读取到空报错
	inboundHandlers := make([]InboundHandler, 2, 20)

	inboundHandlers[0] = packetHandler{
		name:       "报文接收器",
		// 通过实际需要进行字节长度控制， 当前为了传输数据小一点，这里只使用4字节，int32（如果需要可以设置8字节，使用int64）
		contentLen: 4,
	}
	inboundHandlers[1] = exchangeXmlHandler{name: "xml报文解析器"}

	// 初始化输出处理器
	outboundHandlers := make([]OutboundHandler, 2, 20)
	outboundHandlers[0] = exchangeXmlHandler{name: "xml报文转文本器"}
	outboundHandlers[1] = packetHandler{name: "报文发送器", contentLen: 4}

	return inboundHandlers, outboundHandlers
}
