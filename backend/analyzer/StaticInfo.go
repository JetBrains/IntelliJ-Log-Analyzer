package analyzer

import (
	"bytes"
	"log"
	"text/template"
)

type StaticInfo struct {
	IDE         string
	Build       string
	JRE         string
	OS          string
	PluginsList []IDEPlugin
}
type IDEPlugin struct {
	Version string
	Name    string
	Link    string
}

//AggregatedStaticInfo is a source (such as troubleshooting.txt or idea.log) mapped to collected info
type AggregatedStaticInfo map[string]StaticInfo

//ConvertToHTML Represents logs as HTML based on Logs.gohtml template
func (a AggregatedStaticInfo) ConvertToHTML() string {
	//todo: combine static info from several sources
	if a.IsEmpty() {
		return ""
	}
	var tpl bytes.Buffer
	t := template.Must(template.New("StaticInfo.gohtml").
		ParseFS(tmplFS, "StaticInfo.gohtml"))
	err := t.Execute(&tpl, a)
	if err != nil {
		log.Println(err.Error())
	}
	return tpl.String()
}

func (a *AggregatedStaticInfo) IsEmpty() bool {
	for _, info := range *a {
		if len(info.IDE) > 0 || len(info.Build) > 0 || len(info.JRE) > 0 {
			return false
		}
	}
	return true
}
