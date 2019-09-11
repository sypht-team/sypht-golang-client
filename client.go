package sypht

import (
	"net/http"
	"os"
	"sync"
	"time"
)

/*
Client to the HTTP API of Sypht.
*/

// Client ...
type Client struct {
	httpClient     *http.Client
	config         *config
	apiToken       string
	tokenUpdatedAt time.Time
	mutex          sync.RWMutex
}

type config struct {
	clientID     string
	clientSecret string
	apiBaseURL   string
	authURL      string
}

var defaultTimeout = 30

// fieldSets const
const (
	Generic  = "\"sypht.generic\""
	Document = "\"sypht.document\""
	Invoice  = "\"sypht.invoice\""
	Bill     = "\"sypht.bill\""
	Bank     = "\"sypht.bank\""
)

//NewSyphtClient returns a Sypht client instance,
// default request timeout is set to 30 seconds, change it as needed
func NewSyphtClient(clientID string, clientSecret string, timeout *int) (client *Client, err error) {
	if timeout == nil || *timeout < 0 {
		timeout = &defaultTimeout
	}
	client = &Client{
		httpClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Timeout: time.Second * time.Duration(*timeout),
		},
		config: &config{
			clientID:     clientID,
			clientSecret: clientSecret,
			apiBaseURL:   "https://api.sypht.com",
			authURL:      "https://login.sypht.com/oauth/token",
		},
	}
	_, err = client.RefreshToken()
	return
}

//NewSyphtClientFromEnv same as NewSyphtClient except it reads creds from ENV
func NewSyphtClientFromEnv(timeout *int) (client *Client, err error) {
	client, err = NewSyphtClient(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), timeout)
	return
}
