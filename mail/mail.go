package mail

import (
	"net/http"
	"strings"

	gmail "google.golang.org/api/gmail/v1"
)

type Email struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
}

func Get(client *http.Client, max int64) ([]Email, int64, error) {
	var mailList []Email

	srvGmail, err := gmail.New(client)
	if err != nil {
		return nil, 0, err
	}

	r, err := srvGmail.Users.Messages.List("me").Q("is:unread").MaxResults(max).Do()
	if err != nil {
		return nil, 0, err
	}

	toBeRead := r.ResultSizeEstimate //unread emails

	var maxMail int64

	if toBeRead < max {
		maxMail = toBeRead
	} else {
		maxMail = max
	}

	for i := 0; i < int(maxMail); i++ {
		m, err := srvGmail.Users.Messages.Get("me", r.Messages[i].Id).Do()
		if err != nil {
			return nil, 0, err
		}

		var mail Email
		//find senders, emails and subjects
		for _, h := range m.Payload.Headers {
			switch h.Name {
			case "From":
				mail.Name = h.Value[:strings.LastIndex(h.Value, "<")-1]
				mail.Email = h.Value[strings.LastIndex(h.Value, "<")+1 : len(h.Value)-1]
			case "Subject":
				mail.Subject = h.Value
			}
		}

		mailList = append(mailList, mail)
	}

	return mailList, toBeRead, nil
}
