package exchange

import (
	"log"
)

type nodeCheckHandler struct {
	name string
}

func (nodeCheckHandler nodeCheckHandler) getName() string {
	return nodeCheckHandler.name
}

func (nodeCheckHandler nodeCheckHandler) outboundHandle(ctx *context) (int, string) {

	// 非节点校验阶段直接跳过
	if ctx.percent != "000" {
		return 0, ""
	}

	var errCode int
	var msg  string
	streamProcessor := ctx.streamProcessor
	// 调用输出返回需要相应的信息
	outboundHandlers := streamProcessor.nodeOutboundHandlers
	var outboundHandler OutboundHandler
	handlerLen := streamProcessor.nodeOutboundLen
	if handlerLen > 0 {
		for index := 0; index < handlerLen; index++ {
			outboundHandler = outboundHandlers[index]
			errCode, msg  := outboundHandler.outboundHandle(ctx)
			if errCode != 0 {
				log.Printf("Node业务处理器【%s】报错码【%d】错误信息【%s】\n", outboundHandler.getName(), errCode, msg)
				// 返回客户端报错信息
				break
			}
			log.Printf("Node发送业务处理器【%s】执行完成\n", outboundHandler.getName())
		}
	}
	if errCode != 0 {
		return errCode, msg
	}

	inboundHandlers := streamProcessor.nodeInboundHandlers
	var inboundHandler InboundHandler
	handlerLen = streamProcessor.nodeInboundLen
	if handlerLen > 0 {
		for index := 0; index < handlerLen; index++ {
			inboundHandler = inboundHandlers[index]
			errCode, msg = inboundHandler.inboundHandle(ctx)
			if errCode != 0 {
				log.Printf("Node业务处理器【%s】报错码【%d】错误信息【%s】\n", inboundHandler.getName(), errCode, msg)
				// 返回写入conn，告诉客户端处理报错, 并且直接返回
				break
			}
			log.Printf("Node接收业务处理器【%s】执行完成\n", inboundHandler.getName())
		}
	}
	if errCode != 0 {
		return errCode, msg
	}

	response := ctx.message.response
	return response.Body.Status.ErrorCode, response.Body.Status.ErrorMsg
}

