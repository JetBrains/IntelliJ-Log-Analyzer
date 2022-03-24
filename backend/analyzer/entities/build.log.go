package entities

import (
	"log_analyzer/backend/analyzer"
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
	if strings.Contains(path, "build.log") {
		return true
	}
	return false
}
