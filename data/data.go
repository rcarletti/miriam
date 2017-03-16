package data

import (
	"github.com/rcarletti/miriam/calendar"
	"github.com/rcarletti/miriam/mail"
)

type UserInfo struct {
	Weather     string           `json:"weather"`
	Temperature float64          `json:"temperature"`
	Unread      int64            `json:"unread"`
	EmailList   []mail.Email     `json:"email_list"`
	Events      []calendar.Event `json:"events"`
	UserID      string           `json:"user_id"`
}

type UserSettings struct {
	UserID    string `json:"user_ID"`
	EmailMax  int    `json:"email_max"`
	EventsMax int    `json:"events_max"`
	Location  string `json:"location"`
}
