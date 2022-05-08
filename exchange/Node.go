package exchange

import (
	"encoding/xml"
	"log"
)

/**
  该文件定义相关系统节点信息
 */

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
func initNodes(path string) (*Nodes, int, string) {

	// 解析文件内容
	nodeBytes, errCode, msg := readFile(path)
	if errCode != 0 {
		return nil, errCode, msg
	}
	var nodes Nodes
	err := xml.Unmarshal(nodeBytes, &nodes)
	if err != nil {
		log.Println("字符串不能解析为对应类：", string(nodeBytes))
		return nil, -1, err.Error()
	}

	log.Println("----------节点信息初始化完成----------")

	return &nodes, 0, ""
}