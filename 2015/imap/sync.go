package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorgias/imap-sync/sync"
	"github.com/xarg/imap"

	"database/sql"

	_ "github.com/lib/pq"
)

// config struct read from config.json
type Config struct {
	Db struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"db"`
	Newrelic struct {
		AppName    string `json:"app_name"`
		LicenseKey string `json:"license_key"`
		Verbose    bool   `json:"verbose"`
	} `json:"newrelic"`
	Oauth2 struct {
		Google struct {
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
		} `json:"google"`
	} `json:"oauth2"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

const (
	Addr          = "imap.gmail.com:993"
	TOKEN_URL     = "https://accounts.google.com/o/oauth2/token"
	SYNC_INTERVAL = time.Second // change to minute
)

func main() {
	// setup a logger
	imap.DefaultLogger = log.New(os.Stdout, "", 0)
	// imap.DefaultLogMask = imap.LogConn | imap.LogRaw

	// read config
	cfile, _ := os.Open("config.json")
	decoder := json.NewDecoder(cfile)
	conf := Config{}
	err := decoder.Decode(&conf)
	if err != nil {
		log.Fatal(err)
	}

	// setup newrelic
	/*
		agent := gorelic.NewAgent()
		agent.Verbose = conf.Newrelic.Verbose
		agent.NewrelicName = conf.Newrelic.AppName
		agent.NewrelicLicense = conf.Newrelic.LicenseKey
		agent.Run()
	*/

	// Open db
	db, err := sql.Open(conf.Db.Type, conf.Db.URL)
	if err != nil {
		log.Fatal(err)
	}

	// Ignore some common folders
	ignoredFolders := []string{"[Gmail]/Trash", "[Gmail]/Spam", "[Gmail]/Bin", "Drafts", "[Gmail]/Drafts"}

	for {
		// Get all the accounts to sync
		accounts, err := sync.Accounts(db)
		if err != nil {
			log.Fatal(err)
		}

		for _, acc := range accounts {
			// Skip this account if it's not time yet
			if acc.SyncDatetime != nil && time.Now().UTC().Before(acc.SyncDatetime.Add(SYNC_INTERVAL)) {
				continue
			}
			// Given the refresh token we get a new access_token
			// We need to do this because the access_token expires.
			resp, err := http.PostForm(TOKEN_URL,
				url.Values{
					"client_id":     {conf.Oauth2.Google.ClientID},
					"client_secret": {conf.Oauth2.Google.ClientSecret},
					"refresh_token": {acc.RefreshToken},
					"grant_type":    {"refresh_token"},
				})
			if err != nil {
				log.Fatal("Failed to get new token refresh_token", err)
			}

			// Read the json data
			// XXX: this block could be a function
			refreshData := new(RefreshResponse)
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			err = json.Unmarshal(body, &refreshData)
			if err != nil {
				log.Printf("%T\n%s\n%#v\n", err, err, err)
				switch v := err.(type) {
				case *json.SyntaxError:
					fmt.Println(string(body[v.Offset-40 : v.Offset]))
				}
			}

			// Connect to the imap server
			c := Dial(Addr)

			if c.Caps["STARTTLS"] {
				sync.ReportOK(c.StartTLS(nil))
			}

			if c.Caps["ID"] {
				sync.ReportOK(c.ID("name", "goimap"))
			}

			sync.ReportOK(c.Noop())
			sync.ReportOK(c.Auth(imap.XOAUTH2Auth(acc.Username, refreshData.AccessToken)))

			// get the list of all folders for this account
			folders := acc.Folders(c, "")

			for _, f := range folders {
				// No reason to sync folders that can't be selected
				if f.Noselect {
					continue
				}

				// Don't fetch some common folders (Trash, Spam, etc..)
				isIgnored := false
				for _, ignored := range ignoredFolders {
					if f.Name == ignored {
						isIgnored = true
						break
					}
				}
				if isIgnored {
					continue
				}

				// get folder from the the DB
				syncFolder, err := f.SyncFolder(db)
				if err != nil {
					log.Fatal("Failed to get synced folder: ", err)
				}

				if syncFolder == nil {
					// create a new folder in the db based on the one from IMAP
					syncFolder, err = f.NewSyncFolder(db)
					if err != nil {
						log.Fatal("Failed to create sync folder: ", err)
					}
				}

				if syncFolder.Deleted { // if the folder is deleted then just skip it
					continue
				}

				syncFolder.Account = acc

				// get all messages for this folder - only their metadata without the body
				imapMessages, err := f.Messages(c)
				if err != nil {
					log.Fatal("Failed to fetch imap messages: ", err)
				}

				// check if these messsages already exist and add them to this folder if they are not added already or
				// remove them if they are not present any more
				err = syncFolder.HandleExisting(db, imapMessages)
				if err != nil {
					log.Fatal("Failed to handle existing messages: ", err)
				}

				// get all messages from the db for this folder
				syncedMessages, err := syncFolder.Messages(db)
				if err != nil {
					log.Fatal("Failed to get synced messages: ", err)
				}

				// filter messages that we don't need to fetch
				fetchMessages := sync.FilterMessages(syncedMessages, imapMessages)

				// finally get only the new messages
				syncFolder.SyncMessages(c, db, fetchMessages)

			}
			_, err = db.Exec(`UPDATE sync_account SET sync_datetime = (now() at time zone 'utc') WHERE id = $1`, acc.Id)
			if err != nil {
				panic(err)
			}

			// finally logout
			sync.ReportOK(c.Logout(time.Second * 30))
		}
		time.Sleep(SYNC_INTERVAL)
	}
}

func Dial(addr string) (c *imap.Client) {
	var err error
	if strings.HasSuffix(addr, ":993") {
		c, err = imap.DialTLS(addr, nil)
	} else {
		c, err = imap.Dial(addr)
	}
	if err != nil {
		panic(err)
	}
	return c
}
