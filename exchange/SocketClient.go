package exchange

import (
	"log"
	"net"
)

func init() {
	//初始化设置日志格式
	log.SetFlags(log.Lshortfile |log.Lmicroseconds | log.Ldate)
}
/**
  调用对应服务端进行交易请求
*/
func SendClient(hzRequest *HzbankRequest, files []string) (*HzbankResponse, []string, *Status) {

	// 初始化返回状态
	status := Status{
		ErrorCode: 0,
		ErrorMsg:  "",
	}
	if hzRequest == nil {
		status.ErrorCode = -1
		status.ErrorMsg = "请求对象不能为空"
		return nil, nil, &status
	}

	// 1. 开始校验发送节点和交易是否存在配置中

	// 2. 初始化输入输出执行器
	inHandlers, outHandlers := initExcchange()
	// 3. 连接服务端， 根据节点拿到信息并且连接
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Printf("服务端【%s】连接失败\n", "127.0.0.1:8080")
		log.Println("连接错误信息：", err)
		status.ErrorCode = -1
		status.ErrorMsg = "服务器连接失败"
		return nil, nil, &status
	}

	// 初始化发送上下文信息
	context := context{
		conn:      conn,
		node:      node{},
		transCode: hzRequest.Header.TransCode,
		message:   hzbankParameter{
			request: *hzRequest,
		},
        parameter: make(map[string]string),
		sendFiles: files,
		percent:   "001",
	}
	if len(outHandlers) > 0 {
		for _, outHandler := range outHandlers {
			errCode, msg := outHandler.outboundHandle(&context)
			if errCode != 0 {
				status.ErrorCode =errCode
				status.ErrorMsg = msg
				log.Printf("业务处理器【%s】执行交易【%s】失败>>>>>>>>>\n", outHandler.getName(), context.transCode)
				return nil, nil, &status
			}
		}
	}

	if len(inHandlers) > 0 {
		for _, inHandler := range inHandlers {
			errCode, msg := inHandler.inboundHandle(&context)
			if errCode != 0 {
				status.ErrorCode =errCode
				status.ErrorMsg = msg
				log.Printf("业务处理器【%s】执行交易【%s】失败>>>>>>>>>\n", inHandler.getName(), context.transCode)
				return nil, nil, &status
			}
		}
	}

    return &context.message.response, context.recvFiles, &status
}
