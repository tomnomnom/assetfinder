package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func fetchDnsHistory(domain string) ([]string, error) {
	resp, err := http.Get(
		fmt.Sprintf("https://dnshistory.org/subdomains/1/%s", domain),
	)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return []string{}, err
	}

	output := make([]string, 0)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []string{}, err
	}

	mainarea := doc.Find("#mainarea")

	mainarea.Find("a").Each(func(i int, s *goquery.Selection) {
		hrefText := s.Text()
		output = append(output, hrefText)
	})

	return output, nil
}
