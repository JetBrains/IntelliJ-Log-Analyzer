package analyzer

import (
	"github.com/nxadm/tail"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"log"
	"os"
)

func (a *Analyzer) EnableLogsLiveUpdate() {
	if len(a.fileWatchers) > 0 {
		log.Println("Logs live update already enabled")
		for _, watcher := range a.fileWatchers {
			log.Printf("file watcher: %v", watcher)
		}
		return
	}
	for entityIndex, entity := range a.DynamicEntities {
		for path, instanceProperties := range entity.entityInstances {
			if instanceProperties.Visible {
				if entity.GetChangeablePath != nil {
					if changeablePath := entity.GetChangeablePath(path); changeablePath != "" {
						if s, err := os.Stat(changeablePath); !s.IsDir() && err == nil && entity.ConvertStringToLogs != nil {
							go a.addWatcher(changeablePath, entityIndex)
						}
					}
				}
			}
		}
	}
}

func (a *Analyzer) addWatcher(logFile string, entityIndex int) {
	for _, watcher := range a.fileWatchers {
		if watcher.Filename == logFile {
			return
		}
	}
	t, err := tail.TailFile(logFile, tail.Config{
		Follow:   true,
		Poll:     true,
		Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}, // <- line changed
	})
	if err != nil {
		log.Println(err)
	}
	a.fileWatchers = append(a.fileWatchers, t)
	log.Printf("Enabled File watcher for: %v", logFile)

	previousLogEntry := ""
	// t.Lines channel sends new lines of a file to "for" loop.
	// Inside the loop, lines are being checked by ConvertStringToLogs function.
	// If a line could be converted to log entry (ConvertStringToLogs return nil error), it is being remembered as previousLogEntry.
	// If [the next line is also a log entry],  previousLogEntry is being appended to the logs entry list. If next line is the last line of file, it also is being appended to the logs entry list.
	// If [the next line is not a log entry], it is being appended to the previous log line (as it is part of the same log entry).
	for line := range t.Lines {
		f, _ := os.Open(t.Filename)
		seek, _ := f.Seek(0, 2)
		lineIsLast := line.SeekInfo.Offset == seek
		if _, err := a.DynamicEntities[entityIndex].ConvertStringToLogs(line.Text); err == nil {
			if len(previousLogEntry) != 0 {
				a.attachToLogsStruct(previousLogEntry, entityIndex, t.Filename)
			}
			if lineIsLast {
				a.attachToLogsStruct(line.Text, entityIndex, t.Filename)
				previousLogEntry = ""
			} else {
				previousLogEntry = line.Text
			}
		} else {
			previousLogEntry = previousLogEntry + "\n" + line.Text
			if lineIsLast {
				a.attachToLogsStruct(previousLogEntry, entityIndex, t.Filename)
				previousLogEntry = ""
			}
		}

	}
}
func (a *Analyzer) attachToLogsStruct(s string, i int, path string) {
	name := a.DynamicEntities[i].Name
	properties := a.DynamicEntities[i].entityInstances[path]
	l, e := a.DynamicEntities[i].ConvertStringToLogs(s)
	if e == nil {
		a.AggregatedLogs.Append(name, properties, l)
		wailsruntime.EventsEmit(*a.Context, "LogsUpdated", l.ConvertToHTML())
	} else {
		log.Printf("Cannot convert string to logs: %s \n string: %s", e, s)
	}
}
