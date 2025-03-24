package ssti

import (
	"crypto/tls" 
	"fmt"
	"io" 
	"net" 
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/eikehacker1/sstighost/config"
)

type SSTIPayload struct {
	Payload  string
	Expected *regexp.Regexp
}

var SSTIPayloads = []SSTIPayload{
	// Python (Jinja2, Flask, etc.)
	{"{{ 74521 * 9 }}", regexp.MustCompile(`670689`)},

	// Ruby (ERB, Slim, etc.)
	{"<%= 74521 * 9 %>", regexp.MustCompile(`670689`)},

	// PHP (Twig, Smarty, etc.)
	{"{{ 74521 * 9 }}", regexp.MustCompile(`670689`)},

	// JavaScript (Node.js, EJS, etc.)
	{"<%= 74521 * 9 %>", regexp.MustCompile(`670689`)},
}

func SSTI(urlt string, payload string, expected *regexp.Regexp, proxy string, onlypoc bool) string {
	client := createClient(proxy)
	u, err := url.Parse(urlt)
	if err != nil {
		return "ERROR"
	}

	q := u.Query()
	for x := range q {
		q.Set(x, payload)
	}
	u.RawQuery = q.Encode()
	urlt = u.String()

	res, err := http.NewRequest("GET", urlt, nil)
	if err != nil {
		return "ERROR"
	}
	res.Header.Set("Connection", "close")
	for _, v := range config.Headers {
		s := strings.SplitN(v, ":", 2)
		res.Header.Set(s[0], s[1])
	}

	resp, err := client.Do(res)
	if err != nil {
		return "ERROR"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "ERROR"
	}

	page := string(body)
	if expected.MatchString(page) {
		if onlypoc {
			return fmt.Sprintf("\033[1;31m[VULNERABLE] %s (Payload: %s)\033[0;0m", urlt, payload)
		}
		return fmt.Sprintf("\033[1;31m[VULNERABLE] %s (Payload: %s)\033[0;0m", urlt, payload)
	} else if !onlypoc {
		return fmt.Sprintf("\033[1;30m[NOT VULNERABLE] %s (Payload: %s)\033[0;0m", urlt, payload)
	}
	return "ERROR"
}

func createClient(proxy string) *http.Client {
	trans := &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: time.Second,
		}).DialContext,
	}

	if proxy != "0" {
		if p, err := url.Parse(proxy); err == nil {
			trans.Proxy = http.ProxyURL(p)
		}
	}

	return &http.Client{
		Transport: trans,
		Timeout:   3 * time.Second,
	}
}