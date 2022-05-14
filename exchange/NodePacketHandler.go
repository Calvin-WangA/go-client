package exchange

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
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
func (checkHandler nodePacketHandler) inboundHandle(context *context) (int, string) {

	conn := context.conn
	reader := bufio.NewReader(conn)
	for {
		peek, err := reader.Peek(int(checkHandler.contentLen))
		if err != nil {
			log.Printf("输入处理器【%s】读取长度错误>>>>>>>>>>>\n", checkHandler.getName())
			log.Println("错误信息：", err)
			return -1, err.Error()
		}

		buffer := bytes.NewBuffer(peek)
		var length int32
		err = binary.Read(buffer, binary.BigEndian, &length)
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				log.Printf("输入处理器【%s】读取长度数据出错>>>>>>>>>>>>\n", checkHandler.getName())
				log.Println("错误信息：", err)
				return -2, err.Error()
			}
		}

		if int32(reader.Buffered()) < (length + checkHandler.contentLen) {
			continue
		}

		data := make([]byte, length+checkHandler.contentLen)
		_, err = reader.Read(data)
		if err != nil {
			log.Printf("输入处理器【%s】读取数据内容失败>>>>>>>>>>>>\n", checkHandler.getName())
			log.Println("错误信息：", err)
			return -3, err.Error()
		}

		// 数据封装
		nodePacketHandle(data, context)
		break
	}

	return 0, ""
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
func (checkHandler nodePacketHandler) outboundHandle(ctx *context) (int, string) {
	return sendMessage(ctx.nodeBytes, ctx.conn, ctx.percent, ctx.transCode)
}
