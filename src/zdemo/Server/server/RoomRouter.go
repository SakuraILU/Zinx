package server

import (
	"fmt"
	"main/src/zdemo/Server/siface"
	"main/src/zdemo/Server/utils"
	"main/src/zinx/ziface"
	"strconv"
	"strings"
)

// Create a new room
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
	if len(strs) == 1 {
		conn.SendMsg(utils.NErr, []byte("usage: newroom [room name] [room capacity]"))
		return
	}

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

// switch into a new room
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

// which room you are in
type CurrentRoomRouter struct {
	Router
}

func NewCurrentRouter() (current_room_rt *CurrentRoomRouter) {
	current_room_rt = &CurrentRoomRouter{}
	return
}

func (this *CurrentRoomRouter) Handle(request ziface.IRequest) {
	conn := request.GetConn()
	iuser, err := request.GetConn().GetProperty("user")
	if err != nil {
		conn.SendMsg(utils.NErr, []byte(err.Error()))
		return
	}
	user := iuser.(siface.IUser)
	room := user.GetRoom()

	conn.SendMsg(utils.NCurrentRoom, []byte(room.GetName()))
}

// find all the rooms (name)
type RoomsRouter struct {
	Router
	world siface.IWorld
}

func NewRoomsRouter(world siface.IWorld) (rooms_rt *RoomsRouter) {
	rooms_rt = &RoomsRouter{
		world: world,
	}
	return
}

func (this *RoomsRouter) Handle(request ziface.IRequest) {
	conn := request.GetConn()

	msg := this.world.GetAllRoomMsg()
	conn.SendMsg(utils.NRooms, []byte(msg))
}
