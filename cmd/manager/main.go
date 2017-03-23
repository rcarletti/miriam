package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rcarletti/miriam/data"
)

type UserEntry struct {
	gorm.Model
	data.UserSettings
	MACAddress string `gorm:"primary_key"`
}

var db *gorm.DB
var log = logrus.WithField("app", "manager")

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	var err error

	log.Debugln("Starting manager")

	if db, err = gorm.Open("sqlite3", "usersDB.db"); err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&UserEntry{})

	log.Debugln("Opened database")

	// db.Create(&UserEntry{
	// 	MACAddress: "aaa",
	// 	UserSettings: data.UserSettings{
	// 		EmailMax:  2,
	// 		EventsMax: 2,
	// 		Location:  "pisa",
	// 		UserID:    "miricd ..am",
	// 	},
	// })

	// db.Create(&UserEntry{
	// 	MACAddress: "bbb",
	// 	UserSettings: data.UserSettings{
	// 		EmailMax:  3,
	// 		EventsMax: 3,
	// 		Location:  "Londra",
	// 		UserID:    "Rossella",
	// 	},
	// })

	updates := make(chan data.UserSettings)

	go handleBluetoothUpdates(updates)
	go handleUserData(updates)

	select {}
}
