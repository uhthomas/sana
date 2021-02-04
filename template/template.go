package template

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/eknkc/amber"
)

var templates = map[string]*template.Template{}

func init() {
	t, err := amber.CompileDir("./", amber.DirOptions{".jade", true}, amber.DefaultOptions)
	if err != nil {
		panic(err)
	}

	templates = t
}

func Render(res http.ResponseWriter, name string, data interface{}) {
	t, ok := templates[name]
	if !ok {
		res.WriteHeader(404)
		res.Write([]byte(`404 not found`))
		return
	}

	if err := t.Execute(res, data); err != nil {
		res.WriteHeader(500)
		res.Write([]byte(`500 failed to compile template`))
		return
	}
}

func CustomHandler(name string, data interface{}) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		Render(res, filepath.Join("template", name, "index"), data)
	}
}

func HttpHandler(res http.ResponseWriter, req *http.Request) {
	name := filepath.Join("template", req.URL.Path, "index")
	Render(res, name, nil)
}
