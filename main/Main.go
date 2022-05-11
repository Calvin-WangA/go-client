package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"exchange"
	"log"
	"net"
	"os"
	"strings"
)

func main() {

	request := exchange.HzbankRequest{
		XMLName: xml.Name{
			Space: "",
			Local: "Hzbank",
		},
		Header: exchange.Header{
			SerialNo:  "1234567",
			Name:      "购买交易",
			TransCode: "100200",
		},
		Body: exchange.RequestBody{ClientNo: "400383444"},
	}
	files := []string{"D:\\test\\gotest.txt"}
	response, _, status := exchange.SendClient("FSTS", &request, files)
	if status.ErrorCode != 0 {
		log.Printf("交易失败，错误码【%d】，错误信息【%s】\n", status.ErrorCode, status.ErrorMsg)
	}

	if response != nil {
		log.Printf("响应状态【%d】，信息描述【%s】\n", response.Body.Status.ErrorCode, response.Body.Status.ErrorMsg)
	}

}

// 测试报文
func exchangeTest() {
	// 1. 与服务端建立连接
	var conn, err = net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println("服务端连接失败》》》》》》》")
		return
	}

	pkgBuffer := packet([]byte("001" + `<Hzbank>
   <Header>
      <SerialNo>12345</SerialNo>
	  <Name>理财购买</Name>
	  <TransCode>100200</TransCode>
   </Header>
   <RequestBody>
      <ClientNo>400383444</ClientNo>
   </RequestBody>
</Hzbank>`))
	// 先发送报文
	conn.Write(pkgBuffer.Bytes())
	// 发送文件
	fileBytes := readFile("D:\\test\\gotest.txt")
	if fileBytes != nil {
		pkgBuffer = packet(fileBytes)
		conn.Write(pkgBuffer.Bytes())
	}
	// 发送结束信息
	pkgBuffer = packet([]byte("100"))
	conn.Write(pkgBuffer.Bytes())

	transCodeSocket(conn)
}

// 封装包
func packet(message []byte) bytes.Buffer {
	// 新增
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// 写入包长度
	err := binary.Write(pkg, binary.BigEndian, length)
	if err != nil {
		log.Println("包装报文头失败>>>>>>>>>>>", err)

	}

	log.Println("报文长度大小：", length)
	// 写入包内容
	err = binary.Write(pkg, binary.BigEndian, []byte(message))
	if err != nil {
		log.Println("写入报文内容失败>>>>>>>>", err)
	}

	return *pkg
}

func readFile(path string) []byte {

	file, err := os.Open(path)
	if err != nil {
		log.Println("打开文件出错：", err)
		return nil
	}

	reader := bufio.NewReader(file)
	readBytes := make([]byte, 1024)
	n, err := reader.Read(readBytes)
	if err != nil {
		log.Println("读取文件错误：", err)
		return nil
	}

	// 拼接文件头
	percentBytes := []byte("099")
	var arrayBytes = make([][]byte, 2, 2)
	arrayBytes[0] = percentBytes
	readBytes = bytes.Join(arrayBytes, readBytes[0:n])
	log.Println("文件字节数为：", len(readBytes))

	return readBytes
}

func transCodeSocket(conn net.Conn) {
	for {
		// 接收服务端的消息
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Println("读取服务端消息失败>>>>>>")
			return
		}

		log.Printf("接收到个数【%d】值为【%s】\n", n, string(buf[:n]))
	}
}

/**
  客户端数据测试
*/
func socketInput(conn net.Conn) {

	// 2. 使用conn发送数据和接收数据
	var input = bufio.NewReader(os.Stdin)
	for {
		var s, _ = input.ReadString('\n')
		s = strings.TrimSpace(s)
		// 接收到Q就退出交互
		if strings.ToUpper(s) == "Q" {
			log.Printf("退出交互>>>>>>>>>")
			return
		}

		// 3. 写数据到客户端
		_, err := conn.Write([]byte(s))
		if err != nil {
			log.Println("数据发送失败>>>>>>>>>")
			return
		}

		// 接收服务端的消息
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Println("读取服务端消息失败>>>>>>")
			return
		}

		log.Printf("接收到个数【%d】值为【%s】\n", n, string(buf[:n]))
	}
}
