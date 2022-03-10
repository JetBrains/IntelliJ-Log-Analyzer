package analyzer

import (
	"bytes"
	"html/template"
	"log"
)

type StaticInfo struct {
	IDE         string
	Build       string
	JRE         string
	PluginsList []IDEPlugin
}
type IDEPlugin struct {
	Version string
	Name    string
	Link    string
}
type AggregatedStaticInfo map[string]StaticInfo

//ConvertToHTML Represents logs as HTML based on Logs.gohtml template
func (a AggregatedStaticInfo) ConvertToHTML() string {
	//todo: combine static info from several sources
	var tpl bytes.Buffer
	t := template.Must(template.New("StaticInfo.gohtml").
		ParseFiles("./backend/analyzer/StaticInfo.gohtml"))
	err := t.Execute(&tpl, a)
	if err != nil {
		log.Println(err.Error())
	}
	return tpl.String()
}
