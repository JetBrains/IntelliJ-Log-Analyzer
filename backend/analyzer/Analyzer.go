package analyzer

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

var writeSyncer = sync.Mutex{}

type Analyzer struct {
	FolderToWorkWith     string
	IsFolderTemp         bool
	DynamicEntities      []DynamicEntity
	StaticEntities       []StaticEntity
	Filters              Filters
	AggregatedLogs       Logs
	AggregatedStaticInfo AggregatedStaticInfo
}
type StaticEntity struct {
	Name                string
	ConvertToStaticInfo func(path string) StaticInfo
	CheckPath           func(path string) bool
	CollectedInfo       StaticInfo
}

type DynamicEntity struct {
	entityInstances map[string]string //entityInstances is path:hash map of every instance of entity created for every found path of this entity type.
	Name            string            // Name of the Entity. For example "idea.log", "Thread dump", or "CPU snapshot". It will be used to group same entities.
	ConvertToLogs   func(path string) Logs
	CheckPath       func(path string) bool
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
	visit := func(path string, file os.DirEntry, err error) error {
		wg.Add(1)
		go func() {
			defer wg.Done()
			a.CollectLogsFromDynamicEntities(path)
			a.CollectStaticInfoFromStaticEntities(path)
		}()
		return nil
	}
	_ = filepath.WalkDir(path, visit)
	wg.Wait()
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

func (a *Analyzer) GetFilters() *Filters {
	if !a.Filters.IsEmpty() {
		return &a.Filters
	}
	return nil
}

func (a *Analyzer) CollectStaticInfoFromStaticEntities(path string) {
	for i, entity := range a.StaticEntities {
		if entity.CheckPath(path) == true {
			a.StaticEntities[i].CollectedInfo = entity.ConvertToStaticInfo(path)
		}
	}
}

// CollectLogsFromDynamicEntities Checks if path fulfil the Entity requirements and Adds all the Entity's logEntries to the aggregated logs
func (a *Analyzer) CollectLogsFromDynamicEntities(path string) {
	for i, entity := range a.DynamicEntities {
		if entity.CheckPath(path) == true {
			logEntries := entity.ConvertToLogs(path)
			writeSyncer.Lock()
			a.DynamicEntities[i].addDynamicEntityInstance(path)
			a.AggregatedLogs.AppendSeveral(a.DynamicEntities[i].entityInstances[path], logEntries)
			writeSyncer.Unlock()
		}
	}
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
	for i, _ := range a.DynamicEntities {
		a.DynamicEntities[i].entityInstances = make(map[string]string)
	}
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
