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
	port       uint32

	router ziface.IRouter
}

func NewServer(name string) (server *Server) {
	server = &Server{
		name:       name,
		ip_version: "tcp4",
		ip:         "127.0.0.1",
		port:       8999,
		router:     nil,
	}
	return
}

func (this *Server) Start() {

	fmt.Printf("start server %s at (%s: %d)\n", this.name, this.ip, this.port)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.ip, this.port))
	if err != nil {
		panic(err.Error())
	}

	var conn_id uint32 = 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		connection := NewConnection(conn_id, conn, this.router)
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

func (this *Server) AddRounter(router ziface.IRouter) {
	this.router = router
}
