package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
    "github.com/spiral-sec/assetfinder/scanner"
)

func main() {
	var subsOnly bool
	flag.BoolVar(&subsOnly, "subs-only", false, "Only include subdomains of search domain")
	flag.Parse()

	var domains io.Reader
	domains = os.Stdin

	domain := flag.Arg(0)
	if domain != "" {
		domains = strings.NewReader(domain)
	}

	sources := []fetchFn{
		scanner.CertSpotter,
		scanner.HackerTarget,
		scanner.ThreatCrowd,
		scanner.CrtSh,
		scanner.Facebook,
		//scanner.Wayback, // A little too slow :(
		scanner.VirusTotal,
		scanner.FindSubDomains,
		scanner.Urlscan,
		scanner.BufferOverrun,
	}

	out := make(chan string)
	var wg sync.WaitGroup

	sc := bufio.NewScanner(domains)
	rl := scanner.NewRateLimiter(time.Second)

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
					n = scanner.CleanDomain(n)
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

