package main

import (
	"fmt"
)

type PassiveDNS struct {
	PassiveDNS []struct {
		Hostname string `json:"hostname"`
	} `json:"passive_dns"`
}

func fetchAlienVault(domain string) ([]string, error) {
	out := make([]string, 0)

	fetchURL := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/passive_dns", domain)

	var pdns PassiveDNS

	err := fetchJSON(fetchURL, &pdns)
	if err != nil {
		return out, err
	}
	for _, dns := range pdns.PassiveDNS {
		out = append(out, dns.Hostname)
	}

	return out, nil
}
