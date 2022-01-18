package ziface

type IServer interface {
	//启动服务
	Start()
	//暴露到外部的服务
	Server()
	//终止服务
	Stop()
}
