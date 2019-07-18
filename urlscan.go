package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func fetchUrlscan(domain string) ([]string, error) {
	resp, err := http.Get(
		fmt.Sprintf("https://urlscan.io/api/v1/search/?q=domain:%s", domain),
	)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	output := make([]string, 0)

	dec := json.NewDecoder(resp.Body)

	wrapper := struct {
		Results []struct {
			Task struct {
				URL string `json:"url"`
			} `json:"task"`

			Page struct {
				URL string `json:"url"`
			} `json:"page"`
		} `json:"results"`
	}{}

	err = dec.Decode(&wrapper)
	if err != nil {
		return []string{}, err
	}

	for _, r := range wrapper.Results {
		u, err := url.Parse(r.Task.URL)
		if err != nil {
			continue
		}

		output = append(output, u.Hostname())
	}

	for _, r := range wrapper.Results {
		u, err := url.Parse(r.Page.URL)
		if err != nil {
			continue
		}

		output = append(output, u.Hostname())
	}

	return output, nil
}
