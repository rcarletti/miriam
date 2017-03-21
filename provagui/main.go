package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-mangos/mangos/protocol/pull"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/rcarletti/miriam/data"
)

func main() {

	sock, err := pull.NewSocket()
	defer sock.Close()
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Listen("tcp://localhost:40000"); err != nil {
		panic(err)
	}

	for {
		msg, err := sock.Recv()
		if err != nil {
			panic(err)
		}
		var usrInfo data.UserInfo
		json.Unmarshal(msg, &usrInfo)
		fmt.Println(string(msg))
	}

}
