package znet

type Message struct {
	MsgId   uint32
	DataLen uint32
	Data    []byte
}

func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		MsgId:   id,
		DataLen: uint32(len(data)),
		Data:    data,
	}

}
func (m *Message) GetMsgId() uint32 {
	return m.MsgId
}
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}
func (m *Message) GetData() []byte {
	return m.Data
}
func (m *Message) SetMsgId(id uint32) {
	m.MsgId = id
}
func (m *Message) SetMsgLen(len uint32) {
	m.DataLen = len
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}
