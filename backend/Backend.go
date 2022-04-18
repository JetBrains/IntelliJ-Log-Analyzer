package backend

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"log_analyzer/backend/analyzer"
	"log_analyzer/backend/analyzer/entities"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//InitLogDirectory creates an instance of analyzed directory (all entities combined) and parses them
func InitLogDirectory(path string) (err error) {
	entities.CurrentAnalyzer.Clear()
	entities.CurrentAnalyzer.FolderToWorkWith = path
	timeStart := time.Now()
	entities.CurrentAnalyzer.ParseLogDirectory(path)
	log.Printf("Parsed Logs in %s", time.Now().Sub(timeStart).String())
	log.Printf("Log Entries Count: %d", len(entities.CurrentAnalyzer.AggregatedLogs))
	entities.CurrentAnalyzer.AggregatedLogs.SortByTime()
	entities.CurrentAnalyzer.GenerateFilters()
	if entities.CurrentAnalyzer.IsEmpty() {
		return errors.New("could not find logs elements inside")
	} else {
		return nil
	}
}

func GetLogs() *analyzer.Logs {
	return entities.CurrentAnalyzer.GetLogs()
}
func GetStaticInfo() *analyzer.AggregatedStaticInfo {
	return entities.CurrentAnalyzer.GetStaticInfo()
}
func GetOtherFiles() *analyzer.OtherFiles {
	return entities.CurrentAnalyzer.GetOtherFiles()
}
func GetFilters() *analyzer.Filters {
	return entities.CurrentAnalyzer.GetFilters()
}
func GetThreadDumpFolder(dir string) *analyzer.ThreadDump {
	return entities.CurrentAnalyzer.GetThreadDump(dir)
}

func GetIndexingFilePath(path string) string {
	for _, s := range entities.CurrentAnalyzer.GetIndexingFilesList() {
		if strings.Contains(s, path) {
			return s
		}
	}
	return ""
}

func UnzipToTempFodler(src string) (dest string) {
	dest, err := ioutil.TempDir("", "IntelliJLogsAnalyzer")
	r, err := zip.OpenReader(src)
	if err != nil {
		return ""
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, 0755)
			//log.Println(err)
		} else {
			err = os.MkdirAll(filepath.Dir(path), 0755)
			if err != nil {
				log.Println(err)
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return ""
		}
	}
	entities.CurrentAnalyzer.IsFolderTemp = true
	log.Println("Temp folder to work in: " + dest)
	return dest
}

// Set the Checked values for all FilterEntry elements from frontend
func SetFilters(f map[string]bool) error {
	for _, entries := range entities.CurrentAnalyzer.Filters {
		for i, entry := range entries {
			for id, value := range f {
				if entry.ID == id {
					entries[i].Checked = value
				}
			}
		}
	}
	return nil
}
