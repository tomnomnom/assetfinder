package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func fetchJldc(domain string) ([]string, error) {

	fetchURL := fmt.Sprintf("https://jldc.me/anubis/subdomains/%s", domain)
	resp, err := http.Get(fetchURL)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	var body string
	_, err = fmt.Fscanln(resp.Body, &body)
	if err != nil {
		return []string{}, err
	}

	var lines []string
	err = json.Unmarshal([]byte(body), &lines)
	if err != nil {
		return []string{}, err
	}

	return lines, nil
}
