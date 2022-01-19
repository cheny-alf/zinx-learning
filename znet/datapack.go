package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	//Datalen uint32(4字节) + ID uint32（4字节）
	return 8
}

/*
	1.创建一个存放byte字节的缓冲
	2.将dataLen写入buf中
	3.将msgId写入buf中
	4.将data数据 写入buf中
*/
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuf := bytes.NewBuffer([]byte{})
	//写dataLen
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	//写msgID
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuf.Bytes(), nil
}

//解开
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	dataBuf := bytes.NewBuffer(binaryData)
	msg := &Message{}
	//解压head信息得到len和id
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.MsgId); err != nil {
		return nil, err
	}
	if (utils.GlobalObject.MaxPackageSize > 0) && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}
	return msg, nil
}
