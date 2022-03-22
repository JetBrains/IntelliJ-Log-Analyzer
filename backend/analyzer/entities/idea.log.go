package entities

import (
	"bufio"
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
		if !getTimeStampFromIdeaLog(currentString).IsZero() {
			currentEntry := parseLogString(currentString)
			logToPass = append(logToPass, analyzer.LogEntry{
				Severity: currentEntry.Severity,
				Time:     currentEntry.DateAndTime,
				Text:     currentEntry.Class + " " + currentEntry.Header + " " + currentEntry.Body,
			})
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

func getTimeStampFromIdeaLog(str string) (logTime time.Time) {
	str = getTimeStringFromIdeaLog(str)
	str = strings.Replace(str, " ", "T", 1)
	str = strings.Replace(str, ",", ".", 1)
	str = str + "Z"
	logTime, _ = time.Parse(time.RFC3339Nano, str)
	return logTime
}
func getTimeStringFromIdeaLog(str string) string {
	dateMatcher := regexp.MustCompile("^(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}[.,]\\d{3})")
	if !dateMatcher.MatchString(str) {
		return ""
	}
	return dateMatcher.FindString(str)
}

func parseLogString(logEntryAsString string) (currentEntry LogEntry) {
	logEntryAsString = strings.TrimLeft(logEntryAsString, "\n")
	classEndPosition := 0
	HeaderEndPosition := 0
	currentEntry.DateAndTime = getTimeStampFromIdeaLog(logEntryAsString)
	trimFoundPart(&logEntryAsString, getTimeStringFromIdeaLog(logEntryAsString))
	rawTimeSinceStart := getRawTimeSinceStart(&logEntryAsString)
	currentEntry.TimeSinceStart = getTimeSinceStart(rawTimeSinceStart)
	trimFoundPart(&logEntryAsString, rawTimeSinceStart)
	currentEntry.Severity = getSeverity(&logEntryAsString)
	trimFoundPart(&logEntryAsString, currentEntry.Severity)
	//logEntryAsString = strings.TrimSpace(logEntryAsString)
	//currentEntry.Body = logEntryAsString[0:]
	//todo make that without for loop and unify
	for i := range logEntryAsString {
		if currentEntry.Class == "" && getClass(logEntryAsString[0:i]) != "" {
			classEndPosition = i
			currentEntry.Class = getClass(logEntryAsString[0:i])
		}
		if classEndPosition != 0 {
			HeaderEndPosition = strings.IndexAny(logEntryAsString, "\n")
			if HeaderEndPosition == -1 {
				currentEntry.Header = strings.TrimSpace(logEntryAsString[classEndPosition:])
				break
			} else {
				//todo check if this is needed
				currentEntry.Header = logEntryAsString[classEndPosition:HeaderEndPosition]
			}

		}
		if HeaderEndPosition > 0 {
			currentEntry.Body = logEntryAsString[HeaderEndPosition:]
			break
		}
	}

	return currentEntry
}
func trimFoundPart(stringToCut *string, part string) {
	*stringToCut = strings.TrimSpace(*stringToCut)
	*stringToCut = strings.TrimPrefix(*stringToCut, part)
}
func getRawTimeSinceStart(str *string) string {
	dateMatcher := regexp.MustCompile("\\[(.*?)]")
	if dateMatcher.MatchString(*str) {
		return dateMatcher.FindString(*str)
	}
	return ""
}
func getTimeSinceStart(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 8 && s[0] == '[' && s[len(s)-1] == ']' {
		s = s[1 : len(s)-1]
		return strings.TrimSpace(s)

	}
	return ""
}

func getClass(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 2 && s[0] == '-' && s[len(s)-1] == '-' {
		return strings.TrimSpace(s[0:len(s)])
	}
	return ""
}

func getSeverity(s *string) string {
	*s = strings.TrimSpace(*s)
	severities := []string{"INFO", "ERROR", "DEBUG", "WARN", "SEVERE"}
	for _, severity := range severities {
		if strings.HasPrefix(*s, severity) {
			return severity
		}
	}
	return ""
}
