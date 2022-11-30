package main

import (
	"fmt"
)

func fetchCertSpotter(domain string) ([]string, error) {
	out := make([]string, 0)

	fetchURL := fmt.Sprintf("https://api.certspotter.com/v1/issuances?domain=%s&expand=dns_names&include_subdomains=true&match_wildcards=true", domain)

	wrapper := []struct {
		DNSNames []string `json:"dns_names"`
	}{}
	err := fetchJSON(fetchURL, &wrapper)
	if err != nil {
		return out, err
	}

	for _, w := range wrapper {
		out = append(out, w.DNSNames...)
	}

	return out, nil
}
