package ziface

import "net"

type IConnection interface {
	//开始当前链接的服务
	Start()
	//结束当前链接的服务
	Stop()
	//获取当前链接绑定的socket conn
	GetTCPConnection() *net.TCPConn
	//获取当前链接模块的id
	GetConnID() uint32
	//获取远程客户端的tcp状态  ip 端口号
	RemoteAddr() net.Addr
	//发送数据 将数据发给远程
	SendMsg(msgId uint32, data []byte) error
	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string) (interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}

//定义一个处理绑定链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
