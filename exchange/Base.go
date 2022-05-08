package exchange

import "net"

/**
  连接交互上下文信息
 */
type context struct {

	// 连接
	conn net.Conn
	// 节点信息
	nodes []Node
	// 实际获取的IP
	host string
	// 交易码
	transCode string
	// 新增Map对象，可以同过上下文传递信息
	parameter map[string]string
	// 发送或者接收的报文对象
	message hzbankParameter
	// 发送报文字节
	sendBytes []byte
	// 接收到的文件路径
	recvFiles []string
	// 发送的文件路径
	sendFiles []string
	// 发送的阶段
	percent string
	// 响应头处理器
	streamProcessor StreamProcessor
}

/**
  定义IB2处理接口
*/
type InboundHandler interface {
	// 获取业务处理器名称方法
	getName() string

	// 接收数据处理方法
	inboundHandle(ctx *context) (int, string)
}

type OutboundHandler interface {

	// 获取业务处理器名称方法
	getName() string
	// 发送数据处理方法
	outboundHandle(ctx *context) (int, string)
}
