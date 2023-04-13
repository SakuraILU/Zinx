package znet

import (
	"fmt"
	"main/src/zinx/ziface"
	"net"
)

// implement of IServer interface, the Server mode
type Server struct {
	name string

	ip_version string
	ip         string
	port       uint
}

func NewServer(name string) (server *Server) {
	server = &Server{
		name:       name,
		ip_version: "tcp4",
		ip:         "127.0.0.1",
		port:       8999,
	}
	return
}

func handler(request ziface.IRequest) (err error) {
	request.GetConn().Write(request.GetData())
	return
}

func (this *Server) Start() {

	fmt.Printf("start server %s at (%s: %d)\n", this.name, this.ip, this.port)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.ip, this.port))
	if err != nil {
		panic(err.Error())
	}

	var conn_id uint
	conn_id = 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		connection := NewConnection(conn_id, conn, handler)
		go connection.Start()

		conn_id++
	}
}

func (this *Server) Stop() {
	panic("not implemented!")
}

func (this *Server) Serve() {
	this.Start()

	defer this.Stop()
}
