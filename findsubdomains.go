package main

import (
	"fmt"
	"os"
)

var apiToken = os.Getenv("SPYSE_API_TOKEN")

func callSubdomainsAggregateEndpoint(domain string) []string {
	out := make([]string, 0)

	fetchURL := fmt.Sprintf(
		"https://api.spyse.com/v1/subdomains-aggregate?api_token=%s&domain=%s",
		apiToken, domain,
	)

	type Cidr struct {
		Results []struct {
			Data struct {
				Domains []string `json:"domains"`
			} `json:"data"`
		} `json:"results"`
	}

	type Cidrs struct {
		Cidr16, Cidr24 Cidr
	}

	wrapper := struct {
		Cidrs Cidrs `json:"cidr"`
	}{}

	err := fetchJSON(fetchURL, &wrapper)

	if err != nil {
		// Fail silently
		return []string{}
	}

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

	return out
}

/**

 */
func callSubdomainsEndpoint(domain string) []string {
	out := make([]string, 0)

	// Start querying the Spyse API from page 1
	page := 1

	for {
		wrapper := struct {
			Records []struct {
				Domain string `json:"domain"`
			} `json:"records"`
		}{}

		fetchURL := fmt.Sprintf(
			"https://api.spyse.com/v1/subdomains?api_token=%s&domain=%s&page=%d",
			apiToken, domain, page,
		)

		err := fetchJSON(fetchURL, &wrapper)
		if err != nil {
			// Fail silently, by returning what we got so far
			return out
		}

		// The API does not respond with any paging, nor does it give us any idea of
		// the total amount of domains, so we just have to keep asking for a new page until
		// the returned `records` array is empty
		// NOTE: The free tier always gives you the first page for free, and you get "25 unlimited search requests"
		if len(wrapper.Records) == 0 {
			break
		}

		for _, record := range wrapper.Records {
			out = append(out, record.Domain)
		}

		page++
	}

	return out
}

func fetchFindSubDomains(domain string) ([]string, error) {

	out := make([]string, 0)

	apiToken := os.Getenv("SPYSE_API_TOKEN")
	if apiToken == "" {
		// Must have an API token
		return []string{}, nil
	}

	// The Subdomains-Aggregate endpoint returns some, but not all available domains
	out = append(out, callSubdomainsAggregateEndpoint(domain)...)

	// The Subdomains endpoint only guarantees the first 30 domains, the rest needs credit at Spyze
	out = append(out, callSubdomainsEndpoint(domain)...)

	return out, nil
}
