package geosite

import (
	"regexp"
	"strings"

	"github.com/sagernet/sing/common/domain"
)

type Database struct {
	matchers      map[string]*domain.Matcher
	domainKeyword map[string][]string
	domainRegex   map[string][]*regexp.Regexp
	SourceType    Type
	CodeCount     int
}

func (db Database) LookupCode(domainName string) string {
	domainName = strings.ToLower(domainName)
	for code, matcher := range db.matchers {
		if matcher.Match(domainName) {
			return strings.ToLower(code)
		}
	}
	for code, keywords := range db.domainKeyword {
		for _, keyword := range keywords {
			if strings.Contains(domainName, keyword) {
				return strings.ToLower(code)
			}
		}
	}
	for code, regexList := range db.domainRegex {
		for _, re := range regexList {
			if re.MatchString(domainName) {
				return strings.ToLower(code)
			}
		}
	}
	return ""
}

func (db Database) LookupCodes(domainName string) []string {
	results := map[string]struct{}{}
	domainName = strings.ToLower(domainName)
	for code, matcher := range db.matchers {
		if matcher.Match(domainName) {
			results[strings.ToLower(code)] = struct{}{}
		}
	}
	for code, keywords := range db.domainKeyword {
		for _, keyword := range keywords {
			if strings.Contains(domainName, keyword) {
				results[strings.ToLower(code)] = struct{}{}
			}
		}
	}
	for code, regexList := range db.domainRegex {
		for _, re := range regexList {
			if re.MatchString(domainName) {
				results[strings.ToLower(code)] = struct{}{}
			}
		}
	}

	codes := make([]string, 0, len(results))
	for code := range results {
		codes = append(codes, code)
	}
	return codes
}
