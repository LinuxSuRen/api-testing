package render

import (
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
	}, {
		name: "encryptWithBase64PKIXPublicKey",
		text: `{{encryptWithBase64PKIXPublicKey "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqjNKsGtWsOVGz0cgCEg5dBGgD03fmgXHfClp/C7yRho6Fopyhpz1Jga0rVZsavdt0ULmFWh6XcCC4TML11VVAneAaNXzBb+XRlhfKQeZn5D7K3EHISsj90gcV6jMADOFLvvlups/R0tHhnhbn4+WWJosLRFrm6INPNQRaXPOJsDT+VDoTNNbp71Dlemdg8oNk2abQUnfDY6WzF5gES2ZYS/9eW8qFsA56fdKbOiOAUs9oxVT3Ouyr++VYA/VDLjztGKPSDRvThpEbY7AQe+L2pTkZYwryoLNlvbahAvTTGyCi8bsjrcuZ1/lmjy41qo1LiogCzcwxDQBMkSLIKN55wIDAQAB" "hello"}}`,
		ctx:  map[string]interface{}{},
		verify: func(t *testing.T, s string) {
			assert.NotEmpty(t, s, s)
		},
	}, {
		name: "encryptWithBase64PKIXPublicKey error",
		text: `{{encryptWithBase64PKIXPublicKey "fake" "hello"}}`,
		ctx:  map[string]interface{}{},
		verify: func(t *testing.T, s string) {
			assert.Equal(t, "asn1: syntax error: truncated tag or length", s)
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
