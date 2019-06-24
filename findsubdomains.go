package main

import (
	"fmt"
	"os"
)

func fetchFindSubDomains(domain string) ([]string, error) {

	apiToken := os.Getenv("SPYSE_API_TOKEN")
	if apiToken == "" {
		// Must have an API token
		return []string{}, nil
	}

	fetchURL := fmt.Sprintf(
		"https://api.spyse.com/v1/subdomains-aggregate?api_token=%s&domain=%s",
		apiToken, domain,
	)

	out := make([]string, 0)

	type Cidr struct {
		Results []struct {
			Data struct {
				Domains []string `json:"domains"`
			} `json:"data"`
		} `json:"results"`
		Count int `json:"count"`
	}

	type Cidrs struct {
		Cidr16, Cidr24 Cidr
	}

	wrapper := struct {
		Cidrs Cidrs `json:"cidr"`
	}{}
	
	
	err := fetchJSON(fetchURL, &wrapper)

	for _, result := range wrapper.Cidrs.Cidr16.Results {
		for _, domain := range result.Data.Domains {
			out = append(out, domain)
		}
	}
	for _, result := range wrapper.Cidrs.Cidr24.Results {
		for _, domain := range result.Data.Domains {
			out = append(out, domain)
		}
	}

	return out, err
}
