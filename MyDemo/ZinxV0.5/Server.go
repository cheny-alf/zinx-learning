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
	err := request.GetConnect().SendMsg(1, []byte("ping .....ping ..."))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	server := znet.NewServer("[ZinxV0.5]")
	server.AddRouter(&PingRouter{})
	server.Server()

}
