package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/yzhaoyu/gs/gnet"
)

/*
	模拟客户端
*/

func main() {
	fmt.Println("client1 start...")

	time.Sleep(1 * time.Second)

	// 1.直接链接远程服务器得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client starts err, exit!")
		return
	}

	for {
		// 发送封包的message消息, MsgID: 0
		dp := gnet.NewDataPack()
		binaryMsg, _ := dp.Pack(gnet.NewMsgPackage(1, []byte("gsV0.10 client1 Test Message")))
		_, err := conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("write error", err)
			return
		}

		// 服务器就应该给我们回复一个message数据，MsgID：1 pingpingping
		// 1先读取流中的head部分，得到ID和dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error ", err)
			break
		}

		// 将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error ", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			// 2再根据DataLen进行第二次读取，将data读出来
			msg := msgHead.(*gnet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error, ", err)
				return
			}

			fmt.Println("----> Recv Server Msg: ID = ", msg.Id, ", len = ", msg.DataLen, ", data = ", string(msg.Data))
		}

		// CPU阻塞
		time.Sleep(1 * time.Second)
	}
}
