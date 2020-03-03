package main

import (
	"fmt"

	"github.com/yzhaoyu/gs/giface"
	"github.com/yzhaoyu/gs/gnet"
)

/*
	基于gs框架来开发的服务器端应用程序
*/

// ping test自定义路由
type PingRouter struct {
	gnet.BaseRouter
}

// Test PreHandle
func (this *PingRouter) PreHandle(request giface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Test Handle
func (this *PingRouter) Handle(request giface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping...ping...ping error")
	}
}

// Test PostHandle
func (this *PingRouter) PostHandle(request giface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	// 1.创建一个server句柄，使用gs的API
	s := gnet.NewServer("[gs V0.3]")

	// 2.给当前gs框架添加一个自定义的router
	s.AddRouter(&PingRouter{})

	// 3.启动server
	s.Serve()
}
