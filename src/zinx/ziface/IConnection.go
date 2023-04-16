package ziface

import (
	"net"
)

type IConnection interface {
	Start()
	Stop()
	GetConnID() uint32
	RemoteAddr() net.Addr
	SendMsg(id uint32, data []byte) error

	GetRouterManager() IRouterManager
}
