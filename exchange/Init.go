package exchange

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var properties = make(map[string]string)

// 项目对应配置文件路径
var CONF_PATH string

/**
  当前服务端init方法中，..为当前项目路径。
  通过os.Getwd()得到的路径跟client同样的方法得到的不一致：
      因为可以指定了构建之后运行的工作目录，server和client设置得不一样导致。
      idea可以在运行的设置中修改
 */
func init() {

	initPath()
	// 读取配置文件
	readPropertiesScan(CONF_PATH + "/config/application.properties")

	// 初始化节点信息
	initNodes(CONF_PATH + "/config/IB2/ADDR.xml")
	// 赋值当前系统信息
	NODE_SELF = getNode(properties["node"], IB2_NODES.Nodes)
	if NODE_SELF == nil {
		log.Fatalf("节点【%s】信息不存在", properties["node"])
	}
}

/**
  初始化项目路径
 */
func initPath() {
	dirPath, err := os.Getwd()
	log.Println("Main当前代码路径为：", dirPath)
	if err != nil {
		log.Fatalf("获取当前代码路径失败")
	}
	CONF_PATH = dirPath
}

/**
  读取(properties)配置文件
 */
func readPropertiesScan(path string) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("配置打开失败：", path)
	}
	// 设置文件读取缓存
	scanner := bufio.NewScanner(file)
	// 使用scan循环读取文件
	var line string
	var props []string
	for scanner.Scan() {
		line = scanner.Text()
		fmt.Printf("当前行内容为【%s】\n", line)
		props = strings.Split(line, "=")
		properties[props[0]] = props[1]
	}

	// 关闭文件
	defer file.Close()
}
