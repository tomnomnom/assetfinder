package main

import (
	"fmt"
	"os"
)

func fetchVirusTotal(domain string) ([]string, error) {

	apiKey := os.Getenv("VT_API_KEY")
	if apiKey == "" {
		// swallow not having an API key, just
		// don't fetch
		return []string{}, nil
	}

	fetchURL := fmt.Sprintf(
		"https://www.virustotal.com/vtapi/v2/domain/report?domain=%s&apikey=%s",
		domain, apiKey,
	)

	wrapper := struct {
		Subdomains []string `json:"subdomains"`
	}{}

	err := fetchJSON(fetchURL, &wrapper)
	return wrapper.Subdomains, err
}
