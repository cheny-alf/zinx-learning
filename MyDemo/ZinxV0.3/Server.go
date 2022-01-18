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

//test PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("pre handle")
	_, err := request.GetConnect().GetTCPConnection().Write([]byte("before ping"))
	if err != nil {
		fmt.Println("before callback err")
	}
}

//Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("handle")
	_, err := request.GetConnect().GetTCPConnection().Write([]byte("ping>>>>>>>>>ping"))
	if err != nil {
		fmt.Println("callback err")
	}
}

//PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("after handle")
	_, err := request.GetConnect().GetTCPConnection().Write([]byte("after ping"))
	if err != nil {
		fmt.Println("after callback err")
	}
}
func main() {
	server := znet.NewServer("[ZinxV0.3]")
	server.AddRouter(&PingRouter{})
	server.Server()

}
