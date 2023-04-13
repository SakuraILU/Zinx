package ziface

import "net"

type IRequest interface {
	GetConn() net.Conn
	GetData() []byte
}
