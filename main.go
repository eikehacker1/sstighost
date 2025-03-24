package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"sstighost/config"
	"sstighost/fetch"
	"sstighost/ssti"
)

var (
	concurrency    int
	proxy          string
	poc            bool
	appendMode     bool
	ignorePath     bool
	dates          bool
	noSubs         bool
	getVersionsFlag bool
	endpoint       string 
)

func init() {
	flag.BoolVar(&appendMode, "a", false, "Append the value instead of replacing it")
	flag.BoolVar(&ignorePath, "ignore-path", false, "Ignore the path when considering what constitutes a duplicate")
	flag.BoolVar(&dates, "dates", false, "Show date of fetch in the first column")
	flag.BoolVar(&noSubs, "no-subs", false, "Don't include subdomains of the target domain")
	flag.BoolVar(&getVersionsFlag, "get-versions", false, "List URLs for crawled versions of input URL(s)")
	flag.IntVar(&concurrency, "c", 50, "Set concurrency")
	flag.StringVar(&proxy, "proxy", "0", "Send traffic to a proxy")
	flag.BoolVar(&poc, "only-poc", false, "Show only potentially vulnerable URLs")
	flag.StringVar(&endpoint, "e", "", "Endpoint to test for SSTI vulnerabilities") 
	flag.Var((*config.CustomHeaders)(&config.Headers), "headers", "Custom headers")
}

func main() {
	flag.Parse()

	
	if endpoint != "" {
		testEndpoint(endpoint)
		return
	}

	
	runCrawler()
}


func testEndpoint(endpoint string) {
	fmt.Printf("Testing endpoint: %s\n", endpoint)
	for _, payload := range ssti.SSTIPayloads {
		result := ssti.SSTI(endpoint, payload.Payload, payload.Expected, proxy, poc)
		if result != "ERROR" {
			fmt.Println(result)
		}
	}
}


func runCrawler() {
	var domains []string
	if flag.NArg() > 0 {
		domains = []string{flag.Arg(0)}
	} else {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			domains = append(domains, sc.Text())
		}
		if err := sc.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read input: %s\n", err)
		}
	}

	if getVersionsFlag {
		for _, u := range domains {
			versions, err := fetch.GetVersions(u)
			if err != nil {
				continue
			}
			fmt.Println(strings.Join(versions, "\n"))
		}
		return
	}

	fetchFns := []fetch.FetchFn{
		fetch.GetWaybackURLs,
		fetch.GetCommonCrawlURLs,
		fetch.GetVirusTotalURLs,
	}

	seen := make(map[string]bool)
	var wg sync.WaitGroup
	urls := make(chan string)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range urls {
				for _, payload := range ssti.SSTIPayloads {
					result := ssti.SSTI(u, payload.Payload, payload.Expected, proxy, poc)
					if result != "ERROR" {
						fmt.Println(result)
					}
				}
			}
		}()
	}

	for _, domain := range domains {
		var fetchWg sync.WaitGroup
		wurls := make(chan fetch.Wurl)

		for _, fn := range fetchFns {
			fetchWg.Add(1)
			go func(f fetch.FetchFn) {
				defer fetchWg.Done()
				resp, err := f(domain, noSubs)
				if err != nil {
					return
				}
				for _, r := range resp {
					if noSubs && fetch.IsSubdomain(r.Url, domain) {
						continue
					}
					wurls <- r
				}
			}(fn)
		}

		go func() {
			fetchWg.Wait()
			close(wurls)
		}()

		for w := range wurls {
			if _, ok := seen[w.Url]; ok {
				continue
			}
			seen[w.Url] = true

			if dates {
				d, err := time.Parse("20060102150405", w.Date)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to parse date [%s] for URL [%s]\n", w.Date, w.Url)
				}
				fmt.Printf("%s %s\n", d.Format(time.RFC3339), w.Url)
			} else {
				urls <- w.Url
			}
		}
	}

	close(urls)
	wg.Wait()
}