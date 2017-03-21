package main

import (
	"encoding/json"

	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rcarletti/miriam/data"

	"github.com/go-mangos/mangos/protocol/pull"
	"github.com/go-mangos/mangos/protocol/push"
	"github.com/go-mangos/mangos/protocol/req"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/jinzhu/gorm"
)

type userDB struct {
	gorm.Model
	data.UserSettings
	MACAddress string `gorm:"primary_key"`
}

func main() {
	//creazione del database, da fare una volta sola
	//	os.Remove("usersDB.db")
	for {
		shutDown, _ := json.Marshal(data.UserInfo{"", 0, 0, nil, nil, "", "s"})

		//bluetooth

		sockBluetooth, err := pull.NewSocket()
		if err != nil {
			panic(err)
		}
		sockBluetooth.AddTransport(tcp.NewTransport())
		if err = sockBluetooth.Listen("tcp://localhost:60000"); err != nil {
			panic(err)
		}
		defer sockBluetooth.Close()
		//receive from bluetooth
		bMsg, err := sockBluetooth.Recv()
		if err != nil {
			panic(err)
		}
		//decode
		var userList data.NearUsers
		json.Unmarshal(bMsg, &userList)

		db, err := gorm.Open("sqlite3", "usersDB.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		sockNetwork, err := req.NewSocket()
		if err != nil {
			panic(err)
		}
		sockNetwork.AddTransport(tcp.NewTransport())
		if err = sockNetwork.Dial("tcp://localhost:50000"); err != nil {
			panic(err)
		}
		defer sockNetwork.Close()

		db.AutoMigrate(&userDB{})

		// db.Create(&userDB{
		// 	MACAddress: "aaa",
		// 	UserSettings: data.UserSettings{
		// 		EmailMax:  2,
		// 		EventsMax: 2,
		// 		Location:  "pisa",
		// 		UserID:    "miriam",
		// 	},
		// })

		// db.Create(&userDB{
		// 	MACAddress: "bbb",
		// 	UserSettings: data.UserSettings{
		// 		EmailMax:  3,
		// 		EventsMax: 3,
		// 		Location:  "Londra",
		// 		UserID:    "Rossella",
		// 	},
		// })

		//query by MACAddress
		var info userDB
		//query con il mac_address del primo elemento in coda (piu vicino)
		var tmpMAC string
		for i := range userList.BUsersList {
			if userList.BUsersList[i].Distance < 3 {
				tmpMAC = userList.BUsersList[i].MacAddress
				break
			} else {
				tmpMAC = ""
			}
		}

		if tmpMAC == "" {
			sockNetwork.Send(shutDown)
			continue
		}

		db.First(&info, "mac_address = ?", userList.BUsersList[0].MacAddress)

		//NETWORK

		//encode user settings

		js, _ := json.Marshal(info.UserSettings)

		//send user settings to network

		if err = sockNetwork.Send(js); err != nil {
			panic(err)
		}
		fmt.Println("inviato: ", string(js))

		//receive user info from network
		msg, err := sockNetwork.Recv()
		if err != nil {
			panic(err)
		}

		//GUI

		sockGUI, err := push.NewSocket()
		if err != nil {
			panic(err)
		}
		defer sockGUI.Close()
		sockGUI.AddTransport(tcp.NewTransport())
		if err = sockGUI.Dial("tcp://localhost:40000"); err != nil {
			panic(err)
		}

		//send user info to gui

		if err = sockGUI.Send(msg); err != nil {
			panic(err)
		}
	}

}
