package ziface

import (
	"net"
)

type IConnection interface {
	Start()
	Stop()
	GetConnID() uint
	RemoteAddr() net.Addr
	Send(data []byte) error
}

type Handler func(IRequest) error
