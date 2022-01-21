package znet

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	//当前server的消息管理模块，用来绑定msgId和对应的处理业务api关系
	MsgHandle ziface.IMsgHandler
	//该server的链接管理器
	ConnMgr ziface.IConnManager
	//该server创建链接之后自动调用Hook函数 onConnStart
	OnConnStart func(conn ziface.IConnection)
	//该server创建链接之前自动调用Hook函数 onConnStart
	OnConnStop func(conn ziface.IConnection)
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandle(),
		ConnMgr:   NewConnManager(),
	}
	return s
}

func (s *Server) Start() {
	logrus.Infof("[Zinx] Server Listener at IP :%s, Port %d, is starting", utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

	go func() {
		//开启一个消息队列 worker工作池
		s.MsgHandle.StartWorkPool()

		//获取tcp的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			logrus.Errorf("resolve tcp addr error:%s", err)
		}

		//监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			logrus.Errorf("listen %s error: %s", s.IPVersion, err)
			return
		}
		logrus.Infof("start Zinx server success, %s listening successful ", s.Name)

		var cid uint32
		//阻塞的等待客户链接，处理客户端链接的业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				logrus.Errorf("Accept Error: %s", err)
				continue
			}

			//设置最大链接个数的判断，如果超过最大连接数 那么关闭这个新的链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端相应一个超出最大连接数的数据包
				logrus.Infof("============================Too many Connection! MaxConnection is %s", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			//将处理新连接的业务和conn进行binding 得到需要的模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandle)
			cid++
			//启动当前的链接业务处理
			go dealConn.Start()
		}
	}()

}
func (s *Server) Server() {
	s.Start()
	select {}
}
func (s *Server) Stop() {
	//将一些服务器的资源 状态 或者已经开辟的链接信息 进行停止或者回收
	logrus.Info("[STOP] Zinx server name", s.Name)
	s.ConnMgr.ClearConn()
}
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgID, router)
	logrus.Info("Add Router Success")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

//注册OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		logrus.Info("Call OnConnStart().....")
		s.OnConnStart(conn)
	}
}

//调用OnConnStop钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		logrus.Info("Call OnConnStop().....")
		s.OnConnStop(conn)
	}
}
