package main

import (
	"main/src/zdemo/Server/server"
	"main/src/zdemo/Server/siface"
	"main/src/zdemo/Server/utils"
	"main/src/zinx/znet"
)

// var room siface.IRoom = server.NewRoom("room0", 100)
var world siface.IWorld = server.NewWorld()

func main() {

	chat_server := znet.NewServer()

	// chat
	chat_server.AddRounter(utils.NBroadcast, server.NewBroadcastRouter())
	chat_server.AddRounter(utils.NPrivateChat, server.NewPrivateChatRouter())
	// user
	chat_server.AddRounter(utils.NChangeName, server.NewChangeNameRouter())
	chat_server.AddRounter(utils.NWhos, server.NewWhosRouter())
	// room
	chat_server.AddRounter(utils.NNewRoom, server.NewNewRoomRouter(world))
	chat_server.AddRounter(utils.NSwitchRoom, server.NewSwtichRoomRouter(world))
	chat_server.AddRounter(utils.NCurrentRoom, server.NewCurrentRouter())
	chat_server.AddRounter(utils.NRooms, server.NewRoomsRouter(world))

	chat_server.SetOnConnStart(online)
	chat_server.SetOnConnStop(offline)

	chat_server.Serve()
}
