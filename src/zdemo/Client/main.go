package main

import (
	"fmt"
	"main/src/zinx/znet"
	"net"
	"sync"
	"time"
)

var wg sync.WaitGroup

func establish_dial(msg_id uint32, client_id int) {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		panic(err.Error())
	}
	data_pack := znet.NewDataPack()
	for i := 0; i < 10; i++ {

		msg := znet.NewMessage(msg_id, []byte(fmt.Sprintf("Hello, server..\n I am client %d", client_id)))
		buf, err := data_pack.Pack(msg)
		if err != nil {
			panic(err.Error())
		}
		_, err = conn.Write(buf)
		if err != nil {
			panic(err.Error())
		}

		head := make([]byte, data_pack.GetHeadLen())
		_, err = conn.Read(head)
		if err != nil {
			panic(err.Error())
		}
		msg, err = data_pack.UnpackHead(head)
		if err != nil {
			panic(err.Error())
		}
		conn.Read(msg.GetMsgData())
		fmt.Println(string(msg.GetMsgData()))

		time.Sleep(1 * time.Second)
	}

	wg.Done()
}

func main() {
	for i := 0; i < 10; i++ {
		msg_id := 1
		go establish_dial(uint32(msg_id), i)
		wg.Add(1)
	}
	wg.Wait()
}
