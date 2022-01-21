package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("handle")
	fmt.Println("recv from client :msgId:", request.GetMsgID(),
		",data :", request.GetData())
	err := request.GetConnect().SendMsg(200, []byte("ping .....ping ..."))
	if err != nil {
		fmt.Println(err)
	}
}

//ping test 自定义路由
type HelloRouter struct {
	znet.BaseRouter
}

//Handle
func (this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("HelloRouter handle")
	fmt.Println("recv from client :msgId:", request.GetMsgID(),
		",data :", request.GetData())
	err := request.GetConnect().SendMsg(201, []byte("this is hello router"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	server := znet.NewServer("[ZinxV0.6]")
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloRouter{})
	server.Server()

}
