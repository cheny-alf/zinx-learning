package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

func main() {
	fmt.Println("client start")
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		logrus.Error("client start err")
		return
	}

	for {
		_, err := conn.Write([]byte("hello\n"))
		if err != nil {
			logrus.Error("write conn err")
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("read buf err")
			return
		}

		fmt.Printf("server call back:%s,cnt = %d\n", buf, cnt)

		time.Sleep(1 * time.Second)

	}

}
