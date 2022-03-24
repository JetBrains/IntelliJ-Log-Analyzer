package analyzer

import (
	"bytes"
	"html/template"
	"io/ioutil"
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
		ParseFS(tmplFS, "OtherFiles.gohtml"))
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
		if i := SliceContains(collectedFiles, filepath.Dir(file.FullPath)); i == -1 {
			s = append(s, file)
		}
	}
	return s
}

func (f *OtherFiles) GetContent(fileUUID string) string {
	for _, file := range *f {
		if file.Uuid == fileUUID {
			content, _ := ioutil.ReadFile(file.FullPath)
			return string(content)
		}
	}
	return ""
}
