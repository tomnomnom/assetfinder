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
var verbose bool
var wayBack bool

type Result struct {
	n      string
	domain string
	src    string
}

type fetchFn func(string) ([]string, error)

type fetchSource struct {
	name string
	fn   fetchFn
}

func main() {

	flag.BoolVar(&subsOnly, "subs-only", false, "Only include subdomains of search domain")
	flag.BoolVar(&wayBack, "w", false, "Include a search of the Wayback Machine...slow")
	flag.BoolVar(&verbose, "v", false, "Be verbose, and show the source of the subdomain in the output")
	flag.Parse()

	var domains io.Reader
	domains = os.Stdin

	domain := flag.Arg(0)
	if domain != "" {
		domains = strings.NewReader(domain)
	}

	sources := []fetchSource{
		{"Certspotter", fetchCertSpotter},
		{"HackerTarget", fetchHackerTarget},
		{"ThreatCrowd", fetchThreatCrowd},
		{"crt.sh", fetchCrtSh},
		{"Facebook", fetchFacebook},
		{"VirusTotal", fetchVirusTotal},
		{"UrlScan", fetchUrlscan},
		{"BufferOverrun", fetchBufferOverrun},
		{"RiskIq", fetchRiskIq},
		{"Riddler", fetchRiddler},
		{"DnsSpy", fetchDnsSpy},
		{"AlienVault", fetchAlienVault},
		{"Maltiverse", fetchMaltiverse},
		{"Arquivo", fetchArquivo},
		{"DnsHistory", fetchDnsHistory},
		{"Jldc", fetchJldc},
	}

	// optional add in wayback
	if wayBack {
		sources = append(sources, fetchSource{"Wayback", fetchWayback})
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
			fn := source.fn

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
					res.src = source.name // Use the name of the source

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

		if verbose {
			fmt.Printf("%s,%s\n", res.src, res.n)

		} else {
			fmt.Println(res.n)
		}

		printed[res.n] = true
	}
}

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
