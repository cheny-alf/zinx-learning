package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"time"
	"zinx/znet"
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
		//发送封包msg消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1, []byte("ZinxV0.7 client test message")))
		if err != nil {
			fmt.Println("pack error", err)
			return
		}
		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("write error", err)
			return
		}

		binaryHead := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, binaryHead)
		if err != nil {
			fmt.Println("read file error", err)
			break
		}
		//将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unPack message err:", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("read msg data error")
				return
			}
			fmt.Println("----》Recv Server Msg:id ", msg.MsgId, ",len = ", msg.DataLen, ",data = ", string(msg.Data))

		}

		time.Sleep(1 * time.Second)

	}

}
