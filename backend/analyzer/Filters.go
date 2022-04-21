package analyzer

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"log"
	"reflect"
	"sort"
	"strconv"
	"text/template"
)

type Filters map[string]struct {
	State   string
	Entries FilterEntries
}
type FilterEntries []FilterEntry

type FilterEntry struct {
	Checked                    bool
	GroupLabel                 string
	ID                         string
	EntryLabel                 string
	GroupLineHighlightingColor string //Fills automatically based on the LineHighlightingColor of the DynamicEntity
	//innerFilters []innerFilter
}

//type innerFilter struct {
//
//}
// func AddInnerFilter(){}

func (f *Filters) IsEmpty() bool {
	return reflect.ValueOf(*f).IsZero()
}
func (f Filters) ConvertToHTML() string {
	f.setFiltersGroupsState()
	var tpl bytes.Buffer
	t := template.Must(template.New("Filters.gohtml").
		ParseFS(tmplFS, "Filters.gohtml"))
	err := t.Execute(&tpl, f)
	if err != nil {
		log.Println(err.Error())
	}
	return tpl.String()
}

// Append adds generated filter to slice of filters
func (f Filters) Append(entity DynamicEntity, entityEntryPath string) {
	fi, _ := f[entity.Name]
	fi.Entries = append(fi.Entries, FilterEntry{
		Checked:                    entity.entityInstances[entityEntryPath].Visible,
		ID:                         entity.entityInstances[entityEntryPath].Hash,
		GroupLabel:                 entity.Name,
		GroupLineHighlightingColor: entity.LineHighlightingColor,
		EntryLabel:                 entity.GetDisplayName(entityEntryPath),
	})
	f[entity.Name] = fi
}

func (f *Filters) SortByFilename() {
	for _, filterEntries := range *f {
		sort.Slice(filterEntries.Entries, func(i, j int) bool {
			return sortName(filterEntries.Entries[i].EntryLabel) < sortName(filterEntries.Entries[j].EntryLabel)
		})
	}
}
func (f *Filters) getEntriesWithStates() (filtersList map[string]bool) {
	filtersList = make(map[string]bool)
	for _, entries := range *f {
		for _, entry := range entries.Entries {
			filtersList[entry.ID] = entry.Checked
		}
	}
	return filtersList
}

func (f Filters) setFiltersGroupsState() {
	for name, fi := range f {
		i := 0
		for _, entry := range fi.Entries {
			if entry.Checked {
				i++
			}
		}
		if i == len(fi.Entries) {
			fi.State = "checked"
		} else if i == 0 {
			fi.State = "unchecked"
		} else {
			fi.State = "mixed"
		}
		f[name] = fi
	}
}
func sortName(filename string) string {
	name := filename
	// split numeric suffix
	i := len(name) - 1
	for ; i >= 0; i-- {
		if '0' > name[i] || name[i] > '9' {
			break
		}
	}
	i++
	// string numeric suffix to uint64 bytes
	// empty string is zero, so integers are plus one
	b64 := make([]byte, 64/8)
	s64 := name[i:]
	if len(s64) > 0 {
		u64, err := strconv.ParseUint(s64, 10, 64)
		if err == nil {
			binary.BigEndian.PutUint64(b64, u64+1)
		}
	}
	// prefix + numeric-suffix
	return name[:i] + string(b64)
}
