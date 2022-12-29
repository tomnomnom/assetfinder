package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func fetchDnsSpy(domain string) ([]string, error) {
	resp, err := http.Get(
		fmt.Sprintf("https://dnsspy.io/scan/%s", domain),
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

	re := regexp.MustCompile(`\s+`)
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {

		table_text := s.Find("td").Eq(1).Text()
		table_text = re.ReplaceAllString(table_text, "")

		if strings.Contains(table_text, domain) {
			output = append(output, table_text)
		}

	})

	return output, nil
}
