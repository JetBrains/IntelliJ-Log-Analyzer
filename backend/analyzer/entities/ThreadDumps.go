package entities

import (
	"log"
	"log_analyzer/backend/analyzer"
	"os"
	"path/filepath"
	"regexp"
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
	return filepath.Base(path)
}

//getLogEntry represents ThreadDump folder as a Log entry.
func getLogEntry(path string) analyzer.Logs {
	logToPass := []analyzer.LogEntry{}
	logToPass = append(logToPass, analyzer.LogEntry{
		Severity: "FREEZE",
		Time:     getTimeStampFromThreadDump(path).Add(-5 * time.Second),
		Text:     "Freeze started: " + getThreadDumpDisplayName(path) + "\n",
	})
	return logToPass
}
func getTimeStampFromThreadDump(str string) time.Time {
	dateMatcher := regexp.MustCompile("(\\d{8}-\\d{6})")
	if !dateMatcher.MatchString(str) {
		log.Println("Error parsing time from Thread Dump path: " + str)
		return time.Time{}
	}
	str = dateMatcher.FindString(str)
	str = str[:4] + "-" + str[4:6] + "-" + str[6:8] + "T" + str[9:11] + ":" + str[11:13] + ":" + str[13:] + "Z"
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		log.Println(err)
	}
	return t
}
