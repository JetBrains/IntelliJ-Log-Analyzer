package analyzer

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

//List of files mapped to info from one file of ThreadDump
type ThreadDump map[string]ThreadDumpFile

type ThreadDumpFile struct {
	Content     string
	DateAndTime time.Time
}

//List of ThreadDumps Folders
type AggregatedThreadDumps map[string]ThreadDump

func analyzeThreadDumpsFolder(path string, threadDumpsFolder string) ThreadDump {
	t := make(ThreadDump)
	visit := func(path string, file os.DirEntry, err error) error {
		if strings.Contains(path, threadDumpsFolder) {
			t[path] = ThreadDumpFile{
				Content:     "content",
				DateAndTime: time.Time{},
			}
		}
		return nil
	}
	_ = filepath.WalkDir(path, visit)
	return t
}
func (t *ThreadDump) ConvertToHTML() (html string) {
	for path, file := range *t {
		html = html + "\n" + path + " " + file.Content
	}
	return html
}
