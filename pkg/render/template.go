package render

import (
	"bytes"
	"crypto/rsa"
	"html/template"
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
	funcs["randomKubernetesName"] = func() string {
		return util.String(8)
	}
	funcs["encryptWithBase64PKIXPublicKey"] = func(pubData, data string) string {
		var err error
		var pub *rsa.PublicKey
		var encryptedData []byte

		pub, err = Base64PKIXPublicKey([]byte(pubData))
		if err == nil {
			if encryptedData, err = EncryptWithPublicKey([]byte(data), pub); err == nil {
				return string(encryptedData)
			}
		}
		return err.Error()
	}
	return funcs
}
