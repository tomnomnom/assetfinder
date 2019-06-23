package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func fetchHackerTarget(domain string) ([]string, error) {
	out := make([]string, 0)

	raw, err := httpGet(
		fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain),
	)
	if err != nil {
		return out, err
	}

	sc := bufio.NewScanner(bytes.NewReader(raw))
	for sc.Scan() {
		parts := strings.SplitN(sc.Text(), ",", 2)
		if len(parts) != 2 {
			continue
		}

		out = append(out, parts[0])
	}

	return out, sc.Err()
}
