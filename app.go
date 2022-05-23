package main

import (
	"context"
	"encoding/base64"
	json "encoding/json"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"log_analyzer/backend"
	"log_analyzer/backend/analyzer"
	"log_analyzer/backend/analyzer/entities"
	"log_analyzer/backend/analyzer/installedIDEs"
	"log_analyzer/backend/update"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
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
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc // -> this is setting the global timezone
	b.ctx = ctx
	b.RenderSystemMenu()
	b.CheckForUpdates()
}

// domReady is called after the front-end dom has been loaded
func (b *App) domReady(ctx context.Context) {
	// Add your action here
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
	entities.CurrentAnalyzer.Clear()
}

func (b *App) OpenIndexingSummaryForProject(fileName string) {
	absolutePath := backend.GetIndexingFilePath(fileName)
	wailsruntime.BrowserOpenURL(b.ctx, filepath.Dir(absolutePath)+string(filepath.Separator)+"report.html")
}
func (b *App) OpenIndexingReport(fileName string) {
	absolutePath := backend.GetIndexingFilePath(fileName)
	wailsruntime.BrowserOpenURL(b.ctx, absolutePath)
}
func (b *App) OpenFolder() string {
	path, _ := wailsruntime.OpenDirectoryDialog(b.ctx, wailsruntime.OpenDialogOptions{
		DefaultDirectory: "",
		DefaultFilename:  "",
		Title:            "Open direcotry with logs",
		Filters: []wailsruntime.FileFilter{
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
	return path
}
func (b *App) UploadLogFile(filename string, DataURIScheme string) string {
	data := ConvertDataURISchemeToBase64File(DataURIScheme)
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
	return b.InitLogDirectory(f.Name())
}
func (b *App) InitLogDirectory(path string) string {
	if err := backend.InitLogDirectory(path, &b.ctx); err == nil {
		return path
	} else {
		return ""
	}
}
func (b *App) UploadArchive(DataURIScheme string) string {
	data := ConvertDataURISchemeToBase64File(DataURIScheme)
	f, err := os.CreateTemp("", "IntelliJLogsAnalyzer-temp.zip")
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
	err = os.RemoveAll(f.Name())
	if err != nil {
		log.Println("Could not remove temp archive: " + err.Error())
	}
	return b.InitLogDirectory(unzippedDir)
}
func (b *App) OpenArchive() string {
	path, _ := wailsruntime.OpenFileDialog(b.ctx, wailsruntime.OpenDialogOptions{
		DefaultDirectory: "",
		DefaultFilename:  "",
		Title:            "Open archive with logs",
		Filters: []wailsruntime.FileFilter{
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
	if path == "" {
		return ""
	}
	unzippedDir := backend.UnzipToTempFodler(path)
	log.Println("Unzipped files to path: " + path)
	return unzippedDir
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

func (b *App) GetSummary() string {
	return backend.GetFilters().ConvertToHTML() + backend.GetOtherFiles().ConvertToHTML()
}

// FilterGet returns the values of the filter area
func (b *App) GetFilters() string {
	return backend.GetFilters().ConvertToHTML()
}

func (b *App) GetThreadDumpFileContent(dir string, file string) string {
	return backend.GetThreadDumpFolder(dir).GetFile(file).ConvertToHTML()
}
func (b *App) GetOtherFileContent(fileUUID string) string {
	return backend.GetOtherFiles().GetContent(fileUUID)
}

//GetThreadDumpsFilters returns HTML of the list of files in ThreadDump folder.
func (b *App) GetThreadDumpsFilters(dir string) string {
	return backend.GetThreadDumpFolder(dir).GetFiltersHTML()
}

// SetFilters reads the values of filter area on the left of frontend window
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
		jsonMap[entityName] = entityEntries.Entries[0].GroupLineHighlightingColor
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

func (b *App) CheckForUpdates() bool {
	updateAvailable, version := update.CheckForUpdate(Version)
	if updateAvailable {
		log.Printf("Update to version %s available. ", version)
		dialog, _ := wailsruntime.MessageDialog(b.ctx, wailsruntime.MessageDialogOptions{
			Type:          "QuestionDialog",
			Title:         "Update available",
			Message:       "Version " + version + " is available",
			Buttons:       []string{"Update", "Cancel"},
			DefaultButton: "Update",
			CancelButton:  "Cancel",
		})
		if dialog == "Update" {
			var updated bool
			if runtime.GOOS == "darwin" {
				updated = update.DoSelfUpdateMac()
			} else {
				updated = update.DoSelfUpdate(Version)
			}
			if updated {
				wailsruntime.MessageDialog(b.ctx, wailsruntime.MessageDialogOptions{
					"InfoDialog",
					"Update installed",
					"Update " + version + " is installed. Please restart the program to apply changes.",
					[]string{"Ok", "Cancel"},
					"OK",
					"Cancel",
					nil,
				})
			}
		}
		return true
	} else {
		log.Printf("Version %s is the latest", Version)
		return false
	}
}

func (b *App) ShowNoUpdatesMessage() {
	wailsruntime.MessageDialog(b.ctx, wailsruntime.MessageDialogOptions{
		"InfoDialog",
		"No updates available",
		Version + " is the latest version",
		[]string{"Ok"},
		"",
		"",
		nil,
	})
}

func (b *App) RenderSystemMenu() {
	macMenu := menu.NewMenuFromItems(
		menu.AppMenu(),
		menu.SubMenu("File", menu.NewMenuFromItems(
			menu.Text("Check for updates", nil, func(_ *menu.CallbackData) {
				if !b.CheckForUpdates() {
					b.ShowNoUpdatesMessage()
				}
			}),
			menu.Text("Settings", keys.CmdOrCtrl(","), func(_ *menu.CallbackData) {
				wailsruntime.EventsEmit(b.ctx, "ShowSettings")
			}),
		)),
		menu.EditMenu(),
		menu.SubMenu("Help", menu.NewMenuFromItems(
			menu.Text("Submit a Bug Report", keys.CmdOrCtrl("b"), func(_ *menu.CallbackData) {
				wailsruntime.BrowserOpenURL(b.ctx, "https://github.com/annikovk/IntelliJ-Log-Analyzer/issues/new")
			}),
		)),
	)
	windowsMenu := menu.NewMenuFromItems(
		menu.SubMenu("File", menu.NewMenuFromItems(
			menu.Text("Check for updates", nil, func(_ *menu.CallbackData) {
				if !b.CheckForUpdates() {
					b.ShowNoUpdatesMessage()
				}
			}),
			menu.Text("Settings", keys.Combo("s", keys.ControlKey, keys.OptionOrAltKey), func(_ *menu.CallbackData) {
				wailsruntime.EventsEmit(b.ctx, "ShowSettings")
			}),
		)),
		menu.SubMenu("Help", menu.NewMenuFromItems(
			menu.Separator(),
			menu.Text("Submit Bug", keys.CmdOrCtrl("b"), func(_ *menu.CallbackData) {
				wailsruntime.BrowserOpenURL(b.ctx, "https://github.com/annikovk/IntelliJ-Log-Analyzer/issues/new")
			}),
			menu.Text("Check for updates", keys.CmdOrCtrl("u"), func(_ *menu.CallbackData) {
				if !b.CheckForUpdates() {
					b.ShowNoUpdatesMessage()
				}
			}),
		)),
	)
	if runtime.GOOS == "darwin" {
		wailsruntime.MenuSetApplicationMenu(b.ctx, macMenu)
	} else {
		wailsruntime.MenuSetApplicationMenu(b.ctx, windowsMenu)
	}
}

func (b *App) GetSettingsScreenHTML() string {
	return backend.GetSettingsScreenHTML()
}
func (b *App) GetRunningIDEsDropdownHTML() string {
	return installedIDEs.GetInstalledIDEsDropdownHTML()
}
func (b *App) SaveSetting(key string, value interface{}) {
	backend.GetConfig().SaveSetting(key, value)
	wailsruntime.EventsEmit(b.ctx, "SettingsChanged", backend.GetConfig())
	log.Println(b.ctx)
}
func (b *App) GetSetting(key string) interface{} {
	ptr := reflect.ValueOf(backend.GetConfig())
	s := reflect.Indirect(ptr).FieldByName(key).Interface()
	return s
}
func (b *App) EnableLogsLiveUpdate() {
	backend.EnableLogsLiveUpdate()
}
