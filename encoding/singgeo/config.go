package singgeo

type ItemType = uint8

const (
	RuleTypeDomain ItemType = iota
	RuleTypeDomainSuffix
	RuleTypeDomainKeyword
	RuleTypeDomainRegex
)

type Item struct {
	Type  ItemType
	Value string
}
