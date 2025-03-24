package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func GetVirusTotalURLs(domain string, noSubs bool) ([]Wurl, error) {
	out := make([]Wurl, 0)

	apiKey := os.Getenv("VT_API_KEY")
	if apiKey == "" {
		return out, nil
	}

	resp, err := http.Get(fmt.Sprintf("https://www.virustotal.com/vtapi/v2/domain/report?apikey=%s&domain=%s", apiKey, domain))
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	wrapper := struct {
		URLs []struct {
			URL string `json:"url"`
		} `json:"detected_urls"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return out, err
	}

	for _, u := range wrapper.URLs {
		out = append(out, Wurl{Url: u.URL})
	}

	return out, nil
}