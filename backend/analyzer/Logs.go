package analyzer

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"
)

type Logs []LogEntry

type LogEntry struct {
	entityInstanceId string //Log path
	Severity         string
	Time             time.Time
	Text             string
	Visible          bool
}

func (logs *Logs) ConvertToJSON() string {
	logsAsJson, err := json.Marshal(logs)
	if err != nil {
		log.Println(err.Error())
	}
	return string(logsAsJson)
}

//ConvertToHTML Represents logs as HTML based on Logs.gohtml template
func (logs Logs) ConvertToHTML() string {
	var tpl bytes.Buffer
	t := template.Must(template.New("Logs.gohtml").
		Funcs(
			template.FuncMap{
				"replaceNewline": func(s string) template.HTML {
					return template.HTML(strings.Replace(template.HTMLEscapeString(s), "\n", "</p><p>", -1))
				},
			}).
		ParseFiles("./backend/analyzer/Logs.gohtml"))
	err := t.Execute(&tpl, logs)
	if err != nil {
		log.Println(err.Error())
	}
	return tpl.String()
}

// Append adds one log entry of entity with UUID to the struct of logs
func (logs *Logs) Append(instanceId string, entry LogEntry) {
	entry.entityInstanceId = instanceId
	entry.Visible = true
	*logs = append(*logs, entry)
}
func (logs *Logs) AppendSeveral(instanceId string, logEntry []LogEntry) {
	for _, entry := range logEntry {
		logs.Append(instanceId, entry)
	}

}

func (logs *Logs) IsEmpty() bool {
	return reflect.ValueOf(*logs).IsZero()
}

func (logs Logs) SortByTime() {
	sort.Slice(logs, func(i, j int) bool { return logs[i].Time.Before(logs[j].Time) })
}

func (logs Logs) ApplyFilters(filters *Filters) {
	filtersList := filters.getEntriesWithStates()
	for i, entry := range logs {
		a := filtersList[entry.entityInstanceId]
		logs[i].Visible = a
	}
}
