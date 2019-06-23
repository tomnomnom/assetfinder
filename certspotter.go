package main

import (
	"fmt"
)

func fetchCertSpotter(domain string) ([]string, error) {
	out := make([]string, 0)

	fetchURL := fmt.Sprintf("https://certspotter.com/api/v0/certs?domain=%s", domain)

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
