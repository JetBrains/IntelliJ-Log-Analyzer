package entities

import (
	"bufio"
	"errors"
	"fmt"
	"log_analyzer/backend/analyzer"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func init() {
	CurrentAnalyzer.AddDynamicEntity(analyzer.DynamicEntity{
		Name:                "Idea Log",
		ConvertPathToLogs:   parseIdeaLogFile,
		CheckPath:           isIdeaLog,
		GetDisplayName:      getDisplayName,
		GetChangeablePath:   getIdeaLogChangeablePath,
		ConvertStringToLogs: parseIdeaLogString,
	})
}
func isIdeaLog(path string) bool {
	logMatcher := regexp.MustCompile(`idea\.\d+.log`)
	if strings.Contains(path, "idea.log") || logMatcher.MatchString(path) {
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
func parseIdeaLogFile(path string) analyzer.Logs {
	reader, _ := os.Open(path)
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	logs := []analyzer.LogEntry{}
	for scanner.Scan() {
		currentString := scanner.Text()
		if entry, err := parseIdeaLogString(currentString); err == nil {
			logs = append(logs, entry)
		} else if len(logs) > 0 {
			logs[len(logs)-1].Text = logs[len(logs)-1].Text + "\n" + entry.Text
		}
	}
	return logs
}

func parseIdeaLogString(logEntryAsString string) (currentEntry analyzer.LogEntry, err error) {
	logParts := analyzer.GetRegexNamedCapturedGroups(`(?s)(?P<Year>\d{4})-(?P<Month>\d{2})-(?P<Day>\d{2})\s+(?P<Hours>\d{2}):(?P<Minutes>\d{2}):(?P<Seconds>\d{2})[.,](?P<MiliSeconds>\d{3})\s+\[\s*(?P<Duration>\d+)\]\s+(?P<Severity>[A-Z]+)\s+\-\s*(?P<Class>.*?)\s*-(?P<Body>.*)`, logEntryAsString)
	if logParts["Year"] == "" {
		return analyzer.LogEntry{
			Text: logEntryAsString,
		}, errors.New("Could not parse idea.log string: " + logEntryAsString)
	}

	currentEntry.Time, _ = time.Parse(time.RFC3339Nano, fmt.Sprintf("%s-%s-%sT%s:%s:%s.%sZ", logParts["Year"], logParts["Month"], logParts["Day"], logParts["Hours"], logParts["Minutes"], logParts["Seconds"], logParts["MiliSeconds"]))
	currentEntry.Severity = strings.TrimSpace(logParts["Severity"])
	currentEntry.Text = logParts["Class"] + " —" + logParts["Body"]

	if logParts["Class"] == "STDERR" {
		currentEntry.Text = strings.TrimPrefix(currentEntry.Text, "STDERR —")
		if !strings.Contains(logParts["Body"], "\t") {
			return analyzer.LogEntry{
				Severity: "EXCPT",
				Time:     currentEntry.Time,
				Text:     currentEntry.Text,
				Visible:  true,
			}, nil
		}
		return analyzer.LogEntry{
			Severity: "PARSE_ERROR",
			Time:     currentEntry.Time,
			Text:     currentEntry.Text,
			Visible:  true,
		}, errors.New("found STDERR in the entry body")
	}

	return currentEntry, err
}

func getIdeaLogChangeablePath(path string) string {
	if strings.HasSuffix(path, "idea.log") {
		return path
	}
	return ""
}
