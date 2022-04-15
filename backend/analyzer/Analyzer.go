package analyzer

import (
	"crypto/sha1"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
)

var writeSyncer = sync.Mutex{}

//go:embed *.gohtml
var tmplFS embed.FS

type Analyzer struct {
	FolderToWorkWith      string
	IsFolderTemp          bool
	DynamicEntities       []DynamicEntity
	StaticEntities        []StaticEntity
	Filters               Filters
	OtherFiles            OtherFiles
	AggregatedLogs        Logs
	AggregatedThreadDumps AggregatedThreadDumps
	AggregatedStaticInfo  AggregatedStaticInfo
}
type StaticEntity struct {
	Name                string
	ConvertToStaticInfo func(path string) StaticInfo
	CheckPath           func(path string) bool
	CollectedInfo       StaticInfo
}

type DynamicEntity struct {
	entityInstances       map[string]string //entityInstances is path:hash map of every instance of entity created for every found path of this entity type.
	Name                  string            // Name of the Entity. For example "idea.log", "Thread dump", or "CPU snapshot". It will be used to group same entities.
	ConvertToLogs         func(path string) Logs
	CheckPath             func(path string) bool
	CheckIgnoredPath      func(path string) bool
	GetDisplayName        func(path string) string
	LineHighlightingColor string //Color represents the color that is used to highlight all lines of this entity type in the editor
}

func (e *DynamicEntity) addDynamicEntityInstance(path string) {
	if e.entityInstances == nil {
		e.entityInstances = make(map[string]string)
	}
	e.entityInstances[path] = getHash(path)
}

//AddStaticEntity adds new static Entity to the list of known Entities. Should be Called within the application start.
func (a *Analyzer) AddStaticEntity(entity StaticEntity) {
	a.StaticEntities = append(a.StaticEntities, entity)
}

//AddDynamicEntity adds new dynamic Entity to the list of known Entities. Should be Called within the application start.
func (a *Analyzer) AddDynamicEntity(entity DynamicEntity) {
	a.DynamicEntities = append(a.DynamicEntities, entity)
}

//ParseLogDirectory analyzes provided path for known log elements
func (a *Analyzer) ParseLogDirectory(path string) {
	var wg sync.WaitGroup
	var collectedFiles []string
	visit := func(path string, file os.DirEntry, err error) error {
		wg.Add(1)
		go func() {
			defer wg.Done()
			isDynamic := a.CollectLogsFromDynamicEntities(path)
			isStatic := a.CollectStaticInfoFromStaticEntities(path)
			writeSyncer.Lock()
			if isStatic || isDynamic {
				collectedFiles = append(collectedFiles, path)
			} else {
				if !file.IsDir() && !IsHiddenFile(filepath.Base(path)) {
					a.OtherFiles.Append(path)
				}
			}
			writeSyncer.Unlock()
		}()
		return nil
	}
	_ = filepath.WalkDir(path, visit)
	wg.Wait()
	a.OtherFiles = a.OtherFiles.FilterAnalyzedDirectories(collectedFiles)
}

//IsEmpty checks if config has at least one filled attribute
func (a *Analyzer) IsEmpty() bool {
	return reflect.ValueOf(*a).IsZero()
}

func (a *Analyzer) GetLogs() *Logs {
	if !a.AggregatedLogs.IsEmpty() {
		return &a.AggregatedLogs
	}
	return nil
}

func (a *Analyzer) GetStaticInfo() *AggregatedStaticInfo {
	if !(len(a.AggregatedStaticInfo) == 0) {
		return &a.AggregatedStaticInfo
	}
	a.AggregatedStaticInfo = aggregateStaticInfo(a.StaticEntities)
	return a.GetStaticInfo()
}
func (a *Analyzer) GetOtherFiles() *OtherFiles {
	if !a.AggregatedLogs.IsEmpty() {
		return &a.OtherFiles
	}
	return nil
}

//GetThreadDump returns Analyzed ThreadDumps folder as AggregatedThreadDumps entity. Analyzes it if it was not done already.
func (a *Analyzer) GetThreadDump(threadDumpsFolder string) *ThreadDump {
	t := a.AggregatedThreadDumps[threadDumpsFolder]
	if t != nil {
		return &t
	}
	a.AggregatedThreadDumps[threadDumpsFolder] = make(ThreadDump)
	a.AggregatedThreadDumps[threadDumpsFolder] = analyzeThreadDumpsFolder(a.FolderToWorkWith, threadDumpsFolder)
	return a.GetThreadDump(threadDumpsFolder)
}
func (a *Analyzer) GetFilters() *Filters {
	if !a.Filters.IsEmpty() {
		return &a.Filters
	}
	return nil
}

func (a *Analyzer) CollectStaticInfoFromStaticEntities(path string) (analyzed bool) {
	analyzed = false
	for i, entity := range a.StaticEntities {
		if entity.CheckPath(path) == true {
			a.StaticEntities[i].CollectedInfo = entity.ConvertToStaticInfo(path)
			analyzed = true
		}
	}
	return analyzed
}

// CollectLogsFromDynamicEntities Checks if path fulfil the Entity requirements and Adds all the Entity's logEntries to the aggregated logs
func (a *Analyzer) CollectLogsFromDynamicEntities(path string) (analyzed bool) {
	analyzed = false
	for i, entity := range a.DynamicEntities {
		if entity.CheckIgnoredPath != nil {
			if entity.CheckIgnoredPath(path) == true {
				return true
			}
		}
		if entity.CheckPath(path) == true {
			logEntries := entity.ConvertToLogs(path)
			if logEntries == nil {
				log.Printf("Entity \"%s\" returned nothing for %s. Adding file to other files", entity.Name, path)
			} else {
				writeSyncer.Lock()
				a.DynamicEntities[i].addDynamicEntityInstance(path)
				a.AggregatedLogs.AppendSeveral(a.DynamicEntities[i].Name, a.DynamicEntities[i].entityInstances[path], logEntries)
				writeSyncer.Unlock()
				analyzed = true
			}
		}
	}
	return analyzed
}

// GenerateFilters Generates filters for all Entities and saves them into Filters slice
func (a *Analyzer) GenerateFilters() {
	filter := a.InitFilter()
	for _, entity := range a.DynamicEntities {
		for path, id := range entity.entityInstances {
			filter.Append(entity, path, id)
		}
	}
	filter.SortByFilename()
}

func (a *Analyzer) InitFilter() *Filters {
	a.Filters = make(Filters)
	return &a.Filters
}

func (a *Analyzer) Clear() {
	a.AggregatedLogs = Logs{}
	a.Filters = Filters{}
	a.OtherFiles = OtherFiles{}
	a.AggregatedStaticInfo = AggregatedStaticInfo{}
	a.AggregatedThreadDumps = AggregatedThreadDumps{}
	for i, _ := range a.StaticEntities {
		a.StaticEntities[i].CollectedInfo = StaticInfo{}
	}
	for i, _ := range a.DynamicEntities {
		a.DynamicEntities[i].entityInstances = make(map[string]string)
	}
	if a.IsFolderTemp {
		err := os.RemoveAll(a.FolderToWorkWith)
		if err != nil {
			log.Println(err)
		}
	}
}

func (a *Analyzer) GetThreadDumps(dir string) Logs {
	for _, entity := range a.DynamicEntities {
		for path, _ := range entity.entityInstances {
			if strings.Contains(path, dir) {
				return entity.ConvertToLogs(path)
			}
		}
	}
	return nil
}

func aggregateStaticInfo(entity []StaticEntity) (a AggregatedStaticInfo) {
	a = make(AggregatedStaticInfo)
	for _, staticEntity := range entity {
		a[staticEntity.Name] = staticEntity.CollectedInfo
	}
	return a
}

func getHash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	sh := fmt.Sprintf("%x\n", bs)
	return sh
}

/**
 * Parses string s with the given regular expression and returns the
 * group values defined in the expression.
 *
 */
func GetRegexNamedCapturedGroups(regEx, s string) (paramsMap map[string]string) {

	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(s)
	paramsMap = make(map[string]string)

	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
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

func SliceContains[S comparable](slice []S, element S) int {
	for i, e := range slice {
		if e == element {
			return i
		}
	}
	return -1
}

func IsHiddenFile(filename string) bool {
	if runtime.GOOS != "windows" {
		return filename[0:1] == "."
	}
	return false
}
