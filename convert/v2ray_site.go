package convert

import (
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/metacubex/geo/encoding/clashrule"
	"github.com/metacubex/geo/encoding/singgeo"
	"github.com/metacubex/geo/encoding/v2raygeo"

	"github.com/sagernet/sing/common"
	"gopkg.in/yaml.v3"
)

var OutDir string

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

func V2RayToYamlByCode(geositeList []*v2raygeo.GeoSite, output io.Writer, targetCode string) error {
	domainMap := make(map[string][]string, len(geositeList))
	for _, vGeositeEntry := range geositeList {
		code := strings.ToLower(vGeositeEntry.CountryCode)
		if strings.EqualFold(code, targetCode) {
			domains := make([]string, 0, len(vGeositeEntry.Domain)*2)
			attributes := make(map[string][]*v2raygeo.Domain)
			for _, domain := range vGeositeEntry.Domain {
				if len(domain.Attribute) > 0 {
					for _, attribute := range domain.Attribute {
						attributes[attribute.Key] = append(attributes[attribute.Key], domain)
					}
				}
				ruleType := clashrule.GetClashRule(domain.Type)
				switch domain.Type {
				case v2raygeo.Domain_Plain:
					domains = append(domains, ruleType+domain.Value)
				case v2raygeo.Domain_Domain:
					domains = append(domains, ruleType+domain.Value)
				case v2raygeo.Domain_Full:
					domains = append(domains, ruleType+domain.Value)
				case v2raygeo.Domain_Regex:
					continue
				}
			}
			sort.Strings(domains)
			domainMap[targetCode] = common.Uniq(domains)
		}

	}

	yamlOutput := map[string]interface{}{
		"payload": domainMap[targetCode],
	}
	yamlBytes, err := yaml.Marshal(yamlOutput)
	if err != nil {
		return err
	}
	_, err = output.Write(yamlBytes)
	return err

}

func processGeositeEntry(vGeositeEntry *v2raygeo.GeoSite, targetCode string) {
	code := strings.ToLower(vGeositeEntry.CountryCode)
	var domains []string
	domainMap := map[string][]string{}

	entry := strings.Builder{}

	if targetCode == "" || strings.EqualFold(code, targetCode) {
		for _, domain := range vGeositeEntry.Domain {
			entry.Reset()
			entry.WriteString(domain.Type.String())
			entry.WriteString(":")
			entry.WriteString(domain.Value)

			for _, attribute := range domain.Attribute {
				entry.WriteString(" @" + attribute.Key)
			}

			domains = append(domains, entry.String())
		}
		sort.Strings(domains)
		domainMap[code] = common.Uniq(domains)
		WriteByCode(code, domains)
	}
}

func UnpackByCode(geositeList []*v2raygeo.GeoSite, targetCode string) error {
	for _, vGeositeEntry := range geositeList {
		processGeositeEntry(vGeositeEntry, targetCode)
	}

	return nil
}

func WriteByCode(code string, domains []string) error {
	filename := path.Clean(path.Join(OutDir, code))

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(strings.Join(domains, "\n") + "\n")
	return err
}
