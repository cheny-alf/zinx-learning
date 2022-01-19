package ziface

type IRequest interface {
	//得到当前链接
	GetConnect() IConnection
	//得到请求的数据
	GetData() []byte
	GetMsgID() uint32
}
