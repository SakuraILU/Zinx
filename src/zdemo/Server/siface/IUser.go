package siface

import (
	"net"
)

type IUser interface {
	GetName() string
	SetName(string)

	GetRemoteAddr() net.Addr
	AddMsg([]byte)
	StartWorker()
	StopWorker()

	GetRoom() IRoom
}
