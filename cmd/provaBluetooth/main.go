package main

import (
	"encoding/json"

	"time"

	"github.com/go-mangos/mangos/protocol/push"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/rcarletti/miriam/data"
)

func main() {

	usr1 := data.BluetoothUser{"aaa", 2, "13:04"}
	usr2 := data.BluetoothUser{"bbb", 4, "13:04"}

	var usrList data.NearUsers
	usrList.BUsersList = append(usrList.BUsersList, usr1)
	usrList.BUsersList = append(usrList.BUsersList, usr2)

	encodedUserList, err := json.Marshal(usrList)
	if err != nil {
		panic(err)
	}

	sock, err := push.NewSocket()
	defer sock.Close()
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Dial("tcp://localhost:60000"); err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		if err = sock.Send(encodedUserList); err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Second)
	}

	msg, _ := json.Marshal(data.NearUsers{})
	sock.Send(msg)
}
