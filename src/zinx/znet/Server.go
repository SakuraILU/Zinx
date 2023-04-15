package znet

import (
	"fmt"
	"main/src/zinx/utils"
	"main/src/zinx/ziface"
	"net"
)

// implement of IServer interface, the Server mode
type Server struct {
	name string

	ip_version string
	ip         string
	port       uint32

	rt_manager ziface.IRouterManager
	work_pool  ziface.IWorkPool
}

func NewServer(name string) (server *Server) {
	server = &Server{
		name:       utils.Global_obj.Name,
		ip_version: "tcp4",
		ip:         utils.Global_obj.Ip,
		port:       utils.Global_obj.Port,
		rt_manager: NewRouterManager(),
		work_pool:  NewWorkPool(),
	}
	return
}

func (this *Server) Start() {
	this.work_pool.StartWorkPool()

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

		connection := NewConnection(conn_id, conn, this.rt_manager, this.work_pool)
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

func (this *Server) AddRounter(msg_id uint32, router ziface.IRouter) {
	this.rt_manager.AddRouter(msg_id, router)
}
