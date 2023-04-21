package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"main/src/zdemo/Server/utils"
	"main/src/zinx/znet"
	"net"
	"os"
	"strings"
)

type Cmd struct {
	Cmdtype int
	Arg     string
}

var promt string = ">>>"
var cmdresp_chan chan bool = make(chan bool)
var exit_chan chan bool = make(chan bool)

func reader(conn net.Conn) {
	data_pack := znet.NewDataPack()

	for {
		fmt.Printf(promt)

		head := make([]byte, data_pack.GetHeadLen())
		_, err := conn.Read(head)
		fmt.Printf("\r")
		if errors.Is(err, io.EOF) {
			fmt.Println("you are kicked out because of no activities for too long")
			exit_chan <- true
			return
		} else if err != nil {
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

	for {
		cmd, err := cmdParse()
		if err != nil || cmd.Cmdtype == utils.NErr {
			fmt.Printf(promt)
			continue
		}

		msg := znet.NewMessage(uint32(cmd.Cmdtype), []byte(cmd.Arg))
		data, err := data_pack.Pack(msg)
		if err != nil {
			panic(err.Error())
		}
		if _, err := conn.Write(data); err != nil {
			panic(err.Error())
		}
	}
}

func cmdParse() (cmd *Cmd, err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = rerr.(error)
			fmt.Println("cmd with invalid argument")
		}
	}()

	cmd = &Cmd{Cmdtype: utils.NErr}

	reader := bufio.NewReader(os.Stdin)

	pasrse_txt := func() (arg string) {
		arg += "\n"
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}

			if strings.HasPrefix(line, "eof") {
				arg = strings.TrimRight(arg, "\n")
				return
			}
			arg += "  "
			arg += line
		}
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	if strings.HasPrefix(line, "to") {
		cmd.Cmdtype = utils.NPrivateChat
		cmd.Arg = strings.SplitN(line, " ", 2)[1]
		cmd.Arg = strings.TrimRight(cmd.Arg, "\n") + ":"
		cmd.Arg += pasrse_txt()
	} else if strings.HasPrefix(line, "bc") {
		cmd.Cmdtype = utils.NBroadcast
		cmd.Arg = pasrse_txt()
	} else if strings.HasPrefix(line, "rename") {
		cmd.Cmdtype = utils.NChangeName
		cmd.Arg = strings.SplitN(line, " ", 2)[1]
		cmd.Arg = strings.TrimRight(cmd.Arg, "\n")
	} else if strings.HasPrefix(line, "whos") {
		cmd.Cmdtype = utils.NWhos
	} else if strings.HasPrefix(line, "newroom") {
		cmd.Cmdtype = utils.NNewRoom
		cmd.Arg = strings.SplitN(line, " ", 2)[1]
		cmd.Arg = strings.TrimRight(cmd.Arg, "\n")
	} else if strings.HasPrefix(line, "enter") {
		cmd.Cmdtype = utils.NSwitchRoom
		cmd.Arg = strings.SplitN(line, " ", 2)[1]
		cmd.Arg = strings.TrimRight(cmd.Arg, "\n")
	} else if strings.HasPrefix(line, "curroom") {
		cmd.Cmdtype = utils.NCurrentRoom
	} else if strings.HasPrefix(line, "rooms") {
		cmd.Cmdtype = utils.NRooms
	} else if strings.HasPrefix(line, "exit") {
		exit_chan <- true
		select {}
	} else {
		cmd.Cmdtype = utils.NErr
		fmt.Println("unsupported cmd")
	}

	return
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		panic(err.Error())
	}

	data_pack := znet.NewDataPack()

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

	if msg.GetMsgID() == utils.NErr {
		return
	} else {
		go reader(conn)
		go writer(conn)
	}

	select {
	case <-exit_chan:
		return
	}
}
