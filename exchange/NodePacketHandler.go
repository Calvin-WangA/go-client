package exchange

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"strconv"
)

/**
  报文解析之前进行节点相关信息校验
*/

type nodePacketHandler struct {
	name       string
	contentLen int32
}

func (checkHandler nodePacketHandler) getName() string {
	return checkHandler.name
}

/**
  连接交互0阶段，进行节点信息验证
*/
func (checkHandler nodePacketHandler) inboundHandle(context *context) ExchangeError {

	conn := context.conn
	reader := bufio.NewReader(conn)
	var exchangeError ExchangeError
	for {
		peek, err := reader.Peek(int(checkHandler.contentLen))
		if err != nil {
			log.Printf("输入处理器【%s】读取长度错误>>>>>>>>>>>\n", checkHandler.getName())
			exchangeError = newExchangeErrorByParams(515, []string{strconv.Itoa(int(checkHandler.contentLen))})
			log.Println(GetErrorStackf(err, exchangeError.errMsg))
			return exchangeError
		}

		buffer := bytes.NewBuffer(peek)
		var length int32
		err = binary.Read(buffer, binary.BigEndian, &length)
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				log.Printf("输入处理器【%s】读取长度数据出错>>>>>>>>>>>>\n", checkHandler.getName())
				exchangeError = newExchangeErrorByParams(516, []string{strconv.Itoa(int(checkHandler.contentLen))})
				log.Println(GetErrorStackf(err, exchangeError.errMsg))
				return exchangeError
			}
		}

		if int32(reader.Buffered()) < (length + checkHandler.contentLen) {
			continue
		}

		data := make([]byte, length+checkHandler.contentLen)
		_, err = reader.Read(data)
		if err != nil {
			log.Printf("输入处理器【%s】读取数据内容失败>>>>>>>>>>>>\n", checkHandler.getName())
			exchangeError = newExchangeError(517)
			log.Println(GetErrorStackf(err, exchangeError.errMsg))
			return exchangeError
		}

		// 数据封装
		nodePacketHandle(data, context)
		break
	}

	return newExchangeError(0)
}

func nodePacketHandle(data []byte, context *context) {
	// 阶段标志
	message := string(data[7:])
	context.parameter["recvMessage"] = message
	// 设置阶段信息正式报文交互
	context.percent = "001"
	log.Println("接收到的报文：\n", message)
}

/**
  直接调用公共方法进行使用即可。
*/
func (checkHandler nodePacketHandler) outboundHandle(ctx *context) ExchangeError {
	return sendMessage(ctx.nodeBytes, ctx.conn, ctx.percent, ctx.transCode)
}
