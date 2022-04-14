package entities

import (
	"bufio"
	"fmt"
	"io"
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
		Name:           "Idea Log",
		ConvertToLogs:  parseIdeaLog,
		CheckPath:      isIdeaLog,
		GetDisplayName: getDisplayName,
	})
}
func isIdeaLog(path string) bool {
	if strings.Contains(path, "idea.log") {
		return true
	}
	return false
}

type LogEntries []LogEntry
type LogEntry struct {
	DateAndTime    time.Time
	TimeSinceStart string
	Severity       string
	Class          string
	Header         string
	Body           string
}

func getDisplayName(path string) string {
	return filepath.Base(path)
}
func parseIdeaLog(path string) analyzer.Logs {
	reader, _ := os.Open(path)
	bufReader := bufio.NewReader(reader)
	logToPass := []analyzer.LogEntry{}
	for {
		currentString, err := bufReader.ReadString('\n')
		if getTimeStringFromIdeaLog(currentString) != "" {
			logToPass = append(logToPass, parseIdeaLogString(currentString))
		} else if len(logToPass) > 0 {
			logToPass[len(logToPass)-1].Text = logToPass[len(logToPass)-1].Text + currentString
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
	}

	return logToPass
}

func getTimeStringFromIdeaLog(str string) string {
	dateMatcher := regexp.MustCompile("^(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}[.,]\\d{3})")
	if !dateMatcher.MatchString(str) {
		return ""
	}
	return dateMatcher.FindString(str)
}
func parseIdeaLogString(logEntryAsString string) (currentEntry analyzer.LogEntry) {
	logParts := analyzer.GetRegexNamedCapturedGroups(`(?P<Year>\d{4})-(?P<Month>\d{2})-(?P<Day>\d{2})\s+(?P<Hours>\d{2}):(?P<Minutes>\d{2}):(?P<Seconds>\d{2})[.,](?P<MiliSeconds>\d{3})\s+\[\s*(?P<Duration>\d+)\]\s+(?P<Severity>[A-Z]+)\s+\-(?P<Class>.*?\s)-(?P<Body>.*)\n`, logEntryAsString)
	if logParts["Year"] == "" {
		return analyzer.LogEntry{
			Severity: "PARSE_ERROR",
			Time:     time.Now().AddDate(0, 0, 1),
			Text:     logEntryAsString,
			Visible:  true,
		}
	}
	currentEntry.Time, _ = time.Parse(time.RFC3339Nano, fmt.Sprintf("%s-%s-%sT%s:%s:%s.%sZ", logParts["Year"], logParts["Month"], logParts["Day"], logParts["Hours"], logParts["Minutes"], logParts["Seconds"], logParts["MiliSeconds"]))
	currentEntry.Severity = strings.TrimSpace(logParts["Severity"])
	currentEntry.Text = logParts["Class"] + "—" + logParts["Body"]
	return currentEntry
}
