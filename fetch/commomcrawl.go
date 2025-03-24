package fetch

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetCommonCrawlURLs(domain string, noSubs bool) ([]Wurl, error) {
	subsWildcard := "*."
	if noSubs {
		subsWildcard = ""
	}

	res, err := http.Get(fmt.Sprintf("http://index.commoncrawl.org/CC-MAIN-2018-22-index?url=%s%s/*&output=json", subsWildcard, domain))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	sc := bufio.NewScanner(res.Body)
	out := make([]Wurl, 0)

	for sc.Scan() {
		wrapper := struct {
			URL       string `json:"url"`
			Timestamp string `json:"timestamp"`
		}{}
		if err := json.Unmarshal(sc.Bytes(), &wrapper); err != nil {
			continue
		}
		out = append(out, Wurl{Date: wrapper.Timestamp, Url: wrapper.URL})
	}

	return out, nil
}