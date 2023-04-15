package ziface

type IMessage interface {
	GetMsgId() uint32
	GetMsgData() []byte
	GetDataLen() uint32

	SetDataLen(uint32)
	SetMsgId(uint32)
	SetMsgData([]byte)
}
