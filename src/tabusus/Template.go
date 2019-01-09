package tabusus

import (
	"github.com/labstack/echo"
	"html/template"
	"io"
	"strings"
)

// Ref: https://echo.labstack.com/guide/templates

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	directory          string
	templateFileSuffix string
	templates          map[string]*template.Template
}

func newTemplateRenderer(directory, templateFileSuffix string) *TemplateRenderer {
	return &TemplateRenderer{
		directory:          directory,
		templateFileSuffix: templateFileSuffix,
		templates:          map[string]*template.Template{},
	}
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	sess := getSession(c)
	flash := sess.Flashes()
	sess.Save(c.Request(), c.Response())

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
		viewContext["static"] = staticPath
		viewContext["appInfo"] = AppConfig.Conf.GetConfig("app")
		if len(flash) > 0 {
			viewContext["flash"] = flash[0].(string)
		}
	}
	tpl := t.templates[name]
	tokens := strings.Split(name, ":")
	if tpl == nil {
		var files []string
		for _, v := range tokens {
			files = append(files, t.directory+"/"+v+".html")
		}
		tpl = template.Must(template.New(name).ParseFiles(files...))
		t.templates[name] = tpl
	}
	return tpl.ExecuteTemplate(w, tokens[0]+".html", data)
}
