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
	Name string
	Time string
}

type sender struct {
	Name  string
	Email string
}

type clientInfo struct {
	Weather string
	MailNum int64
	Senders []sender
	Events  []event
}

func bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 0, 0, t.Location())
}

func main() {
	var c clientInfo

	os.Setenv("OWM_API_KEY", "5bf842837d6a00751104eb08c3ace476")
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/gmail-go-quickstart.json
	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := data.GetClient(ctx, config)

	//*****************************************************************
	//PARTE DELLE MAIL
	//*****************************************************************

	srvGmail, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	user := "me"

	r, err := srvGmail.Users.Messages.List(user).Q("is:unread").Do()
	toBeRead := r.ResultSizeEstimate

	c.MailNum = toBeRead

	fmt.Println("numero di mail da leggere:", toBeRead)

	for i := 0; i < int(toBeRead); i++ {
		msg := r.Messages[i].Id
		m, _ := srvGmail.Users.Messages.Get(user, msg).Do()
		//cerco il mittente
		for _, h := range m.Payload.Headers {
			if h.Name == "From" {
				//stampo solo il nome del mittente
				var s sender
				s.Name = h.Value[:strings.LastIndex(h.Value, "<")-1]
				s.Email = h.Value[strings.LastIndex(h.Value, "<"):])
				c.Senders = append(c.Senders, s)
			}
		}
	}


	//*****************************************************************
	//PARTE DEL CALENDARIO
	//*****************************************************************

	srvCalendar, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	//creo un orario con data odierna e ora 23:59
	tonight := bod(time.Now()).Format(time.RFC3339)
	now := time.Now().Format(time.RFC3339)

	println(now)

	//ricavo gli eventi della giornata

	events, err := srvCalendar.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(now).TimeMax(tonight).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve user's events. %v", err)
	}

	//stampo gli eventi

	fmt.Println("Eventi in calendario:")
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
			var e event
			e.Name = i.Summary
			e.Time = when

			c.Events = append(c.Events, e)

		}
	} else {
		fmt.Printf("No upcoming events found.\n")
	}

	w, err := owm.NewCurrent("C", "it")
	if err != nil {
		log.Fatalln(err)
	}

	w.CurrentByName("Pisa")
	c.Weather = w.Weather[0].Description

	//trasformo la struttura in json
	fout, err := os.OpenFile("output.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)

	enc := json.NewEncoder(fout).Encode(c)

	fout.Close()

}
