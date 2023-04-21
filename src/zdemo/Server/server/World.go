package server

import (
	"errors"
	"fmt"
	"main/src/zdemo/Server/siface"
	"strings"
	"sync"
)

type World struct {
	rooms map[string]siface.IRoom
	lock  sync.RWMutex

	cap uint32
}

func NewWorld() (world *World) {
	world = &World{
		rooms: make(map[string]siface.IRoom),
		lock:  sync.RWMutex{},
		cap:   5,
	}

	room_name := "default"
	world.rooms[room_name] = NewRoom(room_name, 10000)
	go world.rooms[room_name].StartRoom()

	return
}

func (this *World) CreateRoom(name string, cap uint32) (err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if len(this.rooms) > int(this.cap) {
		err = errors.New(fmt.Sprintf("hit max limit [%d] of rooms", this.cap))
		return
	}

	_, ok := this.rooms[name]
	if ok {
		err = errors.New(fmt.Sprintf("Room %d already exist", name))
		return
	}

	this.rooms[name] = NewRoom(name, cap)
	go this.rooms[name].StartRoom()

	return
}

func (this *World) GetRoom(name string) (room siface.IRoom, err error) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	room, ok := this.rooms[name]
	if !ok {
		err = errors.New(fmt.Sprintf("Room %d not exist", name))
		return
	}

	return
}

func (this *World) RemoveRoom(room siface.IRoom) {
	this.lock.Lock()
	defer this.lock.Unlock()

	delete(this.rooms, room.GetName())
	room.StopRoom()
}

func (this *World) GetAllRoomMsg() (names string) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	names = fmt.Sprintf("name\t\tnum\tcapacity\n")
	for name, room := range this.rooms {
		names += fmt.Sprintf("%s\t\t%d\t%d\n", name, room.GetUserNum(), room.GetCapacity())
	}

	names = strings.TrimRight(names, "\n")
	return
}

func (this *World) UserSwitchRoom(user siface.IUser, nroom_to string, nroom_from string) (err error) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	room_from, ok := this.rooms[nroom_from]
	if !ok {
		err = errors.New("Room_to is not exist")
		return
	}
	if _, err = room_from.GetUser(user.GetName()); err != nil {
		return
	}

	room_to, ok := this.rooms[nroom_to]
	if !ok {
		err = errors.New("Room_to is not exist")
		return
	}
	if _, err = room_to.GetUser(user.GetName()); err == nil {
		err = errors.New("Room_to already has a user with your same name, please rename your name")
		return
	}

	room_from.RemoveUser(user)
	err = room_to.AddUser(user)
	return
}
