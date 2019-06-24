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

	// Start querying the Spyse API from page 1
	page := 1
	out := make([]string, 0)

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
			return out, err
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
	
	return out, nil
}
