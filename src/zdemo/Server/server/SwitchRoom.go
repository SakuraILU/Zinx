package server

import (
	"fmt"
	"main/src/zdemo/Server/siface"
	"main/src/zdemo/Server/utils"
	"main/src/zinx/ziface"
)

type SwitchRoomRouter struct {
	Router
	world siface.IWorld
}

func NewSwtichRoomRouter(world siface.IWorld) (switch_room_rt *SwitchRoomRouter) {
	switch_room_rt = &SwitchRoomRouter{
		world: world,
	}
	return
}

func (this *SwitchRoomRouter) Handle(request ziface.IRequest) {
	conn := request.GetConn()
	nroom_to_buf := request.GetData()

	iuser, err := conn.GetProperty("user")
	if err != nil {
		panic(err.Error())
	}
	user := iuser.(siface.IUser)
	room := user.GetRoom()

	fmt.Printf("swtch from %s to %s\n", room.GetName(), nroom_to_buf)
	if err = this.world.UserSwitchRoom(user, string(nroom_to_buf), room.GetName()); err != nil {
		conn.SendMsg(utils.NErr, []byte(err.Error()))
		return
	}

	str := fmt.Sprintf("enter new room %s", nroom_to_buf)
	conn.SendMsg(utils.NCmdResponse, []byte(str))
}
