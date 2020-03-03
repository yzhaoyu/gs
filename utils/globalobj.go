package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/yzhaoyu/gs/giface"
)

/*
	存储一切有关gs框架的全局参数，供其他模块使用
	一些参数是可以通过gs.json由用户进行配置
*/

type GlobalObj struct {
	// Server
	TcpServer giface.IServer // 当前gs全局的server对象
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器主机监听的端口号
	Name      string         //当前服务器的名称

	// gs
	Version        string // 当前gs版本号
	MaxConn        int    // 当前服务器主机允许的最大链接数
	MaxPackageSize uint32 // 当前gs框架数据包的最大轴
}

/*
	定义一个全局的对外GlobalObj
*/
var GlobalObject *GlobalObj

/*
	从gs.json去加载用于自定义的参数
*/
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/gs.json")
	if err != nil {
		panic(err)
	}

	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
	提供一个init方法，初始化当前的GlobalObject
*/
func init() {
	// 如果配置文件没有加载。默认的值
	GlobalObject = &GlobalObj{
		Name:           "gsServerApp",
		Version:        "V0.4",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	// 应该尝试从conf/gs.json去加载一些用户自定义的参数
	GlobalObject.Reload()
}
