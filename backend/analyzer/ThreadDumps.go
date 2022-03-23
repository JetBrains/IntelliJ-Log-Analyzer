package analyzer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
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
			fileInfo, _ := os.Stat(path)
			if !fileInfo.IsDir() {
				content, _ := ioutil.ReadFile(path)
				t[path] = ThreadDumpFile{
					Content:     string(content),
					DateAndTime: GetTimeStampFromThreadDump(path),
				}
			}
		}
		return nil
	}
	_ = filepath.WalkDir(path, visit)
	return t
}
func (t *ThreadDumpFile) ConvertToHTML() (html string) {
	return t.Content
}

func (t *ThreadDump) GetFile(path string) *ThreadDumpFile {
	for threadDumpElementPath, file := range *t {
		if strings.Contains(threadDumpElementPath, path) {
			return &file
		}
	}
	return nil
}
func (t *ThreadDump) GetElementWithNumber(number int) *ThreadDumpFile {
	threadDump := *t
	for i, path := range sortedKeys(threadDump) {
		if i == number {
			file := threadDump[path]
			return &file
		}
	}
	return nil
}

func (t *ThreadDump) GetFiltersHTML() (html string) {
	for _, i := range sortedKeys(*t) {
		html = html + "<li filename='" + filepath.Base(i) + "'>" + GetThreadDumpDisplayName(filepath.Base(i)) + "</li>"
	}
	return html
}

func sortedKeys[K string, V any](m map[K]V) []K {
	keys := make([]K, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func getTimeStampFromThreadDumpFilename(filename string) time.Time {
	filename = path.Base(filename)
	timeStamp := GetRegexNamedCapturedGroups(`(?P<Year>\d{4})(?P<Month>\d{2})(?P<Day>\d{2})-(?P<Hours>\d{2})(?P<Minutes>\d{2})(?P<Seconds>\d{2})`, filename)
	if len(timeStamp) == 0 {
		return time.Time{}
	}
	s := fmt.Sprintf("%s-%s-%sT%s:%s:%sZ", timeStamp["Year"], timeStamp["Month"], timeStamp["Day"], timeStamp["Hours"], timeStamp["Minutes"], timeStamp["Seconds"])
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		log.Println(err)
	}
	return t
}

//GetTimeStampFromThreadDump returns zero time.Time{} for file/dir that does not contain timestamp
//							 returns time.Time{} for file/dir that contains timestamp
func GetTimeStampFromThreadDump(path string) (t time.Time) {
	fileinfo, _ := os.Stat(path)
	if getTimeStampFromThreadDumpFilename(filepath.Base(path)).IsZero() {
		if fileinfo == nil {
			return t
		} else if fileinfo.IsDir() {
			return getTimeStampFromFirstInnerFile(path)
		}
		return t
	}
	return getTimeStampFromThreadDumpFilename(filepath.Base(path))
}

//GetThreadDumpDisplayName retuns formatted timestamp string if it is possible to convert the filename to timestamp
// If it is not possible - returns filename
func GetThreadDumpDisplayName(path string) string {
	t := time.Time{}
	if GetTimeStampFromThreadDump(path).IsZero() {
		return filepath.Base(path)
	} else {
		t = GetTimeStampFromThreadDump(path)
	}
	duration := GetRegexNamedCapturedGroups(`(?P<Duration>\d{2})sec$`, path)["Duration"]
	if len(duration) > 0 {
		duration = fmt.Sprintf("(%ss)", duration)
	}
	s := fmt.Sprintf("%s %s", t.Format("2006.01.02 15:04:05"), duration)
	return s
}

//scans all the filenames in ThreadDumps directory and returns the timestamp of the earliest
func getTimeStampFromFirstInnerFile(path string) (threadDumpTime time.Time) {
	threadDumpTime = time.Now()
	var wg sync.WaitGroup
	visit := func(path string, file os.DirEntry, err error) error {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if t := getTimeStampFromThreadDumpFilename(path); !t.IsZero() && t.Before(threadDumpTime) {
				threadDumpTime = t
			}
		}()
		return nil
	}
	_ = filepath.WalkDir(path, visit)
	wg.Wait()
	return threadDumpTime
}
