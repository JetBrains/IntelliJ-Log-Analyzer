package main

import (
	"context"
	"encoding/base64"
	json "encoding/json"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"log_analyzer/backend"
	"log_analyzer/backend/analyzer"
	"os"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (b *App) startup(ctx context.Context) {
	// Perform your setup here
	b.ctx = ctx
}

// domReady is called after the front-end dom has been loaded
func (b *App) domReady(ctx context.Context) {
	// Add your action here
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (b *App) OpenFolder() string {
	path, _ := runtime.OpenDirectoryDialog(b.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: "",
		DefaultFilename:  "",
		Title:            "Open direcotry with logs",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "test Display name",
				Pattern:     "",
			},
		},
		ShowHiddenFiles:            false,
		CanCreateDirectories:       false,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	},
	)

	if err := backend.InitLogDirectory(path); err == nil {
		return path
	} else {
		return ""
	}
}
func (b *App) UploadLogFile(filename string, DataURIScheme string) string {
	data := ConvertDataURISchemeToBase64File(DataURIScheme)
	log.Println(DataURIScheme)
	f, err := os.CreateTemp("", filename)
	if err != nil {
		log.Println("Could not create temp file: " + err.Error())
	}
	_, err = f.Write(data)
	if err != nil {
		log.Println("Could not write to temp file:" + err.Error())
	}
	_ = f.Close()
	log.Println("Created file: " + f.Name())
	if err := backend.InitLogDirectory(f.Name()); err == nil {
		return f.Name()
	} else {
		return ""
	}
}
func (b *App) UploadArchive(DataURIScheme string) string {
	data := ConvertDataURISchemeToBase64File(DataURIScheme)
	f, err := os.CreateTemp("", "logs.zip")
	if err != nil {
		log.Println("Could not create temp file: " + err.Error())
	}
	_, err = f.Write(data)
	if err != nil {
		log.Println("Could not write to temp file:" + err.Error())
	}
	_ = f.Close()
	if err != nil {
		log.Println("Could not get path of the temp file: " + err.Error())
	}
	log.Println("Created file: " + f.Name())

	unzippedDir := backend.UnzipToTempFodler(f.Name())
	if err := backend.InitLogDirectory(unzippedDir); err == nil {
		return unzippedDir
	} else {
		return ""
	}
}
func (b *App) OpenArchive() string {
	path, _ := runtime.OpenFileDialog(b.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: "",
		DefaultFilename:  "",
		Title:            "Open archive with logs",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "",
				Pattern:     "*.zip",
			},
		},
		ShowHiddenFiles:            false,
		CanCreateDirectories:       false,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	},
	)
	unzippedDir := backend.UnzipToTempFodler(path)
	log.Println("Unzipped files to path: " + path)
	if err := backend.InitLogDirectory(unzippedDir); err == nil {
		return unzippedDir
	} else {
		return ""
	}
}

func (b *App) GetLogs() string {
	logs := *backend.GetLogs()
	logsToDisplay := analyzer.Logs{}
	for i, entry := range logs {
		if entry.Visible == true {
			logsToDisplay = append(logsToDisplay, logs[i])
		}
	}
	html := logsToDisplay.ConvertToHTML()
	return html
}
func (b *App) GetStaticInfo() string {
	staticInfo := *backend.GetStaticInfo()
	html := staticInfo.ConvertToHTML()
	return html
}

// FilterGet returns the values of the filter area
func (b *App) GetFilters() string {
	return backend.GetFilters().ConvertToHTML()
}

func (b *App) GetThreadDumpFileContent(dir string, file string) string {
	return backend.GetThreadDumpFolder(dir).GetFile(file).ConvertToHTML()
}
func (b *App) GetThreadDumpsFilters(dir string) string {
	return backend.GetThreadDumpFolder(dir).GetFiltersHTML()
}

// FilterUpdate reads the values of filter area on the left of frontend window
func (b *App) FilterUpdate(value interface{}) string {
	return ""
}
func (b *App) SetFilters(a map[string]bool) string {
	err := backend.SetFilters(a)
	if err == nil {
		backend.GetLogs().ApplyFilters(backend.GetFilters())
	}
	return "failure"
}

func (b *App) GetEntityNamesWithLineHighlightingColors() string {
	jsonMap := make(map[string]string)
	for entityName, entityEntries := range *backend.GetFilters() {
		jsonMap[entityName] = entityEntries[0].GroupLineHighlightingColor
	}
	marshal, _ := json.Marshal(jsonMap)

	return string(marshal)
}
func ConvertDataURISchemeToBase64File(DataURIScheme string) (data []byte) {
	b64data := DataURIScheme[strings.IndexByte(DataURIScheme, ',')+1:]
	if data, err := base64.StdEncoding.DecodeString(b64data); err == nil {
		return data
	} else {
		log.Println("Could not convert Data URI scheme to base64 file")
		return nil
	}
}
