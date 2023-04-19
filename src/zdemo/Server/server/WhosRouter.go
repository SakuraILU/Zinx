package server

import (
	"main/src/zdemo/Server/siface"
	"main/src/zinx/ziface"
	"main/src/zinx/znet"
)

type WhosRouter struct {
	znet.BaseRounter
}

func NewWhoRouter() (who_rt *WhosRouter) {
	who_rt = &WhosRouter{}
	return
}

func (this *WhosRouter) Handle(request ziface.IRequest) {
	iuser, err := request.GetConn().GetProperty("user")
	if err != nil {
		panic(err.Error())
	}
	user := iuser.(siface.IUser)
	room := user.GetRoom()

	names := room.GetUserAll()
	request.GetConn().SendMsg(0, []byte(names))
}
