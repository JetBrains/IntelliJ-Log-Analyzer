package entities

import (
	"fmt"
	"log"
	"log_analyzer/backend/analyzer"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	CurrentAnalyzer.AddDynamicEntity(analyzer.DynamicEntity{
		Name:           "Thread Dumps",
		ConvertToLogs:  getLogEntry,
		CheckPath:      isThreadDump,
		GetDisplayName: getThreadDumpDisplayName,
	})
}
func isThreadDump(path string) bool {
	if strings.Contains(path, "threadDump") {
		fileInfo, _ := os.Stat(path)
		if fileInfo.IsDir() {
			return true
		}
	}
	return false
}

func getThreadDumpDisplayName(path string) string {
	t := getTimeStampFromThreadDump(path)
	duration := getRegexNamedCapturedGroups(`(?P<Duration>\d{2})sec$`, path)["Duration"]
	if len(duration) > 0 {
		duration = fmt.Sprintf("(%ss)", duration)
	}

	s := fmt.Sprintf("%s %s", t.Format("2006.01.02 15:04:05"), duration)
	return s
}

//getLogEntry represents ThreadDump folder as a Log entry.
func getLogEntry(path string) analyzer.Logs {
	logToPass := []analyzer.LogEntry{}
	fileName := filepath.Base(path)
	logToPass = append(logToPass, analyzer.LogEntry{
		Severity: "FREEZE",
		Time:     getTimeStampFromThreadDump(path).Add(-5 * time.Second),
		Text:     "Freeze started: " + fileName + "\n",
	})
	return logToPass
}
func getTimeStampFromThreadDump(filename string) time.Time {
	timeStamp := getRegexNamedCapturedGroups(`(?P<Year>\d{4})(?P<Month>\d{2})(?P<Day>\d{2})-(?P<Hours>\d{2})(?P<Minutes>\d{2})(?P<Seconds>\d{2})`, filename)
	if len(timeStamp) == 0 {
		log.Println("Error parsing time from Thread Dump path: " + filename)
		return time.Time{}
	}
	s := fmt.Sprintf("%s-%s-%sT%s:%s:%sZ", timeStamp["Year"], timeStamp["Month"], timeStamp["Day"], timeStamp["Hours"], timeStamp["Minutes"], timeStamp["Seconds"])
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		log.Println(err)
	}
	return t
}
