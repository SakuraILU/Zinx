package server

import (
	"fmt"
	"main/src/zdemo/Server/siface"
	"main/src/zdemo/Server/utils"
	"main/src/zinx/ziface"
	"strconv"
	"strings"
)

type NewRoomRouter struct {
	Router
	world siface.IWorld
}

func NewNewRoomRouter(world siface.IWorld) (new_room_rt *NewRoomRouter) {
	new_room_rt = &NewRoomRouter{
		world: world,
	}
	return
}

func (this *NewRoomRouter) Handle(request ziface.IRequest) {
	conn := request.GetConn()
	buf := request.GetData()
	strs := strings.SplitN(string(buf), " ", 2)
	name := strs[0]
	cap, err := strconv.Atoi(strs[1])
	if err != nil {
		conn.SendMsg(utils.NErr, []byte(err.Error()))
		return
	}

	err = this.world.CreateRoom(name, uint32(cap))
	if err != nil {
		conn.SendMsg(utils.NErr, []byte(err.Error()))
		return
	}

	str := fmt.Sprintf("create new room %s with capacity %d", name, cap)
	conn.SendMsg(utils.NCmdResponse, []byte(str))
}
