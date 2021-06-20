package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gomodule/oauth1/oauth"
	"github.com/joeshaw/envdecode"
)

// we want to keep tracking the state of a connection
var (
	cnx net.Conn

	reader io.ReadCloser

	authClient *oauth.Client

	creds *oauth.Credentials

	authSetupOnce sync.Once

	client *http.Client
)

const TIMEOUT time.Duration = 3 * time.Second

func dial(network, addr string) (net.Conn, error) {

	if cnx != nil {

		cnx.Close()

		cnx = nil

	}

	// establish a new connection and set timeout 3 sec
	conn, err := net.DialTimeout(network, addr, TIMEOUT)

	if err != nil {

		return nil, err

	}

	cnx = conn

	return conn, nil
}

// this funciton will be invoked periodically because we want to
// reload hashtag options from our db at regular intv
func close_cnx() {

	if cnx != nil {

		cnx.Close()

	}

	if reader != nil {

		reader.Close()

	}

}

func setupTwitterAuth() {
	var ts struct {
		ConsumerKey    string `env:"TWITTER_KEY,required"`
		ConsumerSecret string `env:"TWITTER_SECRET, required"`
		AccessToken    string `env:"TWITTER_ACCESS_TOKEN, required"`
		AccessSecret   string `env:"TWITTER_ACCESS_SECRET, required"`
	}

	if err := envdecode.Decode(&ts); err != nil {
		log.Fatalln(err)
	}

	creds = &oauth.Credentials{
		Token:  ts.AccessToken,
		Secret: ts.AccessSecret,
	}

	authClient = &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  ts.ConsumerKey,
			Secret: ts.ConsumerSecret,
		},
	}
}

func makeRequest(req *http.Request, params url.Values) (*http.Response, error) {
	authSetupOnce.Do(func() {
		setupTwitterAuth()
		client = &http.Client{
			Transport: &http.Transport{
				Dial: dial,
			},
		}
	})

	form := params.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(form)))
	req.Header.Set("Authorization", authClient.AuthorizationHeader(creds, "POST", req.URL, params))

	return client.Do(req)
}
