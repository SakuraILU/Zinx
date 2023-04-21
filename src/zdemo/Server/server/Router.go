package server

import (
	"main/src/zdemo/Server/siface"
	"main/src/zinx/ziface"
	"main/src/zinx/znet"
)

type Router struct {
	znet.BaseRounter
}

// set user to be active if send some msgs which triggers a router
func (this *Router) PreHandle(request ziface.IRequest) {
	iuser, err := request.GetConn().GetProperty("user")
	if err != nil {
		panic(err.Error())
	}
	user := iuser.(siface.IUser)
	user.SetActive(true)
}
