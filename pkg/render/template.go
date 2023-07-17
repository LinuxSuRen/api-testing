package render

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/linuxsuren/api-testing/pkg/util"
)

// Render render then return the result
func Render(name, text string, ctx interface{}) (result string, err error) {
	var tpl *template.Template
	if tpl, err = template.New(name).
		Funcs(FuncMap()).
		Parse(text); err == nil {
		buf := new(bytes.Buffer)
		if err = tpl.Execute(buf, ctx); err == nil {
			result = strings.TrimSpace(buf.String())
		}
	}
	return
}

// FuncMap reutrns all the supported functions
func FuncMap() template.FuncMap {
	funcs := sprig.FuncMap()
	for _, item := range GetAdvancedFuncs() {
		if item.FuncName == "" || item.Func == nil {
			continue
		}
		funcs[item.FuncName] = item.Func
	}
	return funcs
}

// RenderThenPrint renders the template then prints the result
func RenderThenPrint(name, text string, ctx interface{}, w io.Writer) (err error) {
	var report string
	if report, err = Render(name, text, ctx); err == nil {
		fmt.Fprint(w, report)
	}
	return
}

var advancedFuncs = []AdvancedFunc{{
	FuncName:   "generateJSONString",
	Func:       generateJSONString,
	GoDogExper: `^生成对象，字段包含 (.*)$`,
	Generator:  generateJSONObject,
}, {
	FuncName: "randomKubernetesName",
	Func: func() string {
		return util.String(8)
	},
	GoDogExper: `^动态k8s名称(.*)$`,
	Generator: func(ctx context.Context, fields string) (err error) {
		writeWithContext(ctx, `{{randomKubernetesName}}`)
		return
	},
}, {
	GoDogExper: `^生成随机字符串，长度 (.*)$`,
	Generator: func(ctx context.Context, fields string) (err error) {
		writeWithContext(ctx, `{{randAlpha `+fields+`}}`)
		return
	},
}}

// GetAdvancedFuncs returns all the advanced functions
func GetAdvancedFuncs() []AdvancedFunc {
	return advancedFuncs
}

func generateJSONString(fields []string) (result string) {
	data := make(map[string]string)
	for _, item := range fields {
		data[item] = "random"
	}

	if json, err := json.Marshal(data); err == nil {
		result = string(json)
	}
	return
}

type ContextKey string

var ContextBufferKey ContextKey = "ContextBufferKey"

// generateJSONObject generates a json object
// For instance: {{generateJSONString "hello" "world"}}
func generateJSONObject(ctx context.Context, fields string) (err error) {
	items := strings.Split(fields, ",")
	funcExp := "{{generateJSONString"
	for _, item := range items {
		funcExp += " \"" + strings.TrimSpace(item) + "\""
	}
	funcExp += "}}"

	writeWithContext(ctx, funcExp)
	return
}

func writeWithContext(ctx context.Context, text string) {
	buf := ctx.Value(ContextBufferKey)
	writer, ok := buf.(io.Writer)
	if ok && writer != nil {
		_, _ = writer.Write([]byte(text))
	}
	return
}

// AdvancedFunc represents an advanced function
type AdvancedFunc struct {
	FuncName   string
	Func       interface{}
	GoDogExper string
	Generator  func(ctx context.Context, fields string) (err error)
}
