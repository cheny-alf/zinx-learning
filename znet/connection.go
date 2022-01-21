package znet

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	//当前conn是属于那个server
	TCPServer ziface.IServer
	//当前链接的socket tcp套接字
	Conn *net.TCPConn
	//链接的ID
	ConnID uint32
	//当前的链接状态
	isClosed bool
	//告知当前链接已经退出停止的channel
	ExitChan chan bool
	//无缓冲通道 用于读 写Groutine之间的消息通信
	msgChan chan []byte
	//消息的管理msgID 和对应的处理业务api关系
	MsgHandle ziface.IMsgHandler
}

//初始化链接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TCPServer: server,
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		MsgHandle: handler,
		msgChan:   make(chan []byte),
		ExitChan:  make(chan bool, 1),
	}

	//将conn加入connMgr中
	c.TCPServer.GetConnMgr().Add(c)
	return c
}

//开始当前链接的服务
func (c *Connection) Start() {
	logrus.Infof("Conn Start()....ConnID = %d", c.ConnID)
	//启动当前链接的读业务
	go c.StartReader()
	//启动当前链接的写业务
	go c.StartWriter()

	//按照开发者传递进来的 创建链接之后需要调用的处理业务 执行对应hook函数
	c.TCPServer.CallOnConnStart(c)
}

//获取当前链接绑定的socket conn
func (c *Connection) Stop() {
	logrus.Infof("Conn Stop()....ConnID:%d", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true

	//按照开发者传递进来的 创建链接之前需要调用的处理业务 执行对应hook函数
	c.TCPServer.CallOnConnStop(c)
	//关闭socket链接
	c.Conn.Close()
	//告知writer关闭
	c.ExitChan <- true
	//将当前链接从connManager中删除
	c.TCPServer.GetConnMgr().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
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

//提供一个sendMsg方法 将我们要发送给客户端的数据 先进行封包 再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection is closed when send message")
	}

	//将data进行封包
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id:", msgId)
		return errors.New("pack error msg")
	}

	//将数据发给客户端
	c.msgChan <- binaryMsg
	return nil
}

func (c *Connection) StartReader() {
	logrus.Info("[Reader Goroutine is running]")
	defer logrus.Info("ConnID = ", c.ConnID, " [Reader is exit, remote addr is", c.RemoteAddr().String(), "]")
	defer c.Stop()

	for {
		//读取客户端的数据到buf中 最大512
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	logrus.Errorf("recv buf err: %s", err)
		//	continue
		//}
		//
		//创建一个拆包解包对象
		dp := NewDataPack()
		//读取客户端的msg head的二进制流 8个字节
		//拆包 得到msgId 和 msgLen 存入msg中
		//根据len再读取data 存入msg.data
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error:", err)
			break
		}
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)
		//得到当前conn数据的request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制 将消息发送给worker工作池处理即可
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			//从路由中，找到的注册绑定的conn对应的router调用
			//根据绑定好的msgID 找到对应处理api业务 执行
			go c.MsgHandle.DoMsgHandler(&req)
		}
	}

}

//写消息goroutine 专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	logrus.Info("[Write Goroutine is running]")
	defer logrus.Infof("%s [Conn Write Exit!]", c.RemoteAddr().String())
	//不断的阻塞等待channel消息 进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				logrus.Error("Send data error:", err)
				return
			}
		case <-c.ExitChan:
			//代表从退出管道中读到消息 此时reader已经退出 writer也要退出
			return
		}
	}
}
