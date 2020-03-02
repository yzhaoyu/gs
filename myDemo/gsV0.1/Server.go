package main

import (
	"github.com/yzhaoyu/gs/gnet"
)

/*
	基于gs框架来开发的服务器端应用程序
*/
func main() {
	// 1.创建一个server句柄，使用gs的API
	s := gnet.NewServer("[gs V0.2]")
	// 2.启动server
	s.Serve()
}
