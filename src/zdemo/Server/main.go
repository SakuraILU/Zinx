package main

import (
	"fmt"
	"main/src/zdemo/Server/server"
	"main/src/zdemo/Server/siface"
	"main/src/zinx/ziface"
	"main/src/zinx/znet"
)

type OnlineRouter struct {
	znet.BaseRounter
	room siface.IRoom
}

func NewOnlineRouter(room siface.IRoom) (online_rt *OnlineRouter) {
	online_rt = &OnlineRouter{
		room: room,
	}
	return
}

func (this *OnlineRouter) Handle(request ziface.IRequest) {
	user := server.NewUser(string(request.GetData()), request.GetConn(), this.room)
	err := this.room.AddUser(user)
	if err != nil {
		request.GetConn().SendMsg(10, []byte(err.Error()))
		return
	} else {
		go user.StartWorker()
		fmt.Printf("user %s is online...", user.GetName())
		this.room.BroadCastMsg([]byte(fmt.Sprintf("user %s is online...", user.GetName())))
	}
}

type BroadcastRouter struct {
	znet.BaseRounter
	room siface.IRoom
}

func NewBroadcastRouter(room siface.IRoom) (broadcast_rt *BroadcastRouter) {
	broadcast_rt = &BroadcastRouter{
		room: room,
	}

	return
}

func (this *BroadcastRouter) Handle(request ziface.IRequest) {
	this.room.BroadCastMsg(request.GetData())
}

func main() {
	chat_server := znet.NewServer()

	room := server.NewRoom("room0", 100)
	go room.StartRoom()
	chat_server.AddRounter(0, NewOnlineRouter(room))
	chat_server.AddRounter(1, NewBroadcastRouter(room))

	chat_server.SetOnConnStart(func(conn ziface.IConnection) {
		fmt.Printf("Connection %d is established\n", conn.GetConnID())
	})
	chat_server.SetOnConnStop(func(conn ziface.IConnection) {
		fmt.Printf("Connection %d is stopped\n", conn.GetConnID())
	})

	chat_server.Serve()

	room.StopRoom()
}
