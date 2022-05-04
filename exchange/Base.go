package exchange

import "net"

/***
  配置地址信息
 */
type address struct {
	// 主机域名
	Host string
	// 端口
	Port int
}

type encryption struct {
	// 使用的算法
	Name string
	// 对应文件路径
	Path string
	// 公钥名
	PublicKey string
	// 私钥名
	PrivateKey string
    // 密码是否使用文件格式
	Flag bool
}
/**
  IB2处理时上线文结构
 */
type node struct {
	// 节点名
	Name string
	// 节点代码
	Code string
	// 地址信息
    Cddress address
	// 加密信息
	Encryption encryption
	// 编码格式
    Encode string
}
type context struct {

	// 连接
	conn net.Conn
	// 节点信息
	node node
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
	recvFiles   []string
	// 发送的文件路径
	sendFiles   []string
    // 发送的阶段
	percent string
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

