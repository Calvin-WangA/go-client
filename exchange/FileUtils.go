package exchange

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
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
		log.Printf("文件【%s】创建失败【%s】\n", path)
		log.Println("文件创建失败原因：", err)
		return nil, err
	}

	return file, nil
}

/**
  直接先保存文件，防止文件太多内存占爆
  返回0表示保存成功， 并且返回的内容为报错的文件路径，否则返回错误码和报错信息
*/
func saveFile(fileBytes []byte) (int, string) {

	// 文件路径
	path := "D:\\test\\ib2\\recv\\test1.txt"
	file, err := CreateFile(path)
	if err != nil {
		return -1, err.Error()
	}
	writer := bufio.NewWriter(file)
	_, err = writer.Write(fileBytes)
	if err != nil {
		log.Printf("文件【%s】保存失败\n", path, err)
		log.Println("保存失败原因：", err)
		return -2, err.Error()
	}

	err = writer.Flush()
	if err != nil {
		log.Println("缓存刷入磁盘失败：", err)
		return -3, err.Error()
	}
	// 写完关闭文件
	defer file.Close()

	return 0, path
}
