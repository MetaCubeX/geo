package geosite

import (
	"regexp"

	"github.com/metacubex/geo/encoding/singgeo"

	"github.com/sagernet/sing/common/domain"
)

func loadSing(reader *singgeo.GeoSiteReader, codes []string) (db *Database, err error) {
	db = &Database{
		matchers:      make(map[string]*domain.Matcher, len(codes)),
		domainKeyword: make(map[string][]string, len(codes)),
		domainRegex:   make(map[string][]*regexp.Regexp, len(codes)),
		SourceType:    TypeSing,
		CodeCount:     len(codes),
	}

	var rules []singgeo.Item
	for _, code := range codes {
		rules, err = reader.Read(code)
		if err != nil {
			return nil, err
		}
		var (
			domainFull    []string
			domainSuffix  []string
			domainKeyword []string
			domainRegex   []*regexp.Regexp
		)
		for _, rule := range rules {
			switch rule.Type {
			case singgeo.RuleTypeDomain:
				domainFull = append(domainFull, rule.Value)
			case singgeo.RuleTypeDomainSuffix:
				domainSuffix = append(domainSuffix, rule.Value)
			case singgeo.RuleTypeDomainKeyword:
				domainKeyword = append(domainKeyword, rule.Value)
			case singgeo.RuleTypeDomainRegex:
				var re *regexp.Regexp
				re, err = regexp.Compile(rule.Value)
				if err != nil {
					return nil, err
				}
				domainRegex = append(domainRegex, re)
			}
		}

		if len(domainFull)+len(domainSuffix) > 0 {
			db.matchers[code] = domain.NewMatcher(domainFull, domainSuffix)
		}
		db.domainKeyword[code] = domainKeyword
		db.domainRegex[code] = domainRegex
	}

	return db, nil
}
