package entities

import (
	"log_analyzer/backend/analyzer"
	"regexp"
	"strings"
)

func init() {
	CurrentAnalyzer.AddDynamicEntity(analyzer.DynamicEntity{
		Name:                  "Build Log",
		ConvertPathToLogs:     parseIdeaLogFile,
		CheckPath:             isBuildLog,
		GetDisplayName:        getDisplayName,
		LineHighlightingColor: "#72cf99",
		GetChangeablePath:     getBuildLogChangeablePath,
		ConvertStringToLogs:   parseIdeaLogString,
	})
}

func isBuildLog(path string) bool {
	return strings.Contains(path, "build.log") ||
		regexp.MustCompile(`build.\d+.log`).MatchString(path)
}
func getBuildLogChangeablePath(path string) string {
	if strings.HasSuffix(path, "build.log") {
		return path
	}
	return ""
}
