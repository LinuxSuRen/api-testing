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
		group          string
		version        string
		kind           string
		namespacedName namespacedName
		prepare        func()
		expect         map[string]interface{}
	}{{
		name:    "normal",
		kind:    "pods",
		version: "v1",
		namespacedName: namespacedName{
			namespace: "ns",
			name:      "fake",
		},
		prepare: func() {
			gock.New(urlFoo).
				Get("/api/v1/namespaces/ns/pods/fake").
				Reply(http.StatusOK).
				JSON(`{"kind":"pod"}`)
			gock.InterceptClient(kubernetes.GetClient())
		},
		expect: map[string]interface{}{
			"kind": "pod",
		},
	}, {
		name:    "deployments",
		kind:    "deployments",
		version: "v1",
		group:   "apps",
		namespacedName: namespacedName{
			namespace: "ns",
			name:      "fake",
		},
		prepare: func() {
			gock.New(urlFoo).
				Get("/apis/apps/v1/namespaces/ns/deployments/fake").
				Reply(http.StatusOK).
				JSON(`{"kind":"deployment"}`)
			gock.InterceptClient(kubernetes.GetClient())
		},
		expect: map[string]interface{}{
			"kind": "deployment",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Clean()
			tt.prepare()
			reader := kubernetes.NewDefaultReader(urlFoo, "")
			result, err := reader.GetResource(tt.group, tt.kind, tt.version, tt.namespacedName.namespace, tt.namespacedName.name)
			assert.Equal(t, tt.expect, result)
			assert.Nil(t, err)
		})
	}
}

type namespacedName struct {
	namespace string
	name      string
}
