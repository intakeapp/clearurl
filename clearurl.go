package clearurl

import (
	_ "embed"
	"encoding/json"
	"net/url"
	"regexp"
	"strings"
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
	Name              string   `json:"-"`
	UrlPattern        string   `json:"urlPattern"`
	CompleteProvider  bool     `json:"completeProvider"`
	Rules             []string `json:"rules"`
	RawRules          []string `json:"rawRules"`
	ReferralMarketing []string `json:"referralMarketing"`
	Exceptions        []string `json:"exceptions"`
	Redirections      []string `json:"redirections"`

	urlPatternRegexp  *regexp.Regexp   `json:"-"`
	exceptionsRegexps []*regexp.Regexp `json:"-"`
	rulesRegexps      []*regexp.Regexp `json:"-"`
	rawRulesRegexps   []*regexp.Regexp `json:"-"`
}

func (p *Provider) initExceptions() error {
	var err error
	if len(p.Exceptions) > 0 && len(p.exceptionsRegexps) == 0 {
		exceptionsRegexps := make([]*regexp.Regexp, len(p.Exceptions))
		for i, e := range p.Exceptions {
			exceptionsRegexps[i], err = regexp.Compile(e)
			if err != nil {
				return err
			}
		}
		p.exceptionsRegexps = exceptionsRegexps
	}
	return nil
}

func (p *Provider) initRules() error {
	var err error
	if len(p.Rules) > 0 && len(p.rulesRegexps) == 0 {
		rulesRegexps := make([]*regexp.Regexp, len(p.Rules))
		for i, r := range p.Rules {
			rulesRegexps[i], err = regexp.Compile(r)
			if err != nil {
				return err
			}
		}
		p.rulesRegexps = rulesRegexps
	}
	return nil
}

func (p *Provider) initRawRules() error {
	var err error
	if len(p.RawRules) > 0 && len(p.rawRulesRegexps) == 0 {
		rawRulesRegexps := make([]*regexp.Regexp, len(p.RawRules))
		for i, r := range p.RawRules {
			rawRulesRegexps[i], err = regexp.Compile(r)
			if err != nil {
				return err
			}
		}
		p.rawRulesRegexps = rawRulesRegexps
	}
	return nil
}

func (p *Provider) initUrlPattern() error {
	var err error
	if p.urlPatternRegexp == nil {
		p.urlPatternRegexp, err = regexp.Compile(p.UrlPattern)
		if err != nil {
			return err
		}
	}

	return nil
}

type Handler struct {
	providers map[string]*Provider
}

func Init() (*Handler, error) {
	h := &Handler{
		providers: rulesMap["providers"],
	}

	for name, p := range h.providers {
		p.Name = name
		if err := p.initUrlPattern(); err != nil {
			return nil, err
		}
	}

	return h, nil
}

func (h *Handler) Preload() error {
	for _, p := range h.providers {
		if err := p.initExceptions(); err != nil {
			return err
		}
		if err := p.initRawRules(); err != nil {
			return err
		}
		if err := p.initRawRules(); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) Clear(rawURL string) (string, error) {
	cleanedURL := rawURL
	u, err := url.Parse(cleanedURL)
	if err != nil {
		return "", err
	}
	values := u.Query()

	for _, p := range h.providers {
		if !p.urlPatternRegexp.MatchString(rawURL) || p.CompleteProvider {
			continue
		}

		if err := p.initExceptions(); err != nil {
			return "", err
		}

		isException := false
		for _, e := range p.exceptionsRegexps {
			if e.MatchString(rawURL) {
				isException = true
				break
			}
		}
		if isException {
			continue
		}

		if err := p.initRawRules(); err != nil {
			return "", err
		}

		for _, r := range p.rawRulesRegexps {
			cleanedURL = r.ReplaceAllString(cleanedURL, "")
		}

		if err := p.initRules(); err != nil {
			return "", err
		}

		for k := range values {
			remove := false
			for _, rule := range p.rulesRegexps {
				if rule.MatchString(strings.ToLower(k)) {
					remove = true
					break
				}
			}
			if remove {
				values.Del(k)
			}
		}
	}

	u, _ = url.Parse(cleanedURL)
	u.RawQuery = values.Encode()
	return u.String(), nil
}
