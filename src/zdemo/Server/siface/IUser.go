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
	StopConn()

	GetRoom() IRoom

	SetActive(bool)
	IsActive() bool
}
