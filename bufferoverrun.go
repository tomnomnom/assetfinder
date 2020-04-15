package main

import (
	"fmt"
	"strings"
)

func fetchBufferOverrun(domain string) ([]string, error) {
	out := make([]string, 0)

	fetchURL := fmt.Sprintf("https://dns.bufferover.run/dns?q=.%s", domain)

	wrapper := struct {
		Records []string `json:"FDNS_A"`
	}{}
	err := fetchJSON(fetchURL, &wrapper)
	if err != nil {
		return out, err
	}

	for _, r := range wrapper.Records {
		parts := strings.SplitN(r, ",", 2)
		if len(parts) != 2 {
			continue
		}
		out = append(out, parts[1])
	}

	return out, nil
}
