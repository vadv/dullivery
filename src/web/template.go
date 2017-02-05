package web

import (
	"io"
	"path/filepath"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

type Template struct {
	StaticDir string
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tpl := pongo2.Must(pongo2.FromFile(filepath.Join(t.StaticDir, "views/layout.html")))
	return tpl.ExecuteWriter(pongo2.Context{"template": name, "Storage": Storage, "data": data}, w)
}

func LoadTemplates(staticDir string) *Template {
	return &Template{StaticDir: staticDir}
}
