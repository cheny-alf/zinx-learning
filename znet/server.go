package znet

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}

func (s *Server) Start() {
	logrus.Infof("[Start] Server Listener at IP :%s, Port %d, is starting\n", s.IP, s.Port)

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

		//阻塞的等待客户链接，处理客户端链接的业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				logrus.Errorf("Accept Error: %s", err)
				continue
			}
			//与客户端建立链接 做一些业务 回显
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						logrus.Errorf("rece buf err: %s", err)
						continue
					}
					//回显
					if _, err := conn.Write(buf[:cnt]); err != nil {
						logrus.Errorf("write back buf err: %s", err)
						continue
					}
				}
			}()
		}
	}()

}
func (s *Server) Server() {
	s.Start()
	select {}
}
func (s *Server) Stop() {

}
