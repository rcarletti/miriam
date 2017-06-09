package main

import (
	"encoding/json"

	"sort"

	"github.com/go-mangos/mangos/protocol/pull"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/rcarletti/miriam/data"
)

const maxDistance = 200 // cm(?)

func retrieveUserSettings(key string) (data.UserSettings, bool) {
	var info UserEntry
	log.Debugln("Search by key:", key)
	v := db.First(&info, "mac_address = ?", key)
	return info.UserSettings, v.RowsAffected == 0
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
		//receive from presence module
		msg, err := sock.Recv()
		if err != nil {
			log.Errorf("Could not receive from socket: %v", err)
			continue
		}

		var users data.NearUsers

		if err = json.Unmarshal(msg, &users); err != nil {
			log.Errorf("Invalid msg received: %v", err)
			continue
		}
		//sort users list by distance
		sort.Slice(users.BUsersList, func(i, j int) bool {
			return users.BUsersList[i].Distance < users.BUsersList[j].Distance
		})

		log.Debugln("Received:", users)

		if len(users.BUsersList) == 0 || users.BUsersList[0].Distance > maxDistance { //no one nearby
			if userInRange { //someone went away!
				userInRange = false
				userData = data.BluetoothUser{}
				updates <- data.UserSettings{}
			}
		} else { //someone is near
			userInRange = true
			keepOldUser := false

			// current user has priority, so check if he is still there
			for _, u := range users.BUsersList {
				if u.MacAddress == userData.MacAddress && u.Distance <= maxDistance {
					userData = u
					keepOldUser = true
					break
				}
			}

			// if the old user is gone, replace him with the closest one
			if !keepOldUser {
				e := true
				var v data.UserSettings

				for i := 0; i < len(users.BUsersList) && e; i++ {
					userData = users.BUsersList[i]
					v, e = retrieveUserSettings(userData.MacAddress)
					if !e {
						updates <- v
					}
				}
			}
		}
	}
}
