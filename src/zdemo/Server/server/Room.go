package server

import (
	"errors"
	"fmt"
	"main/src/zdemo/Server/siface"
	"strings"
	"sync"
	"time"
)

type Room struct {
	name string

	users map[string]siface.IUser
	lock  sync.RWMutex
	cap   uint32

	broadcast_msg chan []byte
	exit_chan     chan bool

	timeout uint32
}

func NewRoom(name string, cap uint32) (room *Room) {
	room = &Room{
		name:          name,
		users:         make(map[string]siface.IUser),
		lock:          sync.RWMutex{},
		cap:           cap,
		broadcast_msg: make(chan []byte),
		exit_chan:     make(chan bool),
		timeout:       300,
	}
	return
}

func (this *Room) GetName() string {
	return this.name
}

func (this *Room) GetUser(name string) (user siface.IUser, err error) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	user, ok := this.users[name]
	if !ok {
		err = errors.New(fmt.Sprintf("User %s not found", name))
	}
	return
}

func (this *Room) AddUser(user siface.IUser) (err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if uint32(len(this.users)) > this.cap {
		err = errors.New(fmt.Sprintf("fail to add user %s, out of capacity, Room %s with capacity %d", user.GetName(), this.name, this.cap))
		return
	}

	_, ok := this.users[user.GetName()]
	if ok {
		err = errors.New(fmt.Sprintf("user %s alread exist...", user.GetName()))
		return
	}

	this.users[user.GetName()] = user
	user.SetRoom(this)

	return
}
func (this *Room) RemoveUser(user siface.IUser) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.users, user.GetName())
	user.SetRoom(nil)
}

func (this *Room) GetAllUserMsg() (names string) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	for name, _ := range this.users {
		names += name + "\n"
	}

	names = strings.TrimRight(names, "\n")
	return
}

func (this *Room) ClearAll() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for name, user := range this.users {
		delete(this.users, name)
		defer user.StopWorker()
	}
}

func (this *Room) GetUserNum() uint32 {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return uint32(len(this.users))
}
func (this *Room) GetCapacity() uint32 {
	return this.cap
}

func (this *Room) StartRoom() {
	go this.broadCaster()
	this.StartTimeouter()
}

func (this *Room) StopRoom() {
	this.exit_chan <- true
	close(this.broadcast_msg)
	close(this.exit_chan)
}

func (this *Room) StartTimeouter() {
	ticker := time.NewTicker(time.Duration(this.timeout) * time.Second)
	for {
		select {
		case <-ticker.C:
			this.lock.Lock()
			for _, user := range this.users {
				if !user.IsActive() {
					go user.StopConn()
				} else {
					user.SetActive(false)
				}
			}
			this.lock.Unlock()
		case <-this.exit_chan:
			return
		}
	}
}

func (this *Room) broadCaster() {
	for {
		select {
		case data := <-this.broadcast_msg:
			for _, user := range this.users {
				user.AddMsg(data)
			}
		case <-this.exit_chan:
			return
		}
	}
}

func (this *Room) BroadCastMsg(data []byte) {
	this.broadcast_msg <- data
}
