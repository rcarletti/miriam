package main

import (
	"os"

	gcal "google.golang.org/api/calendar/v3"
	gmail "google.golang.org/api/gmail/v1"

	"github.com/rcarletti/miriam/calendar"
	"github.com/rcarletti/miriam/gauth"
	"github.com/rcarletti/miriam/mail"
	"github.com/rcarletti/miriam/weather"

	"encoding/json"
)

type clientInfo struct {
	Weather     string           `json:"weather"`
	Temperature float64          `json:"temperature"`
	Unread      int64            `json:"unread"`
	EmailList   []mail.Email     `json:"email_list"`
	Events      []calendar.Event `json:"events"`
	UserID      string           `json:"user_id"`
}

func init() {
	os.Setenv("OWM_API_KEY", "5bf842837d6a00751104eb08c3ace476")
}

func main() {
	var c clientInfo

	client, err := gauth.New(os.Args[1], "client_secret.json",
		gmail.MailGoogleComScope, gcal.CalendarReadonlyScope)
	if err != nil {
		panic(err)
	}

	c.EmailList, err = mail.Get(client)
	if err != nil {
		panic(err)
	}

	c.Events, err = calendar.Get(client)
	if err != nil {
		panic(err)
	}

	c.Temperature, c.Weather = weather.GetWeather("pisa")

	fout, err := os.OpenFile("output.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	json.NewEncoder(fout).Encode(c)
	fout.Close()

}
