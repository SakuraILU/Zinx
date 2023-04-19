package server

import (
	"fmt"
	"main/src/zdemo/Server/siface"
	"main/src/zinx/ziface"
	"main/src/zinx/znet"
)

type ChangeNameRouter struct {
	znet.BaseRounter
}

func NewChangeNameRouter() (change_name_rt *ChangeNameRouter) {
	change_name_rt = &ChangeNameRouter{}
	return
}

func (this *ChangeNameRouter) Handle(request ziface.IRequest) {
	iuser, err := request.GetConn().GetProperty("user")
	if err != nil {
		panic(err.Error())
	}
	user := iuser.(siface.IUser)
	room := user.GetRoom()

	_, err = room.GetUser(string(request.GetData())) // check existance
	if err == nil {
		request.GetConn().SendMsg(0, []byte("this name already exist"))
		return
	}

	room.RemoveUser(user)

	user.SetName(string(request.GetData()))
	err = room.AddUser(user)
	if err != nil {
		panic(err.Error())
	}

	msg := fmt.Sprintf("new name: %s", user.GetName())
	request.GetConn().SendMsg(0, []byte(msg))
}
