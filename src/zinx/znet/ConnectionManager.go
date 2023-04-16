package znet

import (
	"errors"
	"main/src/zinx/ziface"
	"sync"
)

type ConnectionManager struct {
	conns map[uint32]ziface.IConnection
	size  uint32
	lock  sync.RWMutex
}

func NewConnectionManager() (cm *ConnectionManager) {
	cm = &ConnectionManager{
		conns: make(map[uint32]ziface.IConnection),
		size:  0,
		lock:  sync.RWMutex{},
	}
	return
}

func (this *ConnectionManager) GetConn(id uint32) (conn ziface.IConnection, err error) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	conn, ok := this.conns[id]
	if !ok {
		err = errors.New("connection not found")
	}

	return
}

func (this *ConnectionManager) Add(conn ziface.IConnection) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.conns[conn.GetConnID()] = conn
	this.size++
}

func (this *ConnectionManager) Remove(conn ziface.IConnection) {
	this.lock.Lock()
	defer this.lock.Unlock()

	delete(this.conns, conn.GetConnID())
	this.size--
}

// not only clear the map, but also Stop all the connections to free sockets and other resources
func (this *ConnectionManager) ClearAll() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for id, conn := range this.conns {
		delete(this.conns, id)
		defer conn.Stop()
	}
}

func (this *ConnectionManager) Size() uint32 {
	return this.size
}
