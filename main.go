package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	gcal "google.golang.org/api/calendar/v3"
	gmail "google.golang.org/api/gmail/v1"

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
	defer sock.Close()
	sock.AddTransport(tcp.NewTransport())

	if err = sock.Listen("tcp://localhost:50000"); err != nil {
		panic(err)
	}

	for {
		var usr data.UserSettings
		var waitG sync.WaitGroup
		waitG.Add(3)

		//receive request from manager

		msg, err = sock.Recv()
		json.Unmarshal(msg, &usr)
		//fmt.Println("ricevuto:", string(msg))
		client, ok := clientList[usr.UserID]

		if !ok {
			client, err = gauth.New(usr.UserID, "client_secret.json",
				gmail.MailGoogleComScope, gcal.CalendarReadonlyScope)
			if err != nil {
				panic(err)
			}
			clientList[usr.UserID] = client
		}

		//retrieve userinfo

		go func() {
			user.EmailList, user.Unread, err = mail.Get(client, int64(usr.EmailMax))
			waitG.Done()
		}()

		go func() {
			user.Events, err = calendar.Get(client, int64(usr.EventsMax))
			waitG.Done()
		}()

		go func() {
			user.Temperature, user.Weather = weather.GetWeather(usr.Location)
			waitG.Done()
		}()

		waitG.Wait()
		//encode userinfo
		js, err := json.Marshal(user)
		if err != nil {
			panic(err)
		}

		//send userinfo to manager
		if err = sock.Send(js); err != nil {
			panic(err)
		}
		//fmt.Println("inviato: ", user)
	}

}
