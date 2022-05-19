package installedIDEs

import (
	"bytes"
	"embed"
	"log"
	"text/template"
)

//go:embed *.gohtml
var tmplFS embed.FS

func GetInstalledIDEsDropdownHTML() string {
	installedIDEs := GetIdeInstallations()
	var tpl bytes.Buffer
	t := template.Must(template.New("IdeDropdown.gohtml").
		ParseFS(tmplFS, "IdeDropdown.gohtml"))
	err := t.Execute(&tpl, installedIDEs)
	if err != nil {
		log.Printf("Template IdeDropdown.gohtml parsing failed. Error: %s", err.Error())
	}
	return tpl.String()
}
