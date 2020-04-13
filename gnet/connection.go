package gnet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/yzhaoyu/gs/giface"
	"github.com/yzhaoyu/gs/utils"
)

/*
	链接模块
*/
type Connection struct {
	// 当前Conn隶属于哪个Server
	TcpServer giface.IServer

	// 当前链接的socket TCP套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnID uint32

	// 当前的链接状态
	isClosed bool

	// 当前链接所绑定的处理业务方法API
	handleAPI giface.HandleFunc

	// 告知当前链接已经退出的/停止 channel(由Reader告知Writer退出)
	ExitChan chan bool

	// 无缓冲的管道，用于读、写Goroutine之间的通信
	msgChan chan []byte

	// 消息的管理MsgID和对应的处理业务API关系
	MsgHandler giface.IMsgHandle
}

// 初始化链接模块的方法
func NewConnection(server giface.IServer, conn *net.TCPConn, connID uint32, msgHandler giface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}

	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("[Reader is exiting!] connID = ", c.ConnID, "remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {

		// 创建一个拆包解包对象
		dp := NewDataPack()

		// 读取客户端的Msg Head二进制流8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		// 拆包，得到MsgID和MsgDataLen放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		// 根据DataLen，再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data := make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给Worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中找到注册绑定的Conn对应的router调用
			// 根据绑定好的MsgID找到对应处理api业务
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 写消息Goroutine，专门发送给客户端消息模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!")

	// 不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}

// 启动链接，让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)
	// 启动从当前链接的读数据的业务
	go c.StartReader()

	// 启动从当前链接写数据的业务
	go c.StartWriter()

	// 按照开发者传递进来的 创建链接之后需要调用的处理业务，执行对应Hook函数
	c.TcpServer.CallOnConnStart(c)
}

// 停止链接，结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	// 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 调用开发者注册的 销毁链接之前 需要执行的业务Hook函数
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket链接
	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前链接从ConnMgr中删除掉
	c.TcpServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前链接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 提供一个SendMsg方法，将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将data进行封包 MsgDataLen/MsgID/Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Package error msg Id = ", msgId)
		return errors.New("Pack error msg")
	}

	// 将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}
