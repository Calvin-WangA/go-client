package exchange

import (
	"log"
	"net"
)

/***
  packet发送节点信息之后不能马上调用read读取等待结果，否则接收不到正确的消息问题：
  解決方案：
      1. 采用睡眠1纳秒方案，不太合理。
      2. 需要其他手段正确解決
*/

func init() {
	//初始化设置日志格式
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
}

/**
  调用对应服务端进行交易请求
*/
func SendClient(nodeCode string, hzRequest *HzbankRequest, files []string) (*HzbankResponse, []string, *ExchangeError) {

	// 初始化返回状态
	var exchangeError ExchangeError
	if hzRequest == nil {
		exchangeError = newExchangeError(300)
		return nil, nil, &exchangeError
	}

	// 1. 开始校验发送节点和交易是否存在配置中, 应用初始化的时候去用，而不是每次调用
	node := getNode(nodeCode, IB2_NODES.Nodes)
	if node == nil {
        params := []string {nodeCode}
		exchangeError = newExchangeErrorByParams(301, params)
		return nil, nil, &exchangeError
	}

	// 2. 初始化输入输出执行器
	streamProcessor, err := createStreamProcessor(3)
	if err != nil {
		params := []string {err.Error()}
		exchangeError = newExchangeErrorByParams(302, params)
		log.Println(GetErrorStack(err, exchangeError.errMsg))
		return nil, nil, &exchangeError
	}

	// 3. 连接服务端， 根据节点拿到信息并且连接
	address := getAddress(*node)
	if address == "" {
		params := []string {node.Code}
		exchangeError = newExchangeErrorByParams(303, params)
		return nil, nil, &exchangeError
	}
	conn, err := net.Dial(node.Protocol, address)
	if err != nil {
		params := []string {address}
		exchangeError = newExchangeErrorByParams(404, params)
		log.Println(GetErrorStack(err, exchangeError.errMsg))
		return nil, nil, &exchangeError
	}

	// 初始化发送上下文信息
	context := context{
		conn: conn,
		node: *NODE_SELF,
		nodes:     nil,
		transCode: hzRequest.Header.TransCode,
		message: hzbankParameter{
			request: *hzRequest,
		},
		parameter:       make(map[string]string),
		sendFiles:       files,
		percent:         "000",
		streamProcessor: *streamProcessor,
	}

	outHandlers := streamProcessor.outboundHandlers
	handlerLen := streamProcessor.outboundLen
	var outboundHandler OutboundHandler
	if handlerLen > 0 {
		for index := 0; index < handlerLen; index++ {
			outboundHandler = outHandlers[index]
			exchangeError = outboundHandler.outboundHandle(&context)
			if exchangeError.IsFail() {
				log.Printf("业务处理器【%s】执行交易【%s】失败>>>>>>>>>\n", outboundHandler.getName(), context.transCode)
				return nil, nil, &exchangeError
			}
		}
	}

	inHandlers := streamProcessor.inboundHandlers
	handlerLen = streamProcessor.inboundLen
	var inboundHandler InboundHandler
	if handlerLen > 0 {
		for index := 0; index < handlerLen; index++ {
			inboundHandler = inHandlers[index]
			exchangeError = inboundHandler.inboundHandle(&context)
			if exchangeError.IsFail() {
				log.Printf("业务处理器【%s】执行交易【%s】失败>>>>>>>>>\n", inboundHandler.getName(), context.transCode)
				return nil, nil, &exchangeError
			}
		}
	}

	exchangeError = newExchangeError(0)
	return &context.message.response, context.recvFiles, &exchangeError
}
