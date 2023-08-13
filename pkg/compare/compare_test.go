package compare

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestElement(t *testing.T) {
	exp := `{
		"data": [
		  {
			"key": "hell",
			"value": "func() strin"
		  }
		]
	  }
	`
	act := `
	  {
		"data": [
		  {
			"key": "hello",
			"value": "func() string"
		  }
		]
	  }`
	expect := gjson.Parse(exp)
	actul := gjson.Parse(act)

	err := Element("TestElement", expect, actul)

	expmsg1 := "compare: field TestElement.data.0.value: expect func() strin but got func() string"
	expmsg2 := "compare: field TestElement.data.0.key: expect hell but got hello"
	assert.Contains(t, err.Error(), expmsg1)
	assert.Contains(t, err.Error(), expmsg2)
}
