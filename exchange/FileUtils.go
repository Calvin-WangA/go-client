package exchange

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

/**
  工具文件创建类
 */
func CreateFile (path string) (*os.File, error) {

	// 创建文件目录
	fileDir := filepath.Dir(path)
	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		log.Printf("文件目录【%s】创建失败\n", fileDir)
		log.Println("目录创建失败原因：", err)
		return nil, err
	}
	// 创建文件
    file, err := os.Create(path)
    if err != nil {
		log.Printf("文件【%s】创建失败\n", path)
		log.Println("文件创建失败原因：", err)
		return nil, err
	}

	return file, nil
}

/**
  直接先保存文件，防止文件太多内存占爆
  返回0表示保存成功， 并且返回的内容为报错的文件路径，否则返回错误码和报错信息
*/
func saveFile(transCode string, fileBytes []byte) ExchangeError {

	// 文件路径
	filePath := properties["file_path"]
    filePath = filePath + transCode + GetTimeNanoPlusRandom()
	var exchangeError ExchangeError
	file, err := CreateFile(filePath)
	if err != nil {
        exchangeError = newExchangeErrorByParams(600, []string{filePath})
		log.Println(GetErrorStackf(err, exchangeError.errMsg))
		return exchangeError
	}
	writer := bufio.NewWriter(file)
	_, err = writer.Write(fileBytes)
	if err != nil {
		log.Printf("文件【%s】保存失败\n", filePath)
		exchangeError = newExchangeErrorByParams(601, []string{filePath})
		log.Println(GetErrorStackf(err, exchangeError.errMsg))
		return exchangeError
	}

	err = writer.Flush()
	if err != nil {
		log.Println("缓存刷入磁盘失败：", err)
		exchangeError = newExchangeErrorByParams(602, []string{filePath})
		log.Println(GetErrorStackf(err, exchangeError.errMsg))
		return exchangeError
	}
	// 写完关闭文件
	defer file.Close()

	exchangeError = newExchangeError(0)
	exchangeError.filePath = filePath
	return exchangeError
}

/**
  获取纳秒加一个10000随机数的字符串
 */
func GetTimeNanoPlusRandom() string {

	timeUnixNano := strconv.Itoa(int(time.Now().UnixNano()))
	randNum := strconv.Itoa(int(rand.Int63n(10000)))

	return timeUnixNano + randNum
}
