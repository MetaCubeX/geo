package convert

import (
	"io"
	"strings"

	"github.com/metacubex/geo/encoding/singgeo"
	"github.com/metacubex/geo/encoding/v2raygeo"

	"github.com/sagernet/sing/common"
)

// V2RaySiteToSing is modified from https://github.com/SagerNet/sing-geosite
func V2RaySiteToSing(geositeList []*v2raygeo.GeoSite, output io.Writer) error {
	domainMap := make(map[string][]singgeo.Item, len(geositeList))
	for _, vGeositeEntry := range geositeList {
		code := strings.ToLower(vGeositeEntry.CountryCode)
		domains := make([]singgeo.Item, 0, len(vGeositeEntry.Domain)*2)
		attributes := make(map[string][]*v2raygeo.Domain)
		for _, domain := range vGeositeEntry.Domain {
			if len(domain.Attribute) > 0 {
				for _, attribute := range domain.Attribute {
					attributes[attribute.Key] = append(attributes[attribute.Key], domain)
				}
			}
			switch domain.Type {
			case v2raygeo.Domain_Plain:
				domains = append(domains, singgeo.Item{
					Type:  singgeo.RuleTypeDomainKeyword,
					Value: domain.Value,
				})
			case v2raygeo.Domain_Regex:
				domains = append(domains, singgeo.Item{
					Type:  singgeo.RuleTypeDomainRegex,
					Value: domain.Value,
				})
			case v2raygeo.Domain_Domain:
				if strings.Contains(domain.Value, ".") {
					domains = append(domains, singgeo.Item{
						Type:  singgeo.RuleTypeDomain,
						Value: domain.Value,
					})
				}
				domains = append(domains, singgeo.Item{
					Type:  singgeo.RuleTypeDomainSuffix,
					Value: "." + domain.Value,
				})
			case v2raygeo.Domain_Full:
				domains = append(domains, singgeo.Item{
					Type:  singgeo.RuleTypeDomain,
					Value: domain.Value,
				})
			}
		}
		domainMap[code] = common.Uniq(domains)
		for attribute, attributeEntries := range attributes {
			attributeDomains := make([]singgeo.Item, 0, len(attributeEntries)*2)
			for _, domain := range attributeEntries {
				switch domain.Type {
				case v2raygeo.Domain_Plain:
					attributeDomains = append(attributeDomains, singgeo.Item{
						Type:  singgeo.RuleTypeDomainKeyword,
						Value: domain.Value,
					})
				case v2raygeo.Domain_Regex:
					attributeDomains = append(attributeDomains, singgeo.Item{
						Type:  singgeo.RuleTypeDomainRegex,
						Value: domain.Value,
					})
				case v2raygeo.Domain_Domain:
					if strings.Contains(domain.Value, ".") {
						attributeDomains = append(attributeDomains, singgeo.Item{
							Type:  singgeo.RuleTypeDomain,
							Value: domain.Value,
						})
					}
					attributeDomains = append(attributeDomains, singgeo.Item{
						Type:  singgeo.RuleTypeDomainSuffix,
						Value: "." + domain.Value,
					})
				case v2raygeo.Domain_Full:
					attributeDomains = append(attributeDomains, singgeo.Item{
						Type:  singgeo.RuleTypeDomain,
						Value: domain.Value,
					})
				}
			}
			domainMap[code+"@"+attribute] = common.Uniq(attributeDomains)
		}
	}
	return singgeo.Write(output, domainMap)
}
