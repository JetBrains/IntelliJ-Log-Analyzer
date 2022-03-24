package analyzer

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"
)

type OtherFiles []struct {
	Uuid     string
	FullPath string
	BasePath string
}

func (f *OtherFiles) ConvertToHTML() (html string) {
	var tpl bytes.Buffer
	t := template.Must(template.New("OtherFiles.gohtml").
		ParseFiles("./backend/analyzer/OtherFiles.gohtml"))
	err := t.Execute(&tpl, *f)
	if err != nil {
		log.Println(err.Error())
	}
	return tpl.String()
}

func (f *OtherFiles) Append(path string) {
	*f = append(*f, struct {
		Uuid     string
		FullPath string
		BasePath string
	}{Uuid: getHash(path), FullPath: path, BasePath: filepath.Base(path)})
}

//FilterAnalyzedDirectories
func (f *OtherFiles) FilterAnalyzedDirectories(collectedFiles []string) OtherFiles {
	s := OtherFiles{}
	for _, file := range *f {
		if !sliceContains(collectedFiles, filepath.Dir(file.FullPath)) {
			s = append(s, file)
		}
	}
	return s
}
