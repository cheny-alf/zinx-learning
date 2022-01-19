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
	//当前的server添加一个router service注册的链接对应的处理业务
	Router ziface.IRouter
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}

func (s *Server) Start() {
	logrus.Infof("[Zinx] Server Listener at IP :%s, Port %d, is starting", utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

	go func() {
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

			//将处理新连接的业务和conn进行binding 得到需要的模块
			dealConn := NewConnection(conn, cid, s.Router)
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

}
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	logrus.Info("Add Router Success")
}
