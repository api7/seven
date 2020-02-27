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

// Cors
type Cors struct {
	Origins []string `json:"origins,omitempty"`
	Headers []string `json:"headers,omitempty"`
	Methods []string `json:"methods,omitempty"`
	MaxAge  int64    `json:"max_age,omitempty"`
}

// BuildCors
func BuildCors(enable bool, originStr, headerStr, methodStr *string, maxAge *int64) *Cors {
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
	} else {
		return nil
	}
}

// routex
type Routex struct {
	Rules []Rule `json:"rules,inline"`
}

type Rule struct {
	Priority int64  `json:"priority,omitempty"`
	Upstream string `json:"upstream"`
	Desc     string `json:"desc"`
	Matches []Match `json:"matchs,omitempty"`
}

type Match struct {
	Host   string   `json:"host,omitempty"`
	Uri    string   `json:"uri,omitempty"`
	Use    string   `json:"use"`
	Key    string   `json:"key"`
	Values []string `json:"values,omitempty"`
}

// BuildRoutex
func BuildRoutex(enable bool, rules []Rule) *Routex{
	if enable {
		result := &Routex{Rules: rules}
		return result
	}else {
		return nil
	}
}

// token
type Token struct {
	IgnoreUri []string `json:"ignore_uri,omitempty"`
}

// BuildToken
func BuildToken(enable bool, ignoreUris []string) *Token {
	if enable {
		result := &Token{IgnoreUri: ignoreUris}
		return result
	} else {
		return nil
	}
}
