package znet

import (
	"fmt"
	"io"
	"main/src/zinx/ziface"
	"net"
)

type Connection struct {
	id        uint32
	conn      net.Conn
	router    ziface.IRouter
	is_closed bool

	data_pack ziface.IDataPack
}

func NewConnection(id uint32, conn net.Conn, router ziface.IRouter) (connection *Connection) {
	data_pack := NewDataPack()

	connection = &Connection{
		id:        id,
		conn:      conn,
		router:    router,
		is_closed: false,

		data_pack: data_pack,
	}
	return
}

func (this *Connection) unpackMsg() (msg ziface.IMessage, err error) {

	buf := make([]byte, this.data_pack.GetHeadLen())
	_, err = this.conn.Read(buf)

	if err != nil {
		return
	}
	msg, err = this.data_pack.UnpackHead(buf)
	if err != nil {
		return
	}

	_, err = this.conn.Read(msg.GetMsgData())
	// _, err = io.ReadFull(this.conn, msg.GetMsgData())
	if err != nil {
		return
	}

	return
}

func (this *Connection) Start() {
	for {
		msg, err := this.unpackMsg()
		if err == io.EOF {
			fmt.Printf("Connection %d is ended by client\n", this.id)
			break
		} else if err != nil {
			continue
		}

		request := NewRequest(this, msg)

		this.router.PreHandle(request)
		this.router.Handle(request)
		this.router.PostHandle(request)
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

func (this *Connection) GetConnID() (id uint32) {
	id = this.id
	return
}

func (this *Connection) RemoteAddr() (addr net.Addr) {
	addr = this.conn.RemoteAddr()
	return
}

func (this *Connection) SendMsg(id uint32, data []byte) (err error) {
	msg := NewMessage(id, data)
	buf, err := this.data_pack.Pack(msg)
	if err != nil {
		return
	}
	this.conn.Write(buf)
	return
}
