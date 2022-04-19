package entities

import (
	"log_analyzer/backend/analyzer"
	"regexp"
	"strings"
)

func init() {
	CurrentAnalyzer.AddDynamicEntity(analyzer.DynamicEntity{
		Name:                  "Build Log",
		ConvertToLogs:         parseIdeaLog,
		CheckPath:             isBuildLog,
		GetDisplayName:        getDisplayName,
		LineHighlightingColor: "#72cf99",
	})
}

func isBuildLog(path string) bool {
	return strings.Contains(path, "build.log") ||
		regexp.MustCompile(`build.\d+.log`).MatchString(path)
}
