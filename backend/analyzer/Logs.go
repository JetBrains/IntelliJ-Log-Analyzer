package analyzer

import (
	"bytes"
	"log"
	"reflect"
	"sort"
	"text/template"
	"time"
)

type Logs []LogEntry

type LogEntry struct {
	EntityInstanceId string //Fills automatically. Log path
	EntityName       string //Fills automatically. Type of Dynamic Entity (idea log, Thread Dump, etc)
	Severity         string
	Time             time.Time
	Text             string
	Visible          bool
}

//ConvertToHTML Represents logs as HTML based on Logs.gohtml template
func (logs Logs) ConvertToHTML() string {
	var tpl bytes.Buffer
	t := template.Must(template.New("Logs.gohtml").
		ParseFS(tmplFS, "Logs.gohtml"))
	err := t.Execute(&tpl, logs)
	if err != nil {
		log.Println(err.Error())
	}
	return tpl.String()
}

// Append adds one log entry of entity with UUID to the struct of logs
func (logs *Logs) Append(entityName string, instanceId string, entry LogEntry) {
	entry.EntityInstanceId = instanceId
	entry.EntityName = entityName
	entry.Visible = true
	*logs = append(*logs, entry)
}
func (logs *Logs) AppendSeveral(entityName string, instanceId string, logEntry []LogEntry) {
	for _, entry := range logEntry {
		logs.Append(entityName, instanceId, entry)
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
		a := filtersList[entry.EntityInstanceId]
		logs[i].Visible = a
	}
}
