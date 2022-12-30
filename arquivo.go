package main

import (
	"fmt"
	"net/url"
)

type Arquivo struct {
	ResponseItems []struct {
		OriginalURL string `json:"originalURL"`
	} `json:"response_items"`
}

func fetchArquivo(domain string) ([]string, error) {
	out := make([]string, 0)

	fetchURL := fmt.Sprintf("https://arquivo.pt/textsearch?q=%s", domain)

	var pdns Arquivo

	err := fetchJSON(fetchURL, &pdns)
	if err != nil {
		return out, err
	}
	for _, dns := range pdns.ResponseItems {

		uri, err := url.Parse(dns.OriginalURL)
		if err != nil {
			continue
		}
		out = append(out, uri.Hostname())
	}

	return out, nil
}
