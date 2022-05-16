package exchange

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

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
	// 当前请求节点
	node Node
	nodeBytes []byte
}

/**
  定义IB2处理接口
*/
type InboundHandler interface {
	// 获取业务处理器名称方法
	getName() string

	// 接收数据处理方法
	inboundHandle(ctx *context) ExchangeError
}

type OutboundHandler interface {

	// 获取业务处理器名称方法
	getName() string
	// 发送数据处理方法
	outboundHandle(ctx *context) ExchangeError
}

func sendMessage(message []byte, conn net.Conn, percent string, transCode string) ExchangeError {

	var exchangeError ExchangeError
	if message == nil {
		exchangeError = newExchangeErrorByParams(510, []string{transCode})
		exchangeError.ErrorPrintln(nil)
		return exchangeError
	}

	// 发送报文
	pkg, err := packageMessage(transCode, mergePercent(percent, message))
	if err != nil {
		exchangeError = newExchangeErrorByParams(511, []string{transCode})
		exchangeError.ErrorPrintln(err)
		return exchangeError
	}

	_, err = conn.Write(pkg.Bytes())
	if err != nil {
		exchangeError = newExchangeErrorByParams(512, []string{transCode})
		exchangeError.ErrorPrintln(err)
		return exchangeError
	}

	return newExchangeError(0)

}

/**
  合并阶段信息到包内容中
*/
func mergePercent(percent string, message []byte) []byte {

	if len(message) == 0 {
		return []byte(percent)
	}
	var arrayBytes = make([][]byte, 2, 2)
	arrayBytes[0] = []byte(percent)
	message = bytes.Join(arrayBytes, message)

	return message
}

/**
  封装信息为字节包
*/
func packageMessage(transCode string, message []byte) (*bytes.Buffer, error) {
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// 写入包长度
	err := binary.Write(pkg, binary.BigEndian, length)
	if err != nil {
		log.Printf("交易【%s】写入数据长度失败\n", transCode)
		return nil, err
	}
	// 写入数据内容
	err = binary.Write(pkg, binary.BigEndian, message)
	if err != nil {
		log.Printf("交易【%s】写入数据报错>>>>>>>>\n", transCode)
		return nil, err
	}

	return pkg, nil
}