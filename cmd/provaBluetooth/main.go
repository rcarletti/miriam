package main

import (
	"encoding/json"

	"time"

	"github.com/go-mangos/mangos/protocol/push"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/rcarletti/miriam/data"
)

func main() {
	testVectors := []data.NearUsers{
		{
			BUsersList: []data.BluetoothUser{
				data.BluetoothUser{"aaa", 300, 13},
				data.BluetoothUser{"bbb", 300, 13},
			},
		},
		{
			BUsersList: []data.BluetoothUser{
				data.BluetoothUser{"aaa", 100, 13},
				data.BluetoothUser{"bbb", 300, 13},
			},
		},
		{
			BUsersList: []data.BluetoothUser{
				data.BluetoothUser{"aac", 100, 13},
				data.BluetoothUser{"bbb", 100, 13},
			},
		},
		{},
	}

	sock, err := push.NewSocket()
	defer sock.Close()
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Dial("tcp://localhost:60000"); err != nil {
		panic(err)
	}

	for _, v := range testVectors {
		msg, _ := json.Marshal(v)
		if err = sock.Send(msg); err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Second)
	}
}
