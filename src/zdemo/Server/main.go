package main

import (
	"fmt"
	"main/src/zdemo/Server/server"
	"main/src/zdemo/Server/siface"
	"main/src/zinx/ziface"
	"main/src/zinx/znet"
)

var room siface.IRoom = server.NewRoom("room0", 100)

func main() {
	defer room.StopRoom()
	go room.StartRoom()

	chat_server := znet.NewServer()

	chat_server.AddRounter(1, server.NewBroadcastRouter())
	chat_server.AddRounter(2, server.NewChangeNameRouter())
	chat_server.AddRounter(3, server.NewWhoRouter())

	chat_server.SetOnConnStart(online)
	chat_server.SetOnConnStop(offline)

	chat_server.Serve()
}

func online(conn ziface.IConnection) {
	fmt.Println("here")
	user := server.NewUser(conn.RemoteAddr().String(), conn, room)
	err := room.AddUser(user)
	if err != nil {
		conn.SendMsg(10, []byte(err.Error()))
		return
	} else {
		go user.StartWorker()
		fmt.Printf("user %s is online...", user.GetName())
		room.BroadCastMsg([]byte(fmt.Sprintf("[%s] is online", user.GetName())))

		conn.SetProperty("user", user)
	}
}

func offline(conn ziface.IConnection) {
	iuser, err := conn.GetProperty("user")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	user := iuser.(siface.IUser)
	room := user.GetRoom()

	room.RemoveUser(user)

	msg := fmt.Sprintf("[%s] is offline", user.GetName())
	room.BroadCastMsg([]byte(msg))
}
