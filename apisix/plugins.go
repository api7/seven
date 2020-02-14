package apisix

import "strings"

// ip-restrictio
type IpRestriction struct {
	Whitelist []string `json:"whitelist,omitempty"`
	Blacklist []string `json:"blacklist,omitempty"`
}

// Convert2IpRestriction build IpRestriction
func BuildIpRestriction(whites, blacks *string) *IpRestriction {
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

type Cors struct {
	Origins []string `json:"origins,omitempty"`
	Headers []string `json:"headers,omitempty"`
	Methods []string `json:"methods,omitempty"`
	MaxAge  int64    `json:"max_age,omitempty"`
}

// BuildCors
func BuildCors(enable bool, originStr, headerStr, methodStr *string, maxAge *int64) *Cors{
	result := &Cors{}
	if enable {
		if originStr != nil {
			origins := strings.Split(*originStr, ",")
			result.Origins = origins
		}
		if headerStr != nil {
			headers := strings.Split(*headerStr, ",")
			result.Headers = headers
		}
		if methodStr != nil {
			methods := strings.Split(*methodStr, ",")
			result.Methods = methods
		}
		if maxAge != nil {
			result.MaxAge = *maxAge
		}
		return result
	}else {
		return nil
	}
}