package main

import (
	"encoding/json"
	"time"

	"reflect"

	"github.com/go-mangos/mangos/protocol/push"
	"github.com/go-mangos/mangos/protocol/req"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/rcarletti/miriam/data"
)

const (
	cmdShutdown = "s"
	cmdDisplay  = "d"
)

func handleUserData(updates chan data.UserSettings) {
	log.Debugln("Creating GUI socket")

	sockGUI, err := push.NewSocket()
	if err != nil {
		panic(err)
	}
	sockGUI.AddTransport(tcp.NewTransport())
	defer sockGUI.Close()

	if err = sockGUI.Dial("tcp://localhost:40000"); err != nil {
		panic(err)
	}

	log.Debugln("Creating data socket")

	sockData, err := req.NewSocket()
	if err != nil {
		panic(err)
	}
	sockData.AddTransport(tcp.NewTransport())
	defer sockData.Close()

	if err = sockData.Dial("tcp://localhost:50000"); err != nil {
		panic(err)
	}

	var settings data.UserSettings

	for {
		// wait for a user update or a timeout to fetch new data
		select {
		case settings = <-updates:
			log.Debugln("Received user update:", settings)
		case <-time.After(3 * time.Minute):
			log.Debugln("Received timeout")
		}

		var msg []byte
		var info data.UserInfo

		// if settings is not empty, retrieve user data
		if !reflect.DeepEqual(settings, data.UserSettings{}) {
			// send settings
			js, err := json.Marshal(settings)
			if err != nil {
				log.Errorf("Could not marshal user settings: %v", err)
				continue
			}

			if err = sockData.Send(js); err != nil {
				log.Errorf("Could not send user settings: %v", err)
				continue
			}

			log.Debugln("Sent:", string(js))

			//receive user info from network
			if msg, err = sockData.Recv(); err != nil {
				log.Errorf("Could not retrieve user data: %v", err)
				continue
			}

			if err := json.Unmarshal(msg, &info); err != nil {
				log.Errorf("Could not unmarshal user data: %v", err)
				continue
			}

			info.Command = cmdDisplay
			log.Debugln("Sending display command")
		} else {
			// user is gone, send shutdown command
			info.Command = cmdShutdown
			log.Debugln("Sending shutdown command")
		}

		if msg, err = json.Marshal(info); err != nil {
			log.Errorf("Could not marshal user data: %v", err)
			continue
		}

		//send usr info to GUI
		if err = sockGUI.Send(msg); err != nil {
			log.Errorf("Could not send user data: %v", err)
			continue
		}
	}
}
