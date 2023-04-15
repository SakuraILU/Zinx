package main

import (
	"fmt"
	"main/src/zinx/znet"
	"net"
	"time"
)

func main() {
	var msg_id uint32 = 1 // 0 for echo, 1 for ping

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		panic(err.Error())
	}
	data_pack := znet.NewDataPack()
	for i := 0; i < 10; i++ {

		msg := znet.NewMessage(msg_id, []byte("Hello, server\n..fegfaehfa fefio\n fjwieofhew\n EOF"))
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
}
