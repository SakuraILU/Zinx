package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"main/src/zdemo/Server/utils"
	"main/src/zinx/ziface"
	"main/src/zinx/znet"
	"net"
	"os"
	"strings"
)

var cmdresp_chan chan bool = make(chan bool)

func reader(conn net.Conn) {
	data_pack := znet.NewDataPack()

	for {
		fmt.Printf(">>>")

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
		if err != nil {
			panic(err.Error())
		}

		var msg ziface.IMessage
		switch cmd.Cmdtype {
		case To:
			msg = znet.NewMessage(utils.NPrivateChat, []byte(cmd.Arg))
		case Bc:
			msg = znet.NewMessage(utils.NBroadcast, []byte(cmd.Arg))
		case Rename:
			msg = znet.NewMessage(utils.NChangeName, []byte(cmd.Arg))
		case Whos:
			msg = znet.NewMessage(utils.NWhos, []byte(cmd.Arg))
		case Exit:
			exit_chan <- true
			return
		case NewRoom:
			msg = znet.NewMessage(utils.NNewRoom, []byte(cmd.Arg))
		case SwitchRoom:
			msg = znet.NewMessage(utils.NSwitchRoom, []byte(cmd.Arg))
		default:
			fmt.Println("[Error]: Unsupported cmd")
			fmt.Printf(">>>")
			continue
		}

		data, err := data_pack.Pack(msg)
		if err != nil {
			panic(err.Error())
		}
		if _, err := conn.Write(data); err != nil {
			panic(err.Error())
		}

		// select {
		// case <-cmdresp_chan:
		// }
	}
}

const (
	Bc = 1 + iota // in case..input a empty line, and Cmd.Cmdtype = 0 in default... so treat 0 as invalid..
	To
	Rename
	Whos
	NewRoom
	SwitchRoom
	Exit
)

var exit_chan chan bool = make(chan bool)

type Cmd struct {
	Cmdtype int
	Arg     string
}

func cmdParse() (cmd *Cmd, err error) {
	reader := bufio.NewReader(os.Stdin)

	cmd = &Cmd{}

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
		cmd.Cmdtype = To
		cmd.Arg = strings.SplitN(line, " ", 2)[1]
		cmd.Arg = strings.TrimRight(cmd.Arg, "\n") + ":"
		cmd.Arg += pasrse_txt()
	} else if strings.HasPrefix(line, "bc") {
		cmd.Cmdtype = Bc
		cmd.Arg = pasrse_txt()
		// fmt.Printf(cmd.Arg)
	} else if strings.HasPrefix(line, "rename") {
		cmd.Cmdtype = Rename
		cmd.Arg = strings.SplitN(line, " ", 2)[1]
		cmd.Arg = strings.TrimRight(cmd.Arg, "\n")
	} else if strings.HasPrefix(line, "whos") {
		cmd.Cmdtype = Whos
	} else if strings.HasPrefix(line, "exit") {
		cmd.Cmdtype = Exit
	} else if strings.HasPrefix(line, "rnew") {
		cmd.Cmdtype = NewRoom
		cmd.Arg = strings.SplitN(line, " ", 2)[1]
		cmd.Arg = strings.TrimRight(cmd.Arg, "\n")
	} else if strings.HasPrefix(line, "rswitch") {
		cmd.Cmdtype = SwitchRoom
		cmd.Arg = strings.SplitN(line, " ", 2)[1]
		cmd.Arg = strings.TrimRight(cmd.Arg, "\n")
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
