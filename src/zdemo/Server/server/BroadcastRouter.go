package server

import (
	"fmt"
	"main/src/zdemo/Server/siface"
	"main/src/zinx/ziface"
)

type BroadcastRouter struct {
	Router
}

func NewBroadcastRouter() (broadcast_rt *BroadcastRouter) {
	broadcast_rt = &BroadcastRouter{}
	return
}

func (this *BroadcastRouter) Handle(request ziface.IRequest) {
	iuser, err := request.GetConn().GetProperty("user")
	if err != nil {
		panic(err.Error())
	}
	user := iuser.(siface.IUser)
	room := user.GetRoom()

	msg := fmt.Sprintf("[%s]:%s", user.GetName(), request.GetData())
	room.BroadCastMsg([]byte(msg))
}
