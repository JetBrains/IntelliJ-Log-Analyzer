package entities

import (
	"fmt"
	"log"
	"log_analyzer/backend/analyzer"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	CurrentAnalyzer.AddDynamicEntity(analyzer.DynamicEntity{
		Name:                  "Indexing diagnostic",
		ConvertPathToLogs:     parseIndexingDiagnosticFolder,
		CheckPath:             isIndexingFile,
		CheckIgnoredPath:      isIgnortedIndexingFile,
		GetDisplayName:        getDisplayName,
		LineHighlightingColor: "#94c794",
	})
}
func isIndexingFolder(path string) bool {
	if strings.Contains(path, "indexing-diagnostic") {
		return true
	}
	return false
}

func parseIndexingDiagnosticFolder(path string) (l analyzer.Logs) {
	if isIndexingFile(path) {
		l = append(l, analyzer.LogEntry{
			Severity: "INDEX",
			Time:     getTimeStampFromIndexingFile(path),
			Text:     "Indexing project: " + getIndexingProjectName(path) + " (show report.html). Report: " + filepath.Base(path),
		})
	}
	return l
}

func isIgnortedIndexingFile(path string) bool {
	if isIndexingFolder(path) && (filepath.Ext(path) == ".json" ||
		strings.Contains(filepath.Base(path), "changed-files-pushing-events.json") ||
		strings.Contains(filepath.Base(path), "report.html") ||
		strings.Contains(filepath.Base(path), "shared-index-events.json")) {
		return true
	}
	return false
}
func isIndexingFile(path string) bool {
	timeStamp := analyzer.GetRegexNamedCapturedGroups(`diagnostic-(?P<Year>\d{4})-(?P<Month>\d{2})-(?P<Day>\d{2})-(?P<Hours>\d{2})-(?P<Minutes>\d{2})-(?P<Seconds>\d{2}).*.html`, path)
	if len(timeStamp) > 0 {
		return true
	}
	return false
}

func getTimeStampFromIndexingFile(path string) time.Time {
	timeStamp := analyzer.GetRegexNamedCapturedGroups(`diagnostic-(?P<Year>\d{4})-(?P<Month>\d{2})-(?P<Day>\d{2})-(?P<Hours>\d{2})-(?P<Minutes>\d{2})-(?P<Seconds>\d{2}).*.html`, path)
	if len(timeStamp) == 0 {
		return time.Time{}
	}
	s := fmt.Sprintf("%s-%s-%sT%s:%s:%sZ", timeStamp["Year"], timeStamp["Month"], timeStamp["Day"], timeStamp["Hours"], timeStamp["Minutes"], timeStamp["Seconds"])
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		log.Printf("getTimeStampFromIndexingFile failed. path: %s, error: %s", path, err)
	}
	return t
}

func getIndexingProjectName(path string) string {
	projectFolder := analyzer.GetRegexNamedCapturedGroups(`indexing-diagnostic.{1}(?P<ProjectFolder>.*).{1}diagnostic-(?P<Year>\d{4})-(?P<Month>\d{2})-(?P<Day>\d{2})-(?P<Hours>\d{2})-(?P<Minutes>\d{2})-(?P<Seconds>\d{2}).*.html`, path)["ProjectFolder"]
	projectName := projectFolder[:strings.LastIndex(projectFolder, ".")]
	return projectName
}
