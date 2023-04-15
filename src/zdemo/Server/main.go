package main

import (
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
	err := request.GetConn().SendMsg(0, []byte("PreHandle"))
	if err != nil {
		panic(err.Error())
	}
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	data := request.GetData()
	err := request.GetConn().SendMsg(1, data)
	if err != nil {
		panic(err.Error())
	}
}

func (this *PingRouter) PostHandle(request ziface.IRequest) {
	err := request.GetConn().SendMsg(0, []byte("PostHandle"))
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	server := znet.NewServer("sever_v0.3")
	server.AddRounter(NewPingRouter())
	server.Serve()
}
