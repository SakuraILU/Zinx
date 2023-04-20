package main

import (
	"fmt"
	"main/src/zdemo/Server/server"
	"main/src/zdemo/Server/siface"
	"main/src/zdemo/Server/utils"
	"main/src/zinx/ziface"
)

func online(conn ziface.IConnection) {
	room, err := world.GetRoom("default")
	if err != nil {
		conn.SendMsg(utils.NErr, []byte(err.Error()))
		return
	}

	user := server.NewUser(conn.RemoteAddr().String(), conn, room)
	user.SetActive(true)
	err = room.AddUser(user)
	if err != nil {
		conn.SendMsg(utils.NErr, []byte(err.Error()))
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

	fmt.Printf("[%s] is offline", user.GetName())

	room.RemoveUser(user)

	user.StopWorker()

	msg := fmt.Sprintf("[%s] is offline", user.GetName())
	room.BroadCastMsg([]byte(msg))
}
