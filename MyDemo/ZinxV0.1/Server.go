package main

import "zinx/znet"

func main() {
	server := znet.NewServer("[ZinxV0.1]")
	server.Server()
}
