package znet

import (
	"bytes"
	"encoding/binary"
	"main/src/zinx/ziface"
)

type DataPack struct {
}

func NewDataPack() (data_pack *DataPack) {
	data_pack = &DataPack{}
	return
}

func (this *DataPack) Pack(msg ziface.IMessage) (data []byte, err error) {
	// 打印msg信息
	buf := bytes.NewBuffer([]byte{})
	if err = binary.Write(buf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, msg.GetMsgData()); err != nil {
		return
	}

	data = buf.Bytes()
	// assert len
	if binary.LittleEndian.Uint32(data[0:4]) != msg.GetDataLen() {
		panic("Pack error: len is not correct")
	}
	// assert id
	if binary.LittleEndian.Uint32(data[4:8]) != msg.GetMsgID() {
		panic("Pack error: id is not correct")
	}
	// assert data len should be 8 + msg data len
	if len(data) != 8+int(msg.GetDataLen()) {
		panic("Pack error: data len is not correct")
	}
	return
}

func (this *DataPack) UnpackHead(head []byte) (msg ziface.IMessage, err error) {
	buf := bytes.NewBuffer(head)
	var len uint32 = 0
	var id uint32 = 0
	if err = binary.Read(buf, binary.LittleEndian, &len); err != nil {
		return
	}
	if err = binary.Read(buf, binary.LittleEndian, &id); err != nil {
		return
	}
	msg = NewMessage(id, make([]byte, len))
	return
}

func (this *DataPack) GetHeadLen() uint32 {
	return 8
}
