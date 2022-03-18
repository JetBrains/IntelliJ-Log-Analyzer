package entities

import (
	"log_analyzer/backend/analyzer"
	"regexp"
)

var CurrentAnalyzer = &analyzer.Analyzer{}

/**
 * Parses string s with the given regular expression and returns the
 * group values defined in the expression.
 *
 */
func getRegexNamedCapturedGroups(regEx, s string) (paramsMap map[string]string) {

	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(s)
	paramsMap = make(map[string]string)

	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}
