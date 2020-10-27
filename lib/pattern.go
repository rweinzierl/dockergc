package lib

import (
	"regexp"
	"strings"
)

func stringMatchesPattern(s string, pattern string) bool {
	patternParts := strings.Split(pattern, "*")
	for i, part := range patternParts {
		patternParts[i] = regexp.QuoteMeta(part)
	}
	regexPattern := strings.Join(patternParts, ".*")
	matches, _ := regexp.MatchString(regexPattern, s)
	return matches
}

func replaceEmptyPatternListByWildcard(patterns []string) []string {
	if len(patterns) == 0 {
		return []string{"*"}
	}
	return patterns
}
