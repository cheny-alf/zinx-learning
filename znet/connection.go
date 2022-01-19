package znet

import (
	"github.com/sirupsen/logrus"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	//当前链接的socket tcp套接字
	Conn *net.TCPConn
	//链接的ID
	ConnID uint32
	//当前的链接状态
	isClosed bool
	//告知当前链接已经退出停止的channel
	ExitChan chan bool
	//该链接处理的方法的router
	Router ziface.IRouter
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
	return c
}

//开始当前链接的服务
func (c *Connection) Start() {
	logrus.Info("Conn Start()....ConnID = %s", c.ConnID)
	go c.StartReader()
	//Todo
}

//获取当前链接绑定的socket conn
func (c *Connection) Stop() {
	logrus.Info("Conn Stop()....ConnID:%s", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	c.Conn.Close()
	close(c.ExitChan)
}

//获取当前链接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的tcp状态  ip 端口号
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据 将数据发给远程
func (c *Connection) Send(data []byte) error {
	return nil
}

func (c *Connection) StartReader() {
	logrus.Info("Reader Goroutine is running....")
	defer logrus.Info("ConnID = ", c.ConnID, " Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中 最大512
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			logrus.Errorf("recv buf err: %s", err)
			continue
		}

		//得到当前conn数据的request请求数据
		req := Request{
			conn: c,
			data: buf,
		}
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}

}
