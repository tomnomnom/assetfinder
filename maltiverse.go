package main

import (
	"fmt"
)

type Maltiverse struct {
	Hits struct {
		Hits []struct {
			ID     string      `json:"_id"`
			Index  string      `json:"_index"`
			Score  interface{} `json:"_score"`
			Source struct {
				Hostname string `json:"hostname"`
			} `json:"_source,omitempty"`
		}
	}
}

func fetchMaltiverse(domain string) ([]string, error) {
	out := make([]string, 0)

	fetchURL := fmt.Sprintf("https://api.maltiverse.com/search?format=json&from=0&query=hostname.keyword:*.%s&size=10000&sort=creation_time_desc", domain)

	var pdns Maltiverse

	err := fetchJSON(fetchURL, &pdns)
	if err != nil {
		return out, err
	}
	for _, dns := range pdns.Hits.Hits {
		out = append(out, dns.Source.Hostname)
	}

	return out, nil
}
