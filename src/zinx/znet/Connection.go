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
	is_closed bool

	server       ziface.IServer
	rt_manager   ziface.IRouterManager
	work_pool    ziface.IWorkPool
	conn_manager ziface.IConnectionManager

	data_pack ziface.IDataPack
	msg_chan  chan []byte
	exit_chan chan bool
}

func NewConnection(id uint32, conn net.Conn, server ziface.IServer) (connection *Connection) {
	data_pack := NewDataPack()

	connection = &Connection{
		id:        id,
		conn:      conn,
		is_closed: false,

		server:       server,
		rt_manager:   server.GetRouterManager(),
		work_pool:    server.GetWorkPool(),
		conn_manager: server.GetConnectionManager(),

		data_pack: data_pack,
		msg_chan:  make(chan []byte, 3),
		exit_chan: make(chan bool),
	}

	fmt.Printf("current has %d connections\n", server.GetConnectionManager().Size())
	server.GetConnectionManager().Add(connection)
	return
}

func (this *Connection) unpackMsg() (msg ziface.IMessage, err error) {

	buf := make([]byte, this.data_pack.GetHeadLen())

	// // ReadFull reads exactly len(buf) bytes from r into buf...避免了TCP包的帧头段可能不完整的问题，more safer
	_, err = io.ReadFull(this.conn, buf) // _, err = this.conn.Read(buf)

	if err != nil {
		return
	}
	msg, err = this.data_pack.UnpackHead(buf)
	if err != nil {
		return
	}

	// ReadFull reads exactly len(buf) bytes from r into buf...避免了TCP包的数据段可能不完整的问题, more safer
	_, err = io.ReadFull(this.conn, msg.GetMsgData()) // this.conn.Read(msg.GetMsgData())
	// _, err = io.ReadFull(this.conn, msg.GetMsgData())
	if err != nil {
		return
	}

	return
}

func (this *Connection) Writer() {
end:
	for {
		select {
		case <-this.exit_chan:
			break end
		case buf := <-this.msg_chan:
			if _, err := this.conn.Write(buf); err != nil {
				break end
			}
		}
	}

	// fmt.Printf("writer of connection %d is exit...", this.id)
}

func (this *Connection) Start() {
	defer this.Stop()

	// start goes out a Writer，then serves as Reader itself
	go this.Writer()

	for {
		msg, err := this.unpackMsg()
		if err == io.EOF {
			// fmt.Printf("Connection %d is ended by client\n", this.id)
			break
		} else if err != nil {
			continue
		}

		request := NewRequest(this, msg)

		// 读写分离后，读不需要handler的数据了，因此无需等待handler，把handler给go出去
		go this.work_pool.AddRequest(request)
	}

}

func (this *Connection) Stop() {
	if this.is_closed {
		return
	}
	fmt.Printf("Connection %d is stopped[Reader]\n", this.id)

	// tell writer that reader is ready to exit...
	// otherwise writer may still use the closed tcp connection and cause error
	this.exit_chan <- true

	this.conn.Close()
	this.is_closed = true

	close(this.msg_chan)
	close(this.exit_chan)

	this.conn_manager.Remove(this)
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
	// this.conn.Write(buf)
	this.msg_chan <- buf
	return
}

func (this *Connection) GetRouterManager() ziface.IRouterManager {
	return this.rt_manager
}
