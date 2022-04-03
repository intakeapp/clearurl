package clearurl

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
)

//go:embed data.min.json
var rulesData string

var rulesMap = map[string]map[string]*Provider{}

func init() {
	if err := json.Unmarshal([]byte(rulesData), &rulesMap); err != nil {
		panic(err)
	}
}

type Provider struct {
	Urlpattern       string         `json:"urlPattern"`
	UrlpatternRegexp *regexp.Regexp `json:"-"`
	Completeprovider bool           `json:"completeProvider"`
	Rules            []string       `json:"rules"`
	Exceptions       []string       `json:"exceptions"`
	Redirections     []string       `json:"redirections"`
}

func (p *Provider) init() error {
	var err error
	if p.UrlpatternRegexp == nil {
		p.UrlpatternRegexp, err = regexp.Compile(p.Urlpattern)
		if err != nil {
			return err
		}
	}

	return nil
}

type Handler struct {
	providers map[string]*Provider
	global    *Provider
}

func Init() (*Handler, error) {
	h := &Handler{
		providers: rulesMap["providers"],
	}

	for _, p := range h.providers {
		if err := p.init(); err != nil {
			return nil, err
		}
	}

	h.global = h.providers["globalRules"]
	delete(h.providers, "globalRules")

	return h, nil
}

func (h *Handler) Clear(url string) (string, error) {
	for k, p := range h.providers {
		if p.UrlpatternRegexp.MatchString(url) {
			fmt.Println(k, p.Urlpattern)
		}
	}

	return url, nil
}
