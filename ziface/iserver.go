package ziface

type IServer interface {
	//启动服务
	Start()
	//暴露到外部的服务
	Server()
	//终止服务
	Stop()
	//路由功能 给当前的服务注册一个路由方法 供客户端链接使用
	AddRouter(msgID uint32, router IRouter)
	//获取当前的server链接管理器
	GetConnMgr() IConnManager
	//注册OnConnStart钩子函数的方法
	SetOnConnStart(func(connection IConnection))
	//注册OnConnStop钩子函数的方法
	SetOnConnStop(func(connection IConnection))
	//嗲用OnConnStart钩子函数的方法
	CallOnConnStart(connection IConnection)
	//调用OnConnStop钩子函数的方法
	CallOnConnStop(connection IConnection)
}
