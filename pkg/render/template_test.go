package render

import (
	"bytes"
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
