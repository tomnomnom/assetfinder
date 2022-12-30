package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var subsOnly bool
var wayBack bool

type Result struct {
	n      string
	domain string
}

func main() {

	flag.BoolVar(&subsOnly, "subs-only", false, "Only include subdomains of search domain")
	flag.BoolVar(&wayBack, "w", false, "Include a search of the Wayback Machine...slow")
	flag.Parse()

	var domains io.Reader
	domains = os.Stdin

	domain := flag.Arg(0)
	if domain != "" {
		domains = strings.NewReader(domain)
	}

	sources := []fetchFn{
		fetchCertSpotter,
		fetchHackerTarget,
		fetchThreatCrowd,
		fetchCrtSh,
		fetchFacebook,
		// fetchWayback, // A little too slow :(
		fetchVirusTotal,
		// fetchFindSubDomains,
		fetchUrlscan,
		fetchBufferOverrun,
		fetchRiskIq,
		fetchRiddler,
		fetchDnsSpy,
		fetchAlienVault,
		fetchMaltiverse,
		fetchArquivo,
		fetchDnsHistory,
		fetchJldc,
	}

	// optional add in wayback
	if wayBack {
		sources = append(sources, fetchWayback)
	}

	out := make(chan Result)

	var wg sync.WaitGroup

	sc := bufio.NewScanner(domains)
	rl := newRateLimiter(time.Second)

	for sc.Scan() {
		domain := strings.ToLower(sc.Text())

		// call each of the source workers in a goroutine
		for _, source := range sources {
			wg.Add(1)
			fn := source

			go func() {
				defer wg.Done()

				rl.Block(fmt.Sprintf("%#v", fn))
				names, err := fn(domain)

				if err != nil {
					return
				}

				for _, n := range names {

					n = cleanDomain(n)

					res := new(Result)
					res.n = n
					res.domain = domain

					out <- *res
				}

			}()
		}
	}

	if err := sc.Err(); err != nil {
		fmt.Println(err)
	}

	// close the output channel when all the workers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	// track what we've already printed to avoid duplicates
	printed := make(map[string]bool)

	for res := range out {

		if _, ok := printed[res.n]; ok {
			continue
		}

		/*
			moved this check to here as there appeared to be
			an issue where non subdomains were being returned
			if this check was in the go routine
		*/
		if subsOnly && !strings.HasSuffix(res.n, res.domain) {
			continue
		}

		fmt.Println(res.n)
		printed[res.n] = true
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

	if strings.HasPrefix(d, "*") || strings.HasPrefix(d, "%") {
		d = d[1:]
	}

	if strings.HasPrefix(d, ".") {
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
