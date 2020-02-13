package apisix

import "strings"

// ip-restrictio
type IpRestriction struct {
	Whitelist []string `json:"whitelist,omitempty"`
	Blacklist []string `json:"blacklist,omitempty"`
}

// Convert2IpRestriction build IpRestriction
func BuildIpRestriction(whites, blacks *string) *IpRestriction{
	result := &IpRestriction{}
	if whites != nil {
		whiteIps := strings.Split(*whites, ",")
		result.Whitelist = whiteIps
	}
	if blacks != nil {
		blackIps := strings.Split(*blacks, ",")
		result.Blacklist = blackIps
	}
	return result
}