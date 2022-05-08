package exchange

import (
	"errors"
)

type handler interface {
	/**
	  新增输入流处理器
	 */
	addInboundHandler(inHandler InboundHandler) error

	/**
	  新增输出流处理器
	 */
	addOutboundHandler(outHandler OutboundHandler) error

	/**
	  新增输出流处理器
	*/
	addHeaderHandler(headerHandler OutboundHandler) error
}

/***
  流对象
*/
type StreamProcessor struct {
	inboundLen int
	outboundLen int
	headerLen int
	/** handler最大长度限制 */
	maxLen      int
	inboundHandlers [] InboundHandler
	outboundHandlers [] OutboundHandler
	headerHandlers []InboundHandler
}

func (sp *StreamProcessor) addInboundHandler(inHandler InboundHandler) error {
	if sp.inboundLen >= sp.maxLen {
        return errors.New("输入流处理器超过最大限制个数")
	}
    sp.inboundHandlers[sp.inboundLen] = inHandler
    sp.inboundLen = sp.inboundLen + 1

    return nil
}

func (sp *StreamProcessor) addOutboundHandler(outHandler OutboundHandler) error {
    if sp.outboundLen >= sp.maxLen {
    	return errors.New("输出流处理器超过最大限制个数")
	}
	sp.outboundHandlers[sp.outboundLen] = outHandler
	sp.outboundLen = sp.outboundLen + 1

	return nil
}

func (sp *StreamProcessor) addHeaderHandler(headerHandler InboundHandler) error {
	if sp.headerLen >= sp.maxLen {
		return errors.New("输出流处理器超过最大限制个数")
	}
	sp.headerHandlers[sp.headerLen] = headerHandler
	sp.headerLen = sp.headerLen + 1

	return nil
}

/**
  创建流处理器
*/
func createStreamProcessor(len int) (*StreamProcessor, error) {

	// 初始化流处理器
	streamProcessor := StreamProcessor{
		inboundLen: 0,
		outboundLen: 0,
		headerLen: 0,
		maxLen: len,
		inboundHandlers:  make([]InboundHandler, len, len) ,
		outboundHandlers: make([]OutboundHandler, len, len),
		headerHandlers: make([]InboundHandler, len, len),
	}

	// 不进行错误判断是因为本身新增handler长度已知
	// 新增输入流处理器
	err := streamProcessor.addInboundHandler(packetHandler{
		name: "报文接收器",
		// 通过实际需要进行字节长度控制， 当前为了传输数据小一点，这里只使用4字节，int32（如果需要可以设置8字节，使用int64）
		contentLen: 4,
	})
	err = streamProcessor.addInboundHandler(exchangeXmlHandler{name: "xml报文解析器"})

	// 新增输出流处理器
	err = streamProcessor.addOutboundHandler(exchangeXmlHandler{name: "xml报文转文本器"})
	err = streamProcessor.addOutboundHandler(packetHandler{name: "报文发送器", contentLen: 4})

	// 新增header响应处理器
	// 新增输出流处理器
	err = streamProcessor.addHeaderHandler(nodeCheckHandler{name: "节点响应处理器", contentLen: 4})
	err = streamProcessor.addHeaderHandler(nodeXmlHandler{name: "节点xml报文转文本器"})

	if err != nil {
		return nil, err
	}

	return &streamProcessor, nil
}