package main

import (
	"os"

	"encoding/json"

	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rcarletti/miriam/data"

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
	os.Remove("usersDB.db")

	db, err := gorm.Open("sqlite3", "usersDB.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&userDB{})

	db.Create(&userDB{
		MACAddress: "aaa",
		UserSettings: data.UserSettings{
			EmailMax:  2,
			EventsMax: 2,
			Location:  "pisa",
			UserID:    "miriam",
		},
	})

	var info userDB
	db.First(&info, "user_id = ?", "miriam")

	sockNetwork, err := req.NewSocket() //socket per la parte network
	if err != nil {
		panic(err)
	}
	sockNetwork.AddTransport(tcp.NewTransport())
	if err = sockNetwork.Dial("tcp://localhost:" + os.Args[1]); err != nil {
		panic(err)
	}

	js, _ := json.Marshal(info.UserSettings)

	if err = sockNetwork.Send(js); err != nil {
		panic(err)
	}
	fmt.Println("inviato: ", string(js))
	msg, err := sockNetwork.Recv()
	if err != nil {
		panic(err)
	}
	var usrInfo data.UserInfo
	json.Unmarshal(msg, &usrInfo)
	fmt.Println("ricevuto: ", usrInfo)

}
