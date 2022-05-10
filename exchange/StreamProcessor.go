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
	  新增节点输入流链
	*/
	addNodeInboundHandler(nodeInboundHandler InboundHandler) error

	/***
	  新增响应输出流链
	 */
	addNodeOutboundHandler(nodeOutboundHandler OutboundHandler) error
}

/***
  流对象
*/
type StreamProcessor struct {
	inboundLen int
	outboundLen int
	nodeInboundLen int
	nodeOutboundLen int
	/** handler最大长度限制 */
	maxLen      int
	inboundHandlers [] InboundHandler
	outboundHandlers [] OutboundHandler
	nodeInboundHandlers [] InboundHandler
	nodeOutboundHandlers []OutboundHandler
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

func (sp *StreamProcessor) addNodeInboundHandler(nodeInboundHandler InboundHandler) error {
	if sp.nodeInboundLen >= sp.maxLen {
		return errors.New("node输入流处理器超过最大限制个数")
	}
	sp.nodeInboundHandlers[sp.nodeInboundLen] = nodeInboundHandler
	sp.nodeInboundLen = sp.nodeInboundLen + 1

	return nil
}

func (sp *StreamProcessor) addNodeOutboundHandler(nodeOutboundHandler OutboundHandler) error {
	if sp.nodeOutboundLen >= sp.maxLen {
		return errors.New("node输出流处理器超过最大限制个数")
	}
	sp.nodeOutboundHandlers[sp.nodeOutboundLen] = nodeOutboundHandler
	sp.nodeOutboundLen = sp.nodeOutboundLen + 1

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
		nodeInboundLen: 0,
		nodeOutboundLen: 0,
		maxLen: len,
		inboundHandlers:  make([]InboundHandler, len, len) ,
		outboundHandlers: make([]OutboundHandler, len, len),
		nodeInboundHandlers: make([]InboundHandler, len, len),
		nodeOutboundHandlers: make([]OutboundHandler, len, len),
	}

	// 不进行错误判断是因为本身新增handler长度已知
	// 新增输入流处理器

	err := streamProcessor.addInboundHandler(packetHandler{
		name: "报文接收器",
		// 通过实际需要进行字节长度控制， 当前为了传输数据小一点，这里只使用4字节，int32（如果需要可以设置8字节，使用int64）
		contentLen: 4,
	})
	if err != nil {
		return nil ,err
	}
	err = streamProcessor.addInboundHandler(exchangeXmlHandler{name: "xml报文解析器"})
	if err != nil {
		return nil ,err
	}
	// 新增输出流处理器
	err = streamProcessor.addOutboundHandler(exchangeXmlHandler{name: "xml报文转文本器"})
	if err != nil {
		return nil ,err
	}
	err = streamProcessor.addOutboundHandler(nodeCheckHandler{name: "节点校验处理器"})
	if err != nil {
		return nil ,err
	}
	err = streamProcessor.addOutboundHandler(packetHandler{name: "公共报文处理器", contentLen: 4})
	if err != nil {
		return nil ,err
	}

	// 新增节点输入处理链
	err = streamProcessor.addNodeInboundHandler(nodePacketHandler{name: "node信息接收处理器", contentLen: 4})
	if err != nil {
		return nil ,err
	}
	err = streamProcessor.addNodeInboundHandler(exchangeXmlHandler{name: "node信息xml报文解析器"})
	if err != nil {
		return nil ,err
	}

	// 新增节点输出流处理器， 报错不会走该流程
	err = streamProcessor.addNodeOutboundHandler(nodeXmlHandler{name: "node信息xml报文转文本器"})
	if err != nil {
		return nil ,err
	}
	err = streamProcessor.addNodeOutboundHandler(nodePacketHandler{name: "node信息发送处理器", contentLen: 4})
	if err != nil {
		return nil, err
	}

	return &streamProcessor, nil
}