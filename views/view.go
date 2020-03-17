package views

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")

	switch data.(type) {
	case Data:
	default:
		data = Data{
			Yield: data,
		}
	}

	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)

	files = append(files, layoutFiles()...)
	fmt.Println(files)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}
