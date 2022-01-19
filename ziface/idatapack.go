package ziface

type IDataPack interface {
	GetHeadLen()
	Pack(msg IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
}
