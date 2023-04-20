package main

import (
	"main/src/zdemo/Server/server"
	"main/src/zdemo/Server/siface"
	"main/src/zdemo/Server/utils"
	"main/src/zinx/znet"
)

var room siface.IRoom = server.NewRoom("room0", 100)

func main() {
	defer room.StopRoom()
	go room.StartRoom()

	chat_server := znet.NewServer()

	chat_server.AddRounter(utils.NBroadcast, server.NewBroadcastRouter())
	chat_server.AddRounter(utils.NPrivateChat, server.NewPrivateChatRouter())
	chat_server.AddRounter(utils.NChangeName, server.NewChangeNameRouter())
	chat_server.AddRounter(utils.NWhos, server.NewWhosRouter())

	chat_server.SetOnConnStart(online)
	chat_server.SetOnConnStop(offline)

	chat_server.Serve()
}
