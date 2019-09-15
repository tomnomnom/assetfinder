package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
		output = append(output, res.Name)
	}
	return output, nil
}
