package znet

import "net"

type Request struct {
	conn net.Conn
	data []byte
}

func NewRequest(conn net.Conn, data []byte) (request *Request) {
	request = &Request{
		conn: conn,
		data: data,
	}
	return
}

func (this *Request) GetConn() (conn net.Conn) {
	conn = this.conn
	return
}

func (this *Request) GetData() (data []byte) {
	data = this.data
	return
}
