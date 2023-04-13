package znet

import (
	"fmt"
	"io"
	"net"
)

// implement of IServer interface, the Server mode
type Server struct {
	name string

	ip_version string
	ip         string
	port       int
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

func (this *Server) Start() {

	fmt.Printf("start server %s at (%s: %d)\n", this.name, this.ip, this.port)
	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", this.ip, this.port))
	if err != nil {
		panic(err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		go func() {
			for {
				buf := make([]byte, 512)
				cnt, err := conn.Read(buf)
				if err == io.EOF {
					break
				} else if err != nil {
					panic(err.Error())
				}

				_, err = conn.Write(buf[:cnt])
				if err != nil {
					panic(err.Error())
				}
			}
		}()

	}
}

func (this *Server) Stop() {

}

func (this *Server) Serve() {
	go this.Start()
	select {}
}
