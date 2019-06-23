package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	var subsOnly bool
	flag.BoolVar(&subsOnly, "subs-only", false, "Only incluse subdomains of search domain")
	flag.Parse()

	domain := flag.Arg(0)
	if domain == "" {
		fmt.Println("no domain specified")
		return
	}
	domain = strings.ToLower(domain)

	sources := []fetchFn{
		fetchCertSpotter,
		fetchHackerTarget,
		fetchThreatCrowd,
		fetchCrtSh,
		fetchFacebook,
		//fetchWayback, // A little too slow :(
		fetchVirusTotal,
	}

	out := make(chan string)
	var wg sync.WaitGroup

	// call each of the source workers in a goroutine
	for _, source := range sources {
		wg.Add(1)
		fn := source

		go func() {
			defer wg.Done()

			names, err := fn(domain)

			if err != nil {
				fmt.Fprintf(os.Stderr, "err: %s\n", err)
				return
			}

			for _, n := range names {
				out <- n
			}
		}()
	}

	// close the output channel when all the workers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	// track what we've already printed to avoid duplicates
	printed := make(map[string]bool)

	for n := range out {
		n = cleanDomain(n)
		if _, ok := printed[n]; ok {
			continue
		}
		if subsOnly && !strings.HasSuffix(n, domain) {
			continue
		}
		fmt.Println(n)
		printed[n] = true
	}
}

type fetchFn func(string) ([]string, error)

func httpGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	raw, err := ioutil.ReadAll(res.Body)

	res.Body.Close()
	if err != nil {
		return []byte{}, err
	}

	return raw, nil
}

func cleanDomain(d string) string {
	d = strings.ToLower(d)

	// no idea what this is, but we can't clean it ¯\_(ツ)_/¯
	if len(d) < 2 {
		return d
	}

	if d[0] == '*' || d[0] == '%' {
		d = d[1:]
	}

	if d[0] == '.' {
		d = d[1:]
	}

	return d

}

func fetchJSON(url string, wrapper interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	return dec.Decode(wrapper)
}
