package znet

import "zinx/ziface"

type Request struct {
	//已经和客户端建立好链接的conn
	conn ziface.IConnection
	//客户带你请求的数据
	data []byte
}

//得到当前链接
func (r *Request) GetConnect() ziface.IConnection {
	return r.conn
}

//得到请求的数据
func (r *Request) GetData() []byte {
	return r.data
}
