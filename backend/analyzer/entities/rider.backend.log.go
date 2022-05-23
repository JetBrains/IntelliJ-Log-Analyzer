package entities

import (
	"bufio"
	"fmt"
	"log"
	"log_analyzer/backend/analyzer"
	"os"
	"regexp"
	"strings"
	"time"
)

func init() {
	CurrentAnalyzer.AddDynamicEntity(analyzer.DynamicEntity{
		Name:                  "Rider Backend Log",
		ConvertPathToLogs:     parseRiderBackendLog,
		CheckPath:             isRiderBackendLog,
		CheckIgnoredPath:      isIgnoredRiderBackendFile,
		DefaultVisibility:     isBackendLogVisible,
		GetDisplayName:        getDisplayName,
		LineHighlightingColor: "#58c0f1",
	})
	CurrentAnalyzer.AddDynamicEntity(analyzer.DynamicEntity{
		Name:                  "Rider DebuggerWorker",
		ConvertPathToLogs:     parseRiderBackendLog,
		CheckPath:             isRiderDebuggerWorkerLog,
		CheckIgnoredPath:      isIgnoredRiderBackendFile,
		DefaultVisibility:     isPathVisible,
		GetDisplayName:        getDisplayName,
		LineHighlightingColor: "#58c0f1",
	})
	CurrentAnalyzer.AddDynamicEntity(analyzer.DynamicEntity{
		Name:                  "Rider RoslynWorker",
		ConvertPathToLogs:     parseRiderBackendLog,
		CheckPath:             isRiderRoslynWorkerLog,
		CheckIgnoredPath:      isIgnoredRiderBackendFile,
		DefaultVisibility:     isPathVisible,
		GetDisplayName:        getDisplayName,
		LineHighlightingColor: "#58c0f1",
	})
	CurrentAnalyzer.AddDynamicEntity(analyzer.DynamicEntity{
		Name:                  "Rider SolutionBuilder",
		ConvertPathToLogs:     parseRiderBackendLog,
		CheckPath:             isRiderSolutionBuilderLog,
		CheckIgnoredPath:      isIgnoredRiderBackendFile,
		DefaultVisibility:     isPathVisible,
		GetDisplayName:        getDisplayName,
		LineHighlightingColor: "#58c0f1",
	})
	CurrentAnalyzer.AddDynamicEntity(analyzer.DynamicEntity{
		Name:                  "Rider UnitTestLogs",
		ConvertPathToLogs:     parseRiderBackendLog,
		CheckPath:             isRiderUnitTestLog,
		CheckIgnoredPath:      isIgnoredRiderBackendFile,
		DefaultVisibility:     isPathVisible,
		GetDisplayName:        getDisplayName,
		LineHighlightingColor: "#58c0f1",
	})
	CurrentAnalyzer.AddDynamicEntity(analyzer.DynamicEntity{
		Name:                  "Rider MsBuildTask",
		ConvertPathToLogs:     parseRiderBackendLog,
		CheckPath:             isRiderMsBuildTaskLog,
		CheckIgnoredPath:      isIgnoredRiderBackendFile,
		DefaultVisibility:     isPathVisible,
		GetDisplayName:        getDisplayName,
		LineHighlightingColor: "#58c0f1",
	})
}

// ignore files bigger than 49MB
func isBackendLogVisible(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return true
	}
	if f.Size() >= 47*1024*1024 {
		log.Printf("File %s is too big (%d MB) to be displayed", path, f.Size()/(1024*1024))
		return false
	}
	return true
}
func isPathVisible(path string) bool {
	return isRiderBackendLog(path)
}
func isRiderMsBuildTaskLog(path string) bool {
	return strings.Contains(path, "MsBuildTask")
}
func isRiderUnitTestLog(path string) bool {
	return strings.Contains(path, "UnitTestLogs")
}
func isRiderSolutionBuilderLog(path string) bool {
	return strings.Contains(path, "SolutionBuilder") && strings.Contains(path, "ReSharperBuild")
}
func isRiderRoslynWorkerLog(path string) bool {
	logMatcher := regexp.MustCompile(`backend.\d+.log`)
	return (strings.Contains(path, "backend.log") || logMatcher.MatchString(path)) && strings.Contains(path, "RoslynWorker")
}

func isRiderDebuggerWorkerLog(path string) bool {
	logMatcher := regexp.MustCompile(`backend.\d+.log`)
	return (strings.Contains(path, "backend.log") || logMatcher.MatchString(path)) && strings.Contains(path, "DebuggerWorker")
}
func isRiderBackendLog(path string) bool {
	logMatcher := regexp.MustCompile(`backend.\d+.log`)
	return (strings.Contains(path, "backend.log") || logMatcher.MatchString(path)) && !isRiderRoslynWorkerLog(path) && !isRiderDebuggerWorkerLog(path) && !isRiderSolutionBuilderLog(path)
}
func isIgnoredRiderBackendFile(path string) bool {
	return strings.Contains(path, "backend-protocol.log") || strings.Contains(path, "backend-out.log")
}

func parseRiderBackendLog(path string) analyzer.Logs {
	startDate := analyzer.GetFileModTime(path)
	if startDate.IsZero() {
		log.Printf("Could not get creation date for %s", path)
	}
	reader, _ := os.Open(path)
	scanner := bufio.NewScanner(reader)
	logToPass := []analyzer.LogEntry{}
	for scanner.Scan() {
		currentString := scanner.Text()
		if getTimeStringFromRiderBackendLog(currentString) != "" {
			logToPass = append(logToPass, parseRiderBackendLogString(startDate, currentString))
		} else if len(logToPass) > 0 {
			logToPass[len(logToPass)-1].Text = logToPass[len(logToPass)-1].Text + "\n" + currentString
		}
	}

	return logToPass
}

func getTimeStringFromRiderBackendLog(str string) string {
	dateMatcher := regexp.MustCompile(`^(\d{2}:\d{2}:\d{2}[.,]\d{3})`)
	if !dateMatcher.MatchString(str) {
		return ""
	}
	return dateMatcher.FindString(str)
}
func parseRiderBackendLogString(startDate time.Time, logEntryAsString string) (currentEntry analyzer.LogEntry) {
	logEntryAsStringToIdeaLogFormat := fmt.Sprintf("%s %s", startDate.Format("2006-01-02"), logEntryAsString)
	logParts := analyzer.GetRegexNamedCapturedGroups(`(?P<Year>\d{4})-(?P<Month>\d{2})-(?P<Day>\d{2})\s+(?P<Hours>\d{2}):(?P<Minutes>\d{2}):(?P<Seconds>\d{2})[.,](?P<MiliSeconds>\d{3})\s*\|(?P<Severity>\w)\|(?P<Class>.*?)\|(?P<Text>.*)`, logEntryAsStringToIdeaLogFormat)
	if startDate.IsZero() || logParts["Hours"] == "" {
		log.Printf("PARSE_ERROR!\n logEntryAsStringToIdeaLogFormat:%s\n Start Date: %s\n logParts[\"Hours\"]: %s\n", logEntryAsStringToIdeaLogFormat, startDate, logParts["Hours"])
		return analyzer.LogEntry{
			Severity: "PARSE_ERROR",
			Time:     time.Now().AddDate(0, 0, 1),
			Text:     logEntryAsString,
			Visible:  true,
		}
	}
	currentEntry.Time, _ = time.Parse(time.RFC3339Nano, fmt.Sprintf("%s-%s-%sT%s:%s:%s.%sZ", startDate.Format("2006"), startDate.Format("01"), startDate.Format("02"), logParts["Hours"], logParts["Minutes"], logParts["Seconds"], logParts["MiliSeconds"]))
	currentEntry.Severity = getSeverityFromRiderBackendLog(logParts["Severity"])
	currentEntry.Text = logParts["Class"] + " —" + logParts["Text"]
	currentEntry.Visible = true
	return currentEntry
}

func getSeverityFromRiderBackendLog(s string) string {
	if strings.HasPrefix(s, "E") {
		return "ERROR"
	}
	if strings.HasPrefix(s, "W") {
		return "WARN"
	}
	if strings.HasPrefix(s, "V") {
		return "VERB"
	}
	if strings.HasPrefix(s, "T") {
		return "TRACE"
	}
	if strings.HasPrefix(s, "I") {
		return "INFO"
	}
	return "RIDER"
}
