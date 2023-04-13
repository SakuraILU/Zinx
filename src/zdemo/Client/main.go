package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < 10; i++ {
		_, err = conn.Write([]byte("Hello, world"))
		if err != nil {
			panic(err.Error())
		}
		resp := make([]byte, 512)
		_, err = conn.Read(resp)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(string(resp))
		time.Sleep(1 * time.Second)
	}
}
