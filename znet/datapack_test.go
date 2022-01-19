package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err", err)
		return
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("accept error")
			}

			go func(conn net.Conn) {
				//处理客户端的请求
				//拆包的过程 如下
				dp := NewDataPack()
				for {
					//第一次从conn读取 把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack error")
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//说明msg中还是有数据的 需要二次读取
						//第二次从conn处开始读取 根据head中的datalen 再读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						//根据datalen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack err:", err)
							return
						}
						fmt.Println("-->receive msgID:", msg.MsgId, ",dataLen", msg.DataLen, "data:", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	//模拟客户端

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dail err:", err)
		return
	}
	//创建一个封包对象
	/*
		模拟封包过程
		1.封装第一个msg1包
		2.封装第二个msg2包
		3.将两个包粘在一起
		4.一次性发送给服务端
	*/
	dp := NewDataPack()
	msg1 := &Message{
		MsgId:   1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error:", err)
		return
	}
	msg2 := &Message{
		MsgId:   2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error")
		return
	}
	sendData = append(sendData, sendData2...)
	conn.Write(sendData)
	select {}
}
