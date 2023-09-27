package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type CrtShResult struct {
	Name string `json:"name_value"`
}

func fetchCrtSh(domain string) ([]string, error) {
	var results []CrtShResult

	resp, err := http.Get(
		fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain),
	)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	output := make([]string, 0)

	body, _ := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &results); err != nil {
		return []string{}, err
	}

	for _, res := range results {

		/*
			results from crt.sh could product multiple
			hits in name_value field which are separated by a \n
			this then prints these values out which don't
			go through cleanDomain, e.g.
			"*.qa-release.yhs.search.yahoo.com\n*.qa-trunk.yhs.search.yahoo.com"
			will return
			"qa-release.yhs.search.yahoo.com
			*.qa-trunk.yhs.search.yahoo.com"

			using strings.Fields separates on newline and iterating each result to return
			cures this issue
		*/
		s := strings.Fields(res.Name)
		for _, element := range s {
			output = append(output, element)
		}
	}
	return output, nil
}
