package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Wurl struct {
	Date string
	Url  string
}

type FetchFn func(string, bool) ([]Wurl, error)

func GetVersions(u string) ([]string, error) {
	out := make([]string, 0)

	resp, err := http.Get(fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=%s&output=json", u))
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	var r [][]string
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return out, err
	}

	first := true
	seen := make(map[string]bool)
	for _, s := range r {
		if first {
			first = false
			continue
		}
		if seen[s[5]] {
			continue
		}
		seen[s[5]] = true
		out = append(out, fmt.Sprintf("https://web.archive.org/web/%sif_/%s", s[1], s[2]))
	}

	return out, nil
}

func IsSubdomain(rawUrl, domain string) bool {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return false
	}
	return strings.ToLower(u.Hostname()) != strings.ToLower(domain)
}