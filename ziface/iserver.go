package ziface

type IServer interface {
	//启动服务
	Start()
	//暴露到外部的服务
	Server()
	//终止服务
	Stop()
	//路由功能 给当前的服务注册一个路由方法 供客户端链接使用
	AddRouter(router IRouter)
}
