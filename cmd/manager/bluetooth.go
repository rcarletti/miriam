package main

import (
	"encoding/json"

	"sort"

	"github.com/go-mangos/mangos/protocol/pull"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/rcarletti/miriam/data"
)

const maxDistance = 3 // meters

func retrieveUserSettings(key string) data.UserSettings {
	var info UserEntry
	log.Debugln("Search by key:", key)
	db.First(&info, "mac_address = ?", key)
	return info.UserSettings
}

func handleBluetoothUpdates(updates chan data.UserSettings) {
	log.Debugln("Creating BT socket")

	sock, err := pull.NewSocket()
	if err != nil {
		panic(err)
	}

	sock.AddTransport(tcp.NewTransport())
	defer sock.Close()

	if err = sock.Listen("tcp://localhost:60000"); err != nil {
		panic(err)
	}

	var userInRange bool
	var userData data.BluetoothUser

	for {
		//ricevo un nuovo messaggio dal bluetooth
		msg, err := sock.Recv()
		if err != nil {
			log.Errorf("Could not receive from socket: %v", err)
			continue
		}

		var users data.NearUsers //lista utenti nelle vicinanze

		if err = json.Unmarshal(msg, &users); err != nil {
			log.Errorf("Invalid msg received: %v", err)
			continue
		}

		sort.Slice(users.BUsersList, func(i, j int) bool {
			return users.BUsersList[i].Distance < users.BUsersList[j].Distance
		})

		log.Debugln("Received:", users)

		if len(users.BUsersList) == 0 || users.BUsersList[0].Distance > maxDistance {
			if userInRange {
				userInRange = false
				updates <- data.UserSettings{}
			}
		} else {
			userInRange = true
			if users.BUsersList[0] != userData {
				userData = users.BUsersList[0]
				updates <- retrieveUserSettings(userData.MacAddress)
			}
		}
	}
}
