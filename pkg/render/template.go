package render

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/linuxsuren/api-testing/pkg/util"
)

// Render render then return the result
func Render(name, text string, ctx interface{}) (result string, err error) {
	var tpl *template.Template
	if tpl, err = template.New(name).
		Funcs(sprig.FuncMap()).
		Funcs(template.FuncMap{
			"randomKubernetesName": func() string {
				return util.String(8)
			},
		}).Parse(text); err == nil {
		buf := new(bytes.Buffer)
		if err = tpl.Execute(buf, ctx); err == nil {
			result = strings.TrimSpace(buf.String())
		}
	}
	return
}
