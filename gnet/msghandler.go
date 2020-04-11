package gnet

/*
	消息处理模块的实现
*/
import (
	"fmt"
	"strconv"

	"github.com/yzhaoyu/giface"
)

type MsgHandle struct {
	// 存放每个MsgID所对应得处理方法
	Apis map[uint32]giface.IRouter
}

// 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]giface.IRouter),
	}
}

// 调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request giface.IRequest) {
	// 1.从Request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is not found! Need Register!")
	}

	// 2.根据MsgID调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router giface.IRouter) {
	// 1.判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		// id已经注册了
		panic("repeat api, msgID = " + strconv.Itoa(int(msgId)))
	}
	// 2.添加msg与API的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api MsgID = ", msgId, " success!")
}
