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
	err := request.GetConnection().SendMsg(1, []byte("ping...ping..ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 1.创建一个server句柄，使用gs的API
	s := gnet.NewServer("[gs V0.5]")

	// 2.给当前gs框架添加一个自定义的router
	s.AddRouter(&PingRouter{})

	// 3.启动server
	s.Serve()
}
