package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	owm "github.com/briandowns/openweathermap"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("gmail-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 0, 0, t.Location())
}

type event struct {
	name string
	time string
}

type clientInfo struct {
	weather string
	mailNum int64
	senders []string
	events  []string
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
	client := getClient(ctx, config)

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

	c.mailNum = toBeRead

	fmt.Println("numero di mail da leggere:", toBeRead)

	for i := 0; i < int(toBeRead); i++ {
		msg := r.Messages[i].Id
		m, _ := srvGmail.Users.Messages.Get(user, msg).Do()
		//cerco il mittente
		for _, h := range m.Payload.Headers {
			if h.Name == "From" {
				//stampo solo il nome del mittente
				c.senders = append(c.senders, h.Value[:strings.LastIndex(h.Value, "<")-1])
			}
		}
	}

	fmt.Println("*****************************************")

	//*****************************************************************
	//PARTE DEL CALENDARIO
	//*****************************************************************

	srvCalendar, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	//creo un orario con data odierna e ora 23:59
	tonight := Bod(time.Now()).Format(time.RFC3339)
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
			fmt.Println(i.Summary, "ore", when[11:13], ":", when[14:16])

		}
	} else {
		fmt.Printf("No upcoming events found.\n")
	}

	w, err := owm.NewCurrent("C", "it")
	if err != nil {
		log.Fatalln(err)
	}

	w.CurrentByName("Pisa")
	fmt.Println(w.Weather[0])

}
