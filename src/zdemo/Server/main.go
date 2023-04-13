package main

import "main/src/zinx/znet"

func main() {
	server := znet.NewServer("sever_v0.1")
	server.Serve()
}
