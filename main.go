package main

import (
	"encoding/json"
	"net/http"
	"os"

	gcal "google.golang.org/api/calendar/v3"
	gmail "google.golang.org/api/gmail/v1"

	"fmt"

	"github.com/go-mangos/mangos/protocol/rep"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/rcarletti/miriam/calendar"
	"github.com/rcarletti/miriam/data"
	"github.com/rcarletti/miriam/gauth"
	"github.com/rcarletti/miriam/mail"
	"github.com/rcarletti/miriam/weather"
)

func init() {
	os.Setenv("OWM_API_KEY", "5bf842837d6a00751104eb08c3ace476")
}

func main() {
	var user data.UserInfo
	var msg []byte
	clientList := make(map[string]*http.Client)

	sock, err := rep.NewSocket()
	if err != nil {
		panic(err)
	}
	sock.AddTransport(tcp.NewTransport())

	if err = sock.Listen("tcp://localhost:" + os.Args[1]); err != nil {
		panic(err)
	}

	for {
		var usr data.UserSettings
		msg, err = sock.Recv()
		json.Unmarshal(msg, &usr)
		fmt.Println("ricevuto:", string(msg))
		client, ok := clientList[usr.UserID]
		//se non esiste la cartella per il client la creo
		if !ok {
			client, err = gauth.New(usr.UserID, "client_secret.json",
				gmail.MailGoogleComScope, gcal.CalendarReadonlyScope)
			if err != nil {
				panic(err)
			}
			clientList[usr.UserID] = client
		}

		fmt.Println(usr)

		user.EmailList, err = mail.Get(client, int64(usr.EmailMax))
		if err != nil {
			panic(err)
		}

		user.Events, err = calendar.Get(client, int64(usr.EventsMax))
		if err != nil {
			panic(err)
		}

		user.Temperature, user.Weather = weather.GetWeather(usr.Location)
		js, err := json.Marshal(user)
		if err != nil {
			panic(err)
		}
		if err = sock.Send(js); err != nil {
			panic(err)
		}
		fmt.Println("inviato: ", user)
	}

}
