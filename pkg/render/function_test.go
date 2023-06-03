package render_test

import (
	"testing"

	"github.com/linuxsuren/api-testing/pkg/render"
	"github.com/stretchr/testify/assert"
)

func TestEncryptWithPublicKey(t *testing.T) {
	pub, err := render.Base64PKIXPublicKey([]byte(`MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqjNKsGtWsOVGz0cgCEg5dBGgD03fmgXHfClp/C7yRho6Fopyhpz1Jga0rVZsavdt0ULmFWh6XcCC4TML11VVAneAaNXzBb+XRlhfKQeZn5D7K3EHISsj90gcV6jMADOFLvvlups/R0tHhnhbn4+WWJosLRFrm6INPNQRaXPOJsDT+VDoTNNbp71Dlemdg8oNk2abQUnfDY6WzF5gES2ZYS/9eW8qFsA56fdKbOiOAUs9oxVT3Ouyr++VYA/VDLjztGKPSDRvThpEbY7AQe+L2pTkZYwryoLNlvbahAvTTGyCi8bsjrcuZ1/lmjy41qo1LiogCzcwxDQBMkSLIKN55wIDAQAB`))
	if !assert.NoError(t, err) || !assert.NotNil(t, pub) {
		return
	}

	result, err := render.EncryptWithPublicKey([]byte("hello"), pub)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	pub, err = render.BytesToPublicKey([]byte(`-----BEGIN PUBLIC KEY-----
MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAK9rC1d/Z52H2SuN37HE7agkpnIC8V8P
761G516zkPq8YepsjdTqGLa872G0zFsQmOCYQ17j/GHDkJL46mrAXQ0CAwEAAQ==
-----END PUBLIC KEY-----`))
	assert.NoError(t, err)
	assert.NotNil(t, pub)
}
