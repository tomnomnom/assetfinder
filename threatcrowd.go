package main

import (
	"fmt"
)

func fetchThreatCrowd(domain string) ([]string, error) {
	out := make([]string, 0)

	fetchURL := fmt.Sprintf("https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=%s", domain)

	wrapper := struct {
		Subdomains []string `json:"subdomains"`
	}{}
	err := fetchJSON(fetchURL, &wrapper)
	if err != nil {
		return out, err
	}

	out = append(out, wrapper.Subdomains...)

	return out, nil
}
