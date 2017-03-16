package calendar

import (
	"fmt"
	"net/http"
	"time"

	calendar "google.golang.org/api/calendar/v3"
)

type Event struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

func setMidnight(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 0, 0, t.Location())
}

func Get(client *http.Client, max int64) ([]Event, error) {
	var eventList []Event
	srvCalendar, err := calendar.New(client)
	if err != nil {
		return nil, err
	}

	now := time.Now().Format(time.RFC3339)

	events, err := srvCalendar.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(now).MaxResults(max).Do() //today events
	if err != nil {
		return nil, err
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

			eventList = append(eventList, Event{
				Name: i.Summary,
				Time: when,
			})

		}
	} else {
		fmt.Printf("No upcoming events found.\n")
	}
	return eventList, nil
}
