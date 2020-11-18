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

func main() {
	var subsOnly bool
    override := make(map[string]*bool)

	flag.BoolVar(&subsOnly, "subs-only", false, "Only include subdomains of search domain")
    // Flag for all Modules... Sorry Little messy
    override["fetchCertSpotter"]    = flag.Bool("cert", true, "<Boolean> (Default = True) Toggle CertSpotter source")
    override["fetchHackerTarget"]   = flag.Bool("h", true, "<Boolean> (Default = True) Toggle CertSpotter source")
    override["fetchThreatCrowd"]    = flag.Bool("t", true, "<Boolean> (Default = True) Toggle ThreatCrowd source")
    override["fetchCrtSh"]          = flag.Bool("crt", true, "<Boolean> (Default = True) Toggle CrtSh source")
    override["fetchDNSDB"]          = flag.Bool("d", false, "<Boolean> (Default = True) Toggle DNSDB Enterprise source")
    override["fetchDNSDBCommunity"] = flag.Bool("dC", true, "<Boolean> (Default = True) Toggle DNSDB Community source")
    override["fetchFacebook"]       = flag.Bool("f", true, "<Boolean> (Default = True) Toggle Facebook source")
    override["fetchPassiveTotal"]   = flag.Bool("p", true, "<Boolean> (Default = True) Toggle PassiveTotal source")
    override["fetchWayback"]        = flag.Bool("w", false, "<Boolean> (Default = False)  Toggle Wayback source")
    override["fetchVirusTotal"]     = flag.Bool("v", true, "<Boolean> (Default = True) Toggle VirusTotal source")
    override["fetchFindSubDomains"] = flag.Bool("fs", true, "<Boolean> (Default = True) Toggle FindSubDomains source")
    override["fetchUrlscan"]        = flag.Bool("u", true, "<Boolean> (Default = True) Toggle Urlscan source")
    override["fetchBufferOverrun"]  = flag.Bool("b", true, "<Boolean> (Default = True) Toggle BufferOverrun source")
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
		//fetchWayback, // A little too slow :(
        fetchDNSDBCommunity, // Per Hour only 4 domains as Community Edition has limit of 25 query per hour and each domain takes 6 request for complete result for community edition See=>dnsdbCommunity.go Line 13
        fetchPassiveTotal,
        //fetchDNSDB,
		fetchVirusTotal,
		fetchFindSubDomains,
		fetchUrlscan,
		fetchBufferOverrun,
	}

    sources = toggleSources(override, sources)//Toogle Sources according flags
	out := make(chan string)
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
					//fmt.Fprintf(os.Stderr, "err: %s\n", err)
					return
				}

				for _, n := range names {
					n = cleanDomain(n)
					if subsOnly && !strings.HasSuffix(n, domain) {
						continue
					}
					out <- n
				}
			}()
		}
	}

	// close the output channel when all the workers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	// track what we've already printed to avoid duplicates
	printed := make(map[string]bool)

	for n := range out {
		if _, ok := printed[n]; ok {
			continue
		}
		printed[n] = true

		fmt.Println(n)
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

