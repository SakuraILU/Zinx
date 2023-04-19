package main

import (
	"bufio"
	"fmt"
	"main/src/zinx/znet"
	"net"
	"os"
)

func reader(conn net.Conn) {
	data_pack := znet.NewDataPack()

	for {
		head := make([]byte, data_pack.GetHeadLen())
		_, err := conn.Read(head)
		if err != nil {
			panic(err.Error())
		}
		msg, err := data_pack.UnpackHead(head)
		if err != nil {
			panic(err.Error())
		}
		conn.Read(msg.GetMsgData())
		fmt.Println(string(msg.GetMsgData()))
	}
}

func writer(conn net.Conn) {
	data_pack := znet.NewDataPack()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadBytes('\n')
		line = line[:len(line)-1]

		msg := znet.NewMessage(1, line)
		data, err := data_pack.Pack(msg)
		if err != nil {
			panic(err.Error())
		}
		if _, err := conn.Write(data); err != nil {
			panic(err.Error())
		}
	}
}

func main() {
	data_pack := znet.NewDataPack()

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		panic(err.Error())
	}

	head := make([]byte, data_pack.GetHeadLen())
	_, err = conn.Read(head)
	if err != nil {
		panic(err.Error())
	}
	msg, err := data_pack.UnpackHead(head)
	if err != nil {
		panic(err.Error())
	}
	conn.Read(msg.GetMsgData())
	fmt.Println(string(msg.GetMsgData()))
	// change name
	msg = znet.NewMessage(2, []byte("lky"))
	data, err := data_pack.Pack(msg)
	conn.Write(data)

	head = make([]byte, data_pack.GetHeadLen())
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

	// who
	msg = znet.NewMessage(3, []byte(""))
	data, err = data_pack.Pack(msg)
	conn.Write(data)

	head = make([]byte, data_pack.GetHeadLen())
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

	if msg.GetMsgID() == 10 {
		return
	} else {
		go reader(conn)
		go writer(conn)
	}

	select {}
}
