package render

import (
	"github.com/linuxsuren/api-testing/pkg/secret"
)

type nonSecretGetter struct {
	value string
	err   error
}

func (n *nonSecretGetter) GetSecret(name string) (s secret.Secret, err error) {
	s.Value = n.value
	s.Name = name
	err = n.err
	return
}
