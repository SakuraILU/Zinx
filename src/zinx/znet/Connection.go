package znet

import (
	"fmt"
	"io"
	"main/src/zinx/ziface"
	"net"
)

type Connection struct {
	id        uint
	conn      net.Conn
	handler   ziface.Handler
	is_closed bool
}

func NewConnection(id uint, conn net.Conn, handler ziface.Handler) (connection *Connection) {
	connection = &Connection{
		id:        id,
		conn:      conn,
		handler:   handler,
		is_closed: false,
	}
	return
}

func (this *Connection) Start() {
	for {
		buf := make([]byte, 512)
		_, err := this.conn.Read(buf)
		if err == io.EOF {
			fmt.Printf("Connection %d is ended by client\n", this.id)
			break
		} else if err != nil {
			continue
		}
		err = this.handler(NewRequest(this.conn, buf))
		if err != nil {
			panic(err.Error())
		}
	}

	defer this.Stop()
}

func (this *Connection) Stop() {
	if this.is_closed {
		return
	}

	this.conn.Close()
	this.is_closed = true
	fmt.Printf("Connection %d is stopped\n", this.id)
}

func (this *Connection) GetConnID() (id uint) {
	id = this.id
	return
}

func (this *Connection) RemoteAddr() (addr net.Addr) {
	addr = this.conn.RemoteAddr()
	return
}

func (this *Connection) Send(data []byte) error {
	panic("not implemented")
}
