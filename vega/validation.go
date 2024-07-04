package vega

import "regexp"

func IsRetentionPolicyValid(policy string) bool {
	if policy == "standard" || policy == "forever" || policy == "1 day" {
		return true
	}

	aDayMatch, _ := regexp.Match(`^1 (day|month|year)$`, []byte(policy))
	multipleDayMatch, _ := regexp.Match(`^\d+ (days|months|years)$`, []byte(policy))

	return aDayMatch || multipleDayMatch
}
