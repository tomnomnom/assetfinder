package main

import (
	"fmt"
	"os"
)

func fetchCertSpotter(domain string) ([]string, error) {
	out := make([]string, 0)

	apiKey := os.Getenv("CERTSPOTTER_API_KEY")
	fetchURL := fmt.Sprintf("https://%s@api.certspotter.com/v1/issuances?domain=%s&expand=dns_names&include_subdomains=true", apiKey, domain)	
	
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
