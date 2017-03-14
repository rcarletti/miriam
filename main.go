package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	owm "github.com/briandowns/openweathermap"
	"github.com/rcarletti/miriam/data"

	"encoding/json"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
)

type event struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

type email struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
}

type clientInfo struct {
	Weather     string  `json:"weather"`
	Temperature float64 `json:"temperature"`
	Unread      int64   `json:"unread"`
	EmailList   []email `json:"email_list"`
	Events      []event `json:"events"`
	UserID      string  `json:"user_id"`
}

func setMidnight(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 0, 0, t.Location())
}

func main() {

	var c clientInfo

	os.Setenv("OWM_API_KEY", "5bf842837d6a00751104eb08c3ace476")
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json") //read client's secret
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/gmail-go-quickstart.json
	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := data.GetClient(ctx, config, os.Args[1])

	user := "me"

	//*****************************************************************
	//gmail
	//*****************************************************************

	srvGmail, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	r, err := srvGmail.Users.Messages.List(user).Q("is:unread").Do()
	if err != nil {
		panic(err)
	}
	toBeRead := r.ResultSizeEstimate //unread emails
	c.Unread = toBeRead

	for i := 0; i < int(toBeRead); i++ {
		msg := r.Messages[i].Id
		m, _ := srvGmail.Users.Messages.Get(user, msg).Do()
		var name, mail, subject string
		//find senders, emails and subjects
		for _, h := range m.Payload.Headers {
			if h.Name == "From" {
				name = h.Value[:strings.LastIndex(h.Value, "<")-1]
				mail = h.Value[strings.LastIndex(h.Value, "<")+1 : len(h.Value)-1]
			}
			if h.Name == "Subject" {
				subject = h.Value
			}
		}
		c.EmailList = append(c.EmailList, email{
			Name:    name,
			Email:   mail,
			Subject: subject,
		})
	}

	//*****************************************************************
	//calendar
	//*****************************************************************

	srvCalendar, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	tonight := setMidnight(time.Now()).Format(time.RFC3339)
	now := time.Now().Format(time.RFC3339)

	events, err := srvCalendar.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(now).TimeMax(tonight).Do() //today events
	if err != nil {
		log.Fatalf("Unable to retrieve user's events. %v", err)
	}

	if len(events.Items) > 0 {
		for _, i := range events.Items {
			var when string
			// If the DateTime is an empty string the Event is an all-day Event.
			// So only Date is available.
			if i.Start.DateTime != "" {
				when = i.Start.DateTime
			} else {
				when = i.Start.Date
			}

			c.Events = append(c.Events, event{
				Name: i.Summary,
				Time: when,
			})

		}
	} else {
		fmt.Printf("No upcoming events found.\n")
	}

	//*****************************************************************
	//weather
	//*****************************************************************

	w, err := owm.NewCurrent("C", "it")
	if err != nil {
		log.Fatalln(err)
	}
	w.CurrentByName("pisa")
	c.Weather = w.Weather[0].Description
	c.Temperature = w.Main.Temp

	//*****************************************************************
	//json
	//*****************************************************************

	fout, err := os.OpenFile("output.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	json.NewEncoder(fout).Encode(c)
	fout.Close()

}
