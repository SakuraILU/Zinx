package server

import (
	"main/src/zdemo/Server/siface"
	"main/src/zdemo/Server/utils"
	"main/src/zinx/ziface"
	"net"
	"sync"
)

type User struct {
	name string
	conn ziface.IConnection

	room siface.IRoom

	msg_chan  chan []byte
	exit_chan chan bool

	is_active bool
	lock      sync.RWMutex
}

func NewUser(name string, conn ziface.IConnection, room siface.IRoom) (user *User) {
	user = &User{
		name:      name,
		conn:      conn,
		room:      room,
		msg_chan:  make(chan []byte, 10),
		exit_chan: make(chan bool),
		is_active: true,
		lock:      sync.RWMutex{},
	}
	return
}

func (this *User) GetName() string {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.name
}

func (this *User) SetName(name string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.name = name
}

func (this *User) GetRemoteAddr() net.Addr {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.conn.RemoteAddr()
}

func (this *User) AddMsg(data []byte) {
	this.msg_chan <- data
}

func (this *User) StartWorker() {
	for {
		select {
		case data := <-this.msg_chan:
			this.conn.SendMsg(utils.NMsgResponse, data)
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

func (this *User) StopConn() {
	// defer this.StopWorker()
	this.conn.SendMsg(utils.NErr, []byte(""))
	this.conn.Stop()
}

func (this *User) GetRoom() siface.IRoom {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.room
}

func (this *User) SetActive(is_active bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.is_active = is_active
}

func (this *User) IsActive() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.is_active
}
