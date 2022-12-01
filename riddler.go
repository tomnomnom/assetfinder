package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var RIDDLER_BASE_URL = "https://riddler.io"

type Auth struct {
	Meta struct {
		Code int `json:"code"`
	} `json:"meta"`
	Response struct {
		User struct {
			AuthenticationToken string `json:"authentication_token"`
			ID                  string `json:"id"`
		} `json:"user"`
	} `json:"response"`
}

type RiddlerResponse []struct {
	Keywords     []interface{} `json:"keywords"`
	Host         string        `json:"host"`
	Addr         string        `json:"addr"`
	CountryCode  string        `json:"country_code"`
	LastCrawling string        `json:"last_crawling"`
}

func fetchRiddler(domain string) ([]string, error) {
	email := os.Getenv("RIDDLER_EMAIL")
	password := os.Getenv("RIDDLER_PASS")
	if email == "" || password == "" {
		return []string{}, nil
	}

	accessToken, err := riddlerAuth(email, password)
	if err != nil {
		return []string{}, err
	}

	domains, err := getRiddlerHosts(accessToken, domain)
	if err != nil {
		return []string{}, err
	}

	return domains, nil
}

func getRiddlerHosts(accessToken string, domain string) ([]string, error) {

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	post_data, _ := json.Marshal(map[string]string{
		"query":  "pld:" + domain,
		"output": "host",
	})
	reqBody := bytes.NewBuffer(post_data)
	req, err := http.NewRequest("POST", RIDDLER_BASE_URL+"/api/search", reqBody)
	if err != nil {
		return []string{}, err
	}
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authentication-Token", accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	var riddlerResponse RiddlerResponse
	json.Unmarshal(body, &riddlerResponse)

	out := make([]string, 0)

	for _, i := range riddlerResponse {
		out = append(out, i.Host)
	}

	return out, nil
}

func riddlerAuth(email string, password string) (string, error) {

	postBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(RIDDLER_BASE_URL+"/auth/login", "application/json", responseBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var auth Auth
	json.Unmarshal(body, &auth)

	return auth.Response.User.AuthenticationToken, nil

}
