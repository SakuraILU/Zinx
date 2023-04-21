package server

import (
	"fmt"
	"main/src/zdemo/Server/siface"
	"main/src/zdemo/Server/utils"
	"main/src/zinx/ziface"
)

// Change user name
type ChangeNameRouter struct {
	Router
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
		request.GetConn().SendMsg(utils.NCmdResponse, []byte("This name already exist"))
		return
	}

	room.RemoveUser(user)

	user.SetName(string(request.GetData()))
	err = room.AddUser(user)
	if err != nil {
		panic(err.Error())
	}

	msg := fmt.Sprintf("Set new name to %s", user.GetName())
	request.GetConn().SendMsg(utils.NCmdResponse, []byte(msg))
}

// find all users in this room (names)
type WhosRouter struct {
	Router
}

func NewWhosRouter() (whos_rt *WhosRouter) {
	whos_rt = &WhosRouter{}
	return
}

func (this *WhosRouter) Handle(request ziface.IRequest) {
	iuser, err := request.GetConn().GetProperty("user")
	if err != nil {
		panic(err.Error())
	}
	user := iuser.(siface.IUser)
	room := user.GetRoom()

	names := room.GetAllUserMsg()
	request.GetConn().SendMsg(0, []byte(names))
}
