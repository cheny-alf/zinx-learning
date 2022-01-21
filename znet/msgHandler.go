package znet

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	//存放每个msgid对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作worker池的worker数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

//启动一个worker工作池（开启的动作只会发生一次，一个框架 只有一个池子）
func (mh *MsgHandle) StartWorkPool() {
	//根据workerPoolSize 分别开启worker 每个worker用一个go协程来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//根据当前的worker启动对应的消息队列 开辟空间 第0个worker就启动第0个队列
		//启动当前的worker 阻塞等待消息从channel中过来
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	logrus.Infof("Worker ID = %d is started...", workerID)
	//不断阻塞等待对应消息队列的消息
	for {
		select {
		//如果有消息过来 出列的就是一个客户端的request 执行当前request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将消息交给TaskQueue 由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//将消息平均分配给不同的worker
	//根据客户端建立的connID来进行分配
	workerID := request.GetConnect().GetConnID() % mh.WorkerPoolSize
	logrus.Infof("Add ConnID = %d request MsgID = %d to WorkerID = %d", request.GetConnect().GetConnID(), request.GetMsgID(), workerID)
	//将消息发送给对应的worker的taskQueue即可
	mh.TaskQueue[workerID] <- request

}
