package render

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		ctx    interface{}
		expect string
		verify func(*testing.T, string)
	}{{
		name:   "default",
		text:   `{{default "hello" .Bar}}`,
		ctx:    nil,
		expect: "hello",
	}, {
		name:   "trim",
		text:   `{{trim "   hello    "}}`,
		ctx:    "",
		expect: "hello",
	}, {
		name: "randomKubernetesName",
		text: `{{randomKubernetesName}}`,
		verify: func(t *testing.T, s string) {
			assert.Equal(t, 8, len(s))
		},
	}, {
		name:   "md5",
		text:   `{{md5 "linuxsuren"}}`,
		expect: "b559b80ae1ba1c292d9b3265f265e76a",
	}, {
		name:   "base64",
		text:   `{{base64 "linuxsuren"}}`,
		expect: "bGludXhzdXJlbg==",
	}, {
		name:   "base64Decode",
		text:   `{{base64Decode "bGludXhzdXJlbg=="}}`,
		expect: "linuxsuren",
	}, {
		name:   "base64Decode with error",
		text:   `{{base64Decode "error"}}`,
		expect: "illegal base64 data at input byte 4",
	}, {
		name: "complex",
		text: `{{(index .items 0).name}}?a=a&key={{randomKubernetesName}}`,
		ctx: map[string]interface{}{
			"items": []interface{}{map[string]string{
				"name": "one",
			}, map[string]string{
				"name": "two",
			}},
		},
		verify: func(t *testing.T, s string) {
			assert.Equal(t, 20, len(s), s)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Render(tt.name, tt.text, tt.ctx)
			assert.Nil(t, err, err)
			if tt.expect != "" {
				assert.Equal(t, tt.expect, result)
			}
			if tt.verify != nil {
				tt.verify(t, result)
			}
		})
	}
}

func TestRenderThenPrint(t *testing.T) {
	tests := []struct {
		name    string
		tplText string
		ctx     interface{}
		buf     *bytes.Buffer
		expect  string
	}{{
		name:    "simple",
		tplText: `{{max 1 2 3}}`,
		ctx:     nil,
		buf:     new(bytes.Buffer),
		expect:  `3`,
	}, {
		name:    "with a map as context",
		tplText: `{{.name}}`,
		ctx:     map[string]string{"name": "linuxsuren"},
		buf:     new(bytes.Buffer),
		expect:  "linuxsuren",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RenderThenPrint(tt.name, tt.tplText, tt.ctx, tt.buf)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, tt.buf.String())
		})
	}
}

func TestFuncGenerator(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		fields   string
		expect   string
	}{{
		name:     "randomKubernetesName",
		funcName: "randomKubernetesName",
		expect:   `{{randomKubernetesName}}`,
	}, {
		name:     "generateJSONString",
		funcName: "generateJSONString",
		fields:   "name, age",
		expect:   `{{generateJSONString "name" "age"}}`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funcs := GetAdvancedFuncs()
			for _, f := range funcs {
				if f.FuncName == tt.funcName {
					buf := new(bytes.Buffer)
					ctx := context.Background()
					ctx = context.WithValue(ctx, ContextBufferKey, buf)
					err := f.Generator(ctx, tt.fields)
					assert.NoError(t, err)
					assert.Equal(t, tt.expect, buf.String())
				}
			}
		})
	}
}

func TestGoDogGenerator(t *testing.T) {
	tests := []struct {
		name       string
		goDogExper string
		fields     string
		expect     string
	}{{
		name:       "randomKubernetesName",
		goDogExper: `^生成随机字符串，长度 (.*)$`,
		fields:     `3`,
		expect:     `{{randAlpha 3}}`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funcs := GetAdvancedFuncs()
			for _, f := range funcs {
				if f.GoDogExper == tt.goDogExper {
					buf := new(bytes.Buffer)
					ctx := context.Background()
					ctx = context.WithValue(ctx, ContextBufferKey, buf)
					err := f.Generator(ctx, tt.fields)
					assert.NoError(t, err)
					assert.Equal(t, tt.expect, buf.String())
				}
			}
		})
	}
}

func TestGenerateJSONString(t *testing.T) {
	result := generateJSONString([]string{"name", "age"})
	assert.Equal(t, `{"age":"random","name":"random"}`, result)
}

func TestSecret(t *testing.T) {
	SetSecretGetter(nil)
	result, err := Render("", `{{secretValue "pass"}}`, nil)
	assert.NoError(t, err)
	assert.Equal(t, "no secret server", result)

	expected := "password"
	SetSecretGetter(&nonSecretGetter{
		value: expected,
	})
	result, err = Render("", `{{secretValue "pass"}}`, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
