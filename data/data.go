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
	Command     string           `json:"command"`
}

type UserSettings struct {
	UserID    string `json:"user_ID"`
	EmailMax  int    `json:"email_max"`
	EventsMax int    `json:"events_max"`
	Location  string `json:"location"`
}

type NearUsers struct {
	BUsersList []BluetoothUser `json:"BuserList"`
	userTOT    int             `json:"userTOT"`
}

type BluetoothUser struct {
	MacAddress string  `json:"mac_address"`
	Distance   uint    `json:"distance"`    // distance in cm
	Timestamp  uint    `json:"time_stamp"`  // timestamp as seconds
}
