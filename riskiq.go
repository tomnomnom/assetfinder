package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	base_url = "https://api.riskiq.net/pt"
	path     = "/v2/enrichment/subdomains"
)

type Response struct {
	PrimaryDomain string   `json:"primaryDomain"`
	Success       bool     `json:"success"`
	QueryValue    string   `json:"queryValue"`
	Subdomains    []string `json:"subdomains"`
}

func fetchRiskIq(domain string) ([]string, error) {

	username := os.Getenv("RISK_IQ_USER")
	password := os.Getenv("RISK_IQ_KEY")

	// fail gracefully if there's no user API
	if username == "" || password == "" {
		return []string{}, nil
	}

	output := make([]string, 0)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// riskiq has a weird thing of sending json body data in a GET
	var payload = []byte(fmt.Sprintf(`{"query":"%v"}`, domain))

	req, err := http.NewRequest("GET", base_url+path, bytes.NewBuffer(payload))
	if err != nil {
		return []string{}, err
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	var response Response
	json.Unmarshal(body, &response)

	for _, subdomain := range response.Subdomains {
		output = append(output, subdomain+"."+domain)
	}
	return output, nil
}
