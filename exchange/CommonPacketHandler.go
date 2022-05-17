package exchange

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
	"strconv"
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
func (ph packetHandler) inboundHandle(context *context) ExchangeError{

	conn := context.conn
	reader := bufio.NewReader(conn)
	var exchangeError ExchangeError
	for {
		peek, err := reader.Peek(int(ph.contentLen))
		if err != nil {
			log.Println("读取长度错误：", err)
			exchangeError = newExchangeErrorByParams(515, []string{strconv.Itoa(int(ph.contentLen))})
			log.Println(GetErrorStackf(err, exchangeError.errMsg))
			return exchangeError
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
				exchangeError = newExchangeErrorByParams(516, []string{strconv.Itoa(int(ph.contentLen))})
				log.Println(GetErrorStackf(err, exchangeError.errMsg))
				return exchangeError
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
			exchangeError = newExchangeError(517)
			log.Println(GetErrorStackf(err, exchangeError.errMsg))
			return exchangeError
		}

		//处理完整数据
		exchangeError = dataHandle(data, context)
		if exchangeError.IsFail() {
			return exchangeError
		}
		if context.percent == strconv.Itoa(100) {
			break
		}
	}

	return exchangeError
}

/**
  发送包处理信息
*/
func (ph packetHandler) outboundHandle(context *context) ExchangeError {

	conn := context.conn
	// 发送报文
	message := context.sendBytes
	exchangeError := sendMessage(message, conn, "001", context.transCode)
	if exchangeError.IsFail() {
		return exchangeError
	}
	// 发送文件
	files := context.sendFiles
	if len(files) > 0 {
		for _, file := range files {
			fileBytes, exchangeError := readFile(file)
			if exchangeError.IsFail() {
				return exchangeError
			}
			// 发送报文
			pkg, err := packageMessage(context.transCode, mergePercent("099", fileBytes))
			if err != nil {
				exchangeError = newExchangeErrorByParams(511, []string{context.transCode})
				log.Println(GetErrorStackf(err, exchangeError.errMsg))
				return exchangeError
			}
			_, err = conn.Write(pkg.Bytes())
			if err != nil {
				exchangeError = newExchangeErrorByParams(512, []string{context.transCode})
				log.Println(GetErrorStackf(err, exchangeError.errMsg))
				return exchangeError
			}
		}
	}
	// 发送结束标志
	pkg, err := packageMessage(context.transCode, mergePercent("100", make([]byte, 0)))
	if err != nil {
		exchangeError = newExchangeErrorByParams(511, []string{context.transCode})
		log.Println(GetErrorStackf(err, exchangeError.errMsg))
		return exchangeError
	}
	_, err = conn.Write(pkg.Bytes())
	if err != nil {
		exchangeError = newExchangeErrorByParams(512, []string{context.transCode})
		log.Println(GetErrorStackf(err, exchangeError.errMsg))
		return exchangeError
	}

	return newExchangeError(0)
}


/***
  处理发送过来的数据
*/
func dataHandle(data []byte, context *context) ExchangeError {
	// 阶段标志
	percent := string(data[4:7])
	dataLen := len(data)
	var params = []string{context.transCode}
	var exchangeError = newExchangeError(0)
	if percent == "001" {
		// 数据内容
		if dataLen <= 7 {
			exchangeError = newExchangeErrorByParams(512, params)
			log.Println(exchangeError.errMsg)
			return exchangeError
		}
		message := string(data[7:])
		context.parameter["recvMessage"] = message
		context.percent = percent
		log.Println("接收到的报文：\n", message)
		return exchangeError
	} else if percent == "099" {
		// 保存到文件，文件可以是空文件
		context.percent = percent
		var fileData []byte
		if dataLen <= 7 {
			fileData = []byte{}
		} else {
			fileData = data[7:]
		}
		exchangeError := saveFile(context.transCode, fileData)
		if exchangeError.IsFail() {
			return exchangeError
		}
		fileIndex := len(context.recvFiles) - 1
		context.recvFiles[fileIndex] = exchangeError.filePath
		return exchangeError
	} else if percent == "100" {
		// 消息接收完成，跳出循环
		context.percent = percent
		return exchangeError
	}

	exchangeError = newExchangeErrorByParams(514, []string {context.transCode, percent})
	log.Println(exchangeError.errMsg)
	log.Printf("未失败的发送阶段【%s】\n", percent)
	return exchangeError
}

/**
  直接读文本，再转化为字节
  后期可优化为直接读取字节，少一道中转，效率更高
*/
func readFile(path string) ([]byte, ExchangeError) {

	var exchangeError ExchangeError
	file, err := os.Open(path)
	if err != nil {
		exchangeError = newExchangeErrorByParams(603, []string{path})
		log.Println(GetErrorStackf(err, exchangeError.errMsg))
		log.Printf("文件【%s】打开失败\n", path)
		return nil, exchangeError
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
				exchangeError = newExchangeErrorByParams(604, []string{path})
				log.Println(GetErrorStackf(err, exchangeError.errMsg))
				return nil, exchangeError
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

	return fileBytes, newExchangeError(0)
}
