package exchange

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var properties = make(map[string]string)

/**
  文件路径点为当前路径项目路径
 */
func init() {
	// 读取配置文件
	readPropertiesScan("./config/application.properties")

	// 初始化节点信息
	initNodes("./config/IB2/ADDR.xml")
	// 赋值当前系统信息
	NODE_SELF = getNode(properties["node"], IB2_NODES.Nodes)
	if NODE_SELF == nil {
		log.Fatalf("节点【%s】信息不存在", properties["node"])
	}
}

func readPropertiesScan(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("配置打开失败：", err)
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
