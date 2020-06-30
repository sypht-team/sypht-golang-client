package sypht

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var gracePeriod = 10

// RefreshToken refreshes api token and updates Sypht client instance
func (s *Client) RefreshToken() (token string, err error) {
	var req *http.Request
	var payload *strings.Reader
	if strings.Contains(s.config.authURL, "/oauth2") {
		payload = strings.NewReader("client_id=" + s.config.clientID + "&grant_type=client_credentials")
		basicAuthSlug := base64.StdEncoding.EncodeToString([]byte(strings.Join([]string{s.config.clientID, s.config.clientSecret}, ":")))
		req, err = http.NewRequest("POST", s.config.authURL, payload)
		if err != nil {
			log.Println(err)
			return
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Authorization", "Basic "+basicAuthSlug)
	} else {
		payload = strings.NewReader(strings.Join([]string{"{\n \"client_id\":\"",
			s.config.clientID,
			"\",\n \"client_secret\":\"",
			s.config.clientSecret,
			"\",\n \"audience\":\"https://api.sypht.com\",\n \"grant_type\":\"client_credentials\" \n}"}, ""))
		req, err = http.NewRequest("POST", s.config.authURL, payload)
		if err != nil {
			log.Println(err)
			return
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var dat map[string]interface{}
	json.Unmarshal(body, &dat)
	token = string(dat["access_token"].(string))
	s.mutex.Lock()
	s.apiToken = token
	s.tokenUpdatedAt = time.Now()
	s.mutex.Unlock()
	return
}

func (s *Client) getToken() (token string) {
	s.mutex.RLock()
	if time.Since(s.tokenUpdatedAt) >= time.Hour*24-time.Minute*time.Duration(gracePeriod) {
		s.mutex.RUnlock()
		token, err := s.RefreshToken()
		if err != nil {
			log.Printf("Error refreshing token : %v", err)
		}
		return token
	}
	s.mutex.RLock()
	token = s.apiToken
	s.mutex.RUnlock()
	return
}
