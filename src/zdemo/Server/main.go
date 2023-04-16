package main

import (
	"fmt"
	"main/src/zinx/ziface"
	"main/src/zinx/znet"
)

type PingRouter struct {
	znet.BaseRounter
}

func NewPingRouter() (ping_test *PingRouter) {
	ping_test = &PingRouter{}
	return
}

func (this *PingRouter) PreHandle(request ziface.IRequest) {
	// err := request.GetConn().SendMsg(0, []byte("PreHandle"))
	// if err != nil {
	// 	panic(err.Error())
	// }
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	err := request.GetConn().SendMsg(0, []byte(fmt.Sprintf("I am Server connetion %d", request.GetConn().GetConnID())))
	if err != nil {
		panic(err.Error())
	}
}

func (this *PingRouter) PostHandle(request ziface.IRequest) {
	// err := request.GetConn().SendMsg(0, []byte("PostHandle"))
	// if err != nil {
	// 	panic(err.Error())
	// }
}

type EchoRouter struct {
	znet.BaseRounter
}

func NewEchoRouter() (ping_test *EchoRouter) {
	ping_test = &EchoRouter{}
	return
}

func (this *EchoRouter) PreHandle(request ziface.IRequest) {
	// err := request.GetConn().SendMsg(0, []byte("PreHandle"))
	// if err != nil {
	// 	panic(err.Error())
	// }
}

func (this *EchoRouter) Handle(request ziface.IRequest) {
	data := request.GetData()
	err := request.GetConn().SendMsg(0, data)
	if err != nil {
		panic(err.Error())
	}
}

func (this *EchoRouter) PostHandle(request ziface.IRequest) {
	// err := request.GetConn().SendMsg(0, []byte("PostHandle"))
	// if err != nil {
	// 	panic(err.Error())
	// }
}

func main() {
	server := znet.NewServer()

	server.AddRounter(0, NewEchoRouter())
	server.AddRounter(1, NewPingRouter())

	server.SetOnConnStart(func(conn ziface.IConnection) {
		fmt.Printf("Connection %d is established\n", conn.GetConnID())
	})
	server.SetOnConnStop(func(conn ziface.IConnection) {
		fmt.Printf("Connection %d is stopped\n", conn.GetConnID())
	})

	server.Serve()
}
