package exchange

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"io"
	"log"
	"net"
	"os"
	"time"
)

/**
  定义粘包结构体
*/
type packetHandler struct {
	name string
	// 头部长度 ,当前最大只能用int32, 如果特殊需要可以设置int64, 长度字节相应也设置为8
	contentLen int32
}

func (ph packetHandler) getName() string {
	return ph.name
}

/**
  接收包处理信息
*/
func (ph packetHandler) inboundHandle(context *context) (int, string) {

	conn := context.conn
	reader := bufio.NewReader(conn)
	for {
		peek, err := reader.Peek(int(ph.contentLen))
		if err != nil {
			log.Println("读取长度错误：", err)
			return -1, "Peek长度错误：" + err.Error()
		}

		// 读取长度
		buffer := bytes.NewBuffer(peek)
		var length int32
		err = binary.Read(buffer, binary.BigEndian, &length)
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				log.Println("读取长度出错，错误信息为：", err)
				return -2, "读取数据长度错误：" + err.Error()
			}
		}
		// 判断是否读完整，否则继续等待
		if int32(reader.Buffered()) < (length + ph.contentLen) {
			log.Printf("数据长度【%d】，继续等待数据\n", length)
			continue
		}

		// 获取传递的数据内容
		data := make([]byte, length+ph.contentLen)
		_, err = reader.Read(data)
		if err != nil {
			log.Println("读取数据失败，", err)
			return -3, "读取数据内容错误：" + err.Error()
		}

		//处理完整数据
		n, msg := dataHandle(data, context)
		if n == -1 {
			return n, msg
		}
		if n == 100 {
			break
		}
	}

	return 0, ""
}

/**
  发送包处理信息
*/
func (ph packetHandler) outboundHandle(context *context) (int, string) {

	conn := context.conn
	// 发送节点信息进行校验，校验通过在发送报文和文件
	nodeBytes, errCode, msg := getNodeBytes()
	if errCode != 0 {
		return errCode, msg
	}
	errCode, msg = sendMessage(nodeBytes, conn, "000", context.transCode)
	if errCode != 0 {
		return errCode, msg
	}
    // 解決不能正确读取服务端返回业务处理结果的问题，需要其他手段正确解決
	time.Sleep(1)
	// 接收响应，正常才进行发送参数信息
	errCode, msg = headerCheck(context)
	if errCode != 0 {
		return errCode, msg
	}

	// 发送报文
	message := context.sendBytes
	errCode, msg = sendMessage(message, conn, "001", context.transCode)
	if errCode != 0 {
		return errCode, msg
	}
	// 发送文件
	files := context.sendFiles
	if len(files) > 0 {
		for _, file := range files {
			fileBytes, errCode, msg := readFile(file)
			if errCode != 0 {
				return errCode, msg
			}
			// 发送报文
			pkg, err := packageMessage(context.transCode, mergePercent("099", fileBytes))
			if err != nil {
				return -2, err.Error()
			}
			_, err = conn.Write(pkg.Bytes())
			if err != nil {
				log.Printf("交易【%s】报文发送失败>>>>>>>>>>\n", context.transCode)
				log.Println("文件发送错误信息：", err)
				return -2, err.Error()
			}
		}
	}
	// 发送结束标志
	pkg, err := packageMessage(context.transCode, mergePercent("100", make([]byte, 0)))
	if err != nil {
		return -3, err.Error()
	}
	_, err = conn.Write(pkg.Bytes())
	if err != nil {
		log.Printf("交易【%s】结束报文发送失败>>>>>>>>>>\n", context.transCode)
		log.Println("结束报文发送错误信息：", err)
		return -3, err.Error()
	}

	return 0, ""
}

func sendMessage(message []byte, conn net.Conn, percent string, transCode string) (int, string) {

	if message == nil {
		return -2, "待发送信息为空"
	}

	// 发送报文
	pkg, err := packageMessage(transCode, mergePercent(percent, message))
	if err != nil {
		return -1, err.Error()
	}

	_, err = conn.Write(pkg.Bytes())
	if err != nil {
		log.Printf("交易【%s】报文发送失败>>>>>>>>>>\n", transCode)
		log.Println("报文发送错误信息：", err)
		return -1, err.Error()
	}

	return 0, ""

}

func getNodeBytes() ([]byte, int, string) {
	node := Node{
		Name:       "手机银行",
		Code:       "FSTS",
		Encode:     "GBK",
		Addresses:  Addresses{},
		Encryption: Encryption{},
	}
	nodeBytes, err := xml.Marshal(&node)
	if err != nil {
		return nil, -1, err.Error()
	}

	return nodeBytes, 0, ""
}

/**
  通过输入流获取结果
*/
func headerCheck(context *context) (int, string) {

	streamProcessor := context.streamProcessor
	headerHandlers := streamProcessor.headerHandlers
	headerLen := streamProcessor.headerLen
	var headerHandler InboundHandler
	if len(headerHandlers) > 0 {
		for index:= 0; index < headerLen; index++ {
			headerHandler = headerHandlers[index]
			errCode, msg := headerHandler.inboundHandle(context)
			if errCode != 0 {
				log.Printf("业务处理器【%s】执行交易【%s】失败>>>>>>>>>\n", headerHandler.getName(), context.transCode)
				return errCode, msg
			}
			log.Printf("节点处理器【%s】执行完成>>>>>>>\n", headerHandler.getName())
		}
	}

	return 0, ""
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
		log.Println("写入数据长度错误信息：", err)
		return nil, err
	}
	// 写入数据内容
	err = binary.Write(pkg, binary.BigEndian, message)
	if err != nil {
		log.Printf("交易【%s】写入数据报错>>>>>>>>\n", transCode)
		log.Println("写入数据错误信息：", err)
		return nil, err
	}

	return pkg, nil
}

/***
  处理发送过来的数据
*/
func dataHandle(data []byte, context *context) (int, string) {
	// 阶段标志
	percent := string(data[4:7])
	dataLen := len(data)
	if percent == "001" {
		// 数据内容
		if dataLen <= 7 {
			return -2, "接收报文无内容"
		}
		message := string(data[7:])
		context.parameter["recvMessage"] = message
		context.percent = percent
		log.Println("接收到的报文：\n", message)
		return 1, ""
	} else if percent == "099" {
		// 保存到文件，文件可以是空文件
		context.percent = percent
		var fileData []byte
		if dataLen <= 7 {
			fileData = []byte{}
		} else {
			fileData = data[7:]
		}
		errCode, msg := saveFile(fileData)
		if errCode != 0 {
			return errCode, msg
		}
		fileIndex := len(context.recvFiles) - 1
		context.recvFiles[fileIndex] = msg
		return 99, ""
	} else if percent == "100" {
		// 消息接收完成，跳出循环
		context.percent = percent
		return 100, ""
	}

	log.Printf("未失败的发送阶段【%s】\n", percent)
	return -1, "未识辨的发送阶段错误信息"
}

/**
  直接读文本，再转化为字节
  后期可优化为直接读取字节，少一道中转，效率更高
*/
func readFile(path string) ([]byte, int, string) {

	file, err := os.Open(path)
	if err != nil {
		log.Printf("文件【%s】打开失败\n", path)
		log.Println("文件打开错误信息：", err)
		return nil, -1, err.Error()
	}

	// 读取文件
	reader := bufio.NewReader(file)
	var arrayBytes = make([][]byte, 2, 2)
	var fileBytes []byte
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Printf("文件【%s】内容读取失败, 错误信息\n", path)
				log.Println("内如读取错误信息：", err)
				return nil, -2, err.Error()
			}
		}
		// 合并字节信息
		if arrayBytes[0] == nil {
			fileBytes = []byte(line)
			arrayBytes[0] = fileBytes
		} else {
			fileBytes = bytes.Join(arrayBytes, []byte(line))
		}
		arrayBytes[0] = fileBytes

	}

	return fileBytes, 0, ""
}
