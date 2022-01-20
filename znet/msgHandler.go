package znet

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"zinx/ziface"
)

type MsgHandle struct {
	//存放每个msgid对应的处理方法
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

//调度对应的router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1。从request中找到msgID
	//2.根据MsgID 调度对应的router业务
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		logrus.Error("aps msgID = ", request.GetMsgID(), "is NOT FOUND! Need Register!")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理router
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//首先判断msg绑定的api处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api ,msgId:" + strconv.Itoa(int(msgID)))
	}
	//添加msg于api的绑定关系
	mh.Apis[msgID] = router
	logrus.Info("Add api MsgID = ", msgID, "succ!")
}
