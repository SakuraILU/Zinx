package server

import (
	"main/src/zdemo/Server/siface"
	"main/src/zinx/ziface"
	"net"
)

type User struct {
	name string
	conn ziface.IConnection

	room siface.IRoom

	msg_chan  chan []byte
	exit_chan chan bool
}

func NewUser(name string, conn ziface.IConnection, room siface.IRoom) (user *User) {
	user = &User{
		name:      name,
		conn:      conn,
		room:      room,
		msg_chan:  make(chan []byte, 10),
		exit_chan: make(chan bool),
	}
	return
}

func (this *User) GetName() string {
	return this.name
}

func (this *User) SetName(name string) {
	this.name = name
}

func (this *User) GetRemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *User) AddMsg(data []byte) {
	this.msg_chan <- data
}

func (this *User) StartWorker() {
	for {
		select {
		case data := <-this.msg_chan:
			this.conn.SendMsg(0, data)
		case <-this.exit_chan:
			return
		}
	}
}

func (this *User) StopWorker() {
	this.exit_chan <- true
	close(this.msg_chan)
	close(this.exit_chan)
}

func (this *User) GetRoom() siface.IRoom {
	return this.room
}
