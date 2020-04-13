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

// Test Handle
func (this *PingRouter) Handle(request giface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping..ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	gnet.BaseRouter
}

func (this *HelloRouter) Handle(request giface.IRequest) {
	fmt.Println("Call Hello gs Router Handle...")
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hello Welcome to gs"))
	if err != nil {
		fmt.Println(err)
	}
}

// 创建链接之后执行钩子函数
func DoConnectionBegin(conn giface.IConnection) {
	fmt.Println("===> DoConnectionBegin is called...")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}
}

// 链接断开之前需要执行的函数
func DoConnectionLost(conn giface.IConnection) {
	fmt.Println("===> DoConnectionLost is Called...")
	fmt.Println("connID = ", conn.GetConnID(), " is lost...")
}

func main() {
	// 1.创建一个server句柄，使用gs的API
	s := gnet.NewServer("[gs V0.9]")

	// 2. 注册链接Hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 3.给当前gs框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	// 4.启动server
	s.Serve()
}
