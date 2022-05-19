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
		log.Println(err.Error())
	}
	return tpl.String()
}
