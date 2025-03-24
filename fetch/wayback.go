package fetch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetWaybackURLs(domain string, noSubs bool) ([]Wurl, error) {
	subsWildcard := "*."
	if noSubs {
		subsWildcard = ""
	}

	res, err := http.Get(fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=%s%s/*&output=json&collapse=urlkey", subsWildcard, domain))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var wrapper [][]string
	if err := json.Unmarshal(raw, &wrapper); err != nil {
		return nil, err
	}

	out := make([]Wurl, 0, len(wrapper))
	skip := true
	for _, urls := range wrapper {
		if skip {
			skip = false
			continue
		}
		out = append(out, Wurl{Date: urls[1], Url: urls[2]})
	}

	return out, nil
}