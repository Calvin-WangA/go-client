package exchange

import (
	"encoding/xml"
	"log"
	"strconv"
)

/**
  该文件定义相关系统节点信息
 */

/** 系统级节点信息 */
var IB2_NODES Nodes
/** 系统自己节点信息 */
var NODE_SELF *Node

/**
  地址对象
 */
type Address struct {
	// 主机名
	Host string `xml:"host"`
	// 端口
	Port int `xml:"port"`
}

/**
  映射为地址数组
 */
type Addresses struct {
	Addresses []Address `xml:"address"`
}

/**
  加解密算法对象
 */
type Encryption struct {
	// 使用的算法名称
	Name string `xml:"name"`
	// 公钥key或者文件地址
	PubKey string `xml:"pubKey"`
	// 私钥key或者对应地址
	PriKey string `xml:"priKey"`
	// key使用标志 0 key值，1文件路径
	Flag bool `xml:"flag"`
}

/**
  节点对象信息
 */
type Node struct {
	// 系统名称
    Name string `xml:"name"`
    // 系统代码
    Code string `xml:"code"`
    // 支持的编码格式
    Encode string `xml:"encode"`
    // 地址对象集合
    Addresses Addresses `xml:"addresses"`
    // 加密算法信息
    Encryption Encryption `xml:"encryption"`
    // 交互协议
	Protocol string `xml:"protocol"`
}

/**
  节点数组对象封装
 */
type Nodes struct {
	Nodes []Node `xml:"node"`
}

/**
  解析节点配置文件内容
 */
func initNodes(path string) {

	// 解析文件内容
	nodeBytes, exchangeError := readFile(path)
	if exchangeError.IsFail() {
		log.Fatalf("配置文件【%s】解析失败代码>>>>>>>>>\n", path)
	}
	err := xml.Unmarshal(nodeBytes, &IB2_NODES)
	if err != nil {
		log.Fatalln("字符串不能解析为对应类：", string(nodeBytes))
	}

	log.Println("----------节点信息初始化完成----------")
}

/**
  获取单个节点
 */
func getNode(nodeCode string, nodes []Node) *Node {
	for _, node := range nodes {
		if nodeCode == node.Code {
			return &node
		}
	}

	return nil
}

// 可进行负载算法处理
func getAddress (node Node) string {
	for _, address := range node.Addresses.Addresses {
		if address.Port != -1 {
			return address.Host + ":" + strconv.Itoa(address.Port)
		}
	}

	return ""
}
