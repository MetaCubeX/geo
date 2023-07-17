package geosite

import (
	"regexp"
	"strings"

	"github.com/metacubex/geo/encoding/v2raygeo"

	"github.com/sagernet/sing/common/domain"
)

func loadV2RayFromBytes(data []byte) (*Database, bool) {
	geositeList, err := v2raygeo.LoadSite(data)
	if err != nil {
		return nil, false
	}

	db := &Database{
		matchers:      make(map[string]*domain.Matcher, len(geositeList)),
		domainKeyword: make(map[string][]string, len(geositeList)),
		domainRegex:   make(map[string][]*regexp.Regexp, len(geositeList)),
		SourceType:    TypeV2Ray,
		CodeCount:     len(geositeList),
	}

	for _, code := range geositeList {
		var (
			domainFull    []string
			domainSuffix  []string
			domainKeyword []string
			domainRegex   []*regexp.Regexp
		)
		for _, rule := range code.Domain {
			switch rule.Type {
			case v2raygeo.Domain_Full:
				domainFull = append(domainFull, rule.Value)
			case v2raygeo.Domain_Domain:
				if strings.Contains(rule.Value, ".") {
					domainFull = append(domainFull, rule.Value)
				}
				domainSuffix = append(domainSuffix, "."+rule.Value)
			case v2raygeo.Domain_Plain:
				domainKeyword = append(domainKeyword, rule.Value)
			case v2raygeo.Domain_Regex:
				var re *regexp.Regexp
				re, err = regexp.Compile(rule.Value)
				if err != nil {
					return nil, false
				}
				domainRegex = append(domainRegex, re)
			}
		}

		if len(domainFull)+len(domainSuffix) > 0 {
			db.matchers[code.CountryCode] = domain.NewMatcher(domainFull, domainSuffix)
		}
		db.domainKeyword[code.CountryCode] = domainKeyword
		db.domainRegex[code.CountryCode] = domainRegex
	}

	return db, true
}
