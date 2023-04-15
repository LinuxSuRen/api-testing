package kubernetes_test

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/linuxsuren/api-testing/pkg/runner/kubernetes"
	"github.com/stretchr/testify/assert"
)

func TestGetPod(t *testing.T) {
	tests := []struct {
		name           string
		namespacedName namespacedName
		prepare        func()
		expect         map[string]interface{}
	}{{
		name: "normal",
		namespacedName: namespacedName{
			namespace: "ns",
			name:      "fake",
		},
		prepare: func() {
			gock.New("http://foo").
				Get("/api/v1/namespaces/ns/pods/fake").
				Reply(http.StatusOK).
				JSON(`{"kind":"pod"}`)
			gock.InterceptClient(kubernetes.GetClient())
		},
		expect: map[string]interface{}{
			"kind": "pod",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Clean()
			tt.prepare()
			reader := kubernetes.NewDefaultReader("http://foo", "")
			result := reader.GetPod(tt.namespacedName.namespace, tt.namespacedName.name)
			assert.Equal(t, tt.expect, result)
		})
	}
}

type namespacedName struct {
	namespace string
	name      string
}
