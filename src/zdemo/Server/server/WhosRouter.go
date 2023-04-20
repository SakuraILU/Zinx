package server

import (
	"main/src/zdemo/Server/siface"
	"main/src/zinx/ziface"
)

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
