package kubernetes_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/antonmedv/expr"
	"github.com/h2non/gock"
	"github.com/linuxsuren/api-testing/pkg/runner/kubernetes"
	"github.com/stretchr/testify/assert"
)

func TestKubernetesValidatorFunc(t *testing.T) {
	os.Setenv("KUBERNETES_SERVER", "http://foo")
	os.Setenv("KUBERNETES_TOKEN", "token")
	gock.InterceptClient(kubernetes.GetClient())
	defer gock.RestoreClient(http.DefaultClient)
	defer gock.Off()

	tests := []struct {
		name       string
		prepare    func()
		expression string
		expectBool bool
		expectErr  bool
	}{{
		name:       "pod exist expr",
		prepare:    preparePod,
		expression: `pod('ns', 'foo').Exist()`,
		expectBool: true,
	}, {
		name:       "pod expectField expr",
		prepare:    preparePod,
		expression: `pod('ns', 'foo').ExpectField('pod', 'kind')`,
		expectBool: true,
	}, {
		name:       "pod expectField expr, not match",
		prepare:    preparePod,
		expression: `pod('ns', 'foo').ExpectField('pods', 'kind')`,
		expectBool: false,
	}, {
		name:       "pod expectField expr, not find field",
		prepare:    preparePod,
		expression: `pod('ns', 'foo').ExpectField('pods', 'kinds')`,
		expectBool: false,
	}, {
		name:       "no enough params",
		expression: `k8s('crd')`,
		prepare:    emptyPrepare,
		expectBool: false,
		expectErr:  true,
	}, {
		name:       "crd",
		expression: `k8s({"kind":"vms","group":"bar","version":"v2"}, "ns", "foo").Exist()`,
		prepare:    prepareCRDVM,
		expectBool: true,
	}, {
		name:       "deploy",
		expression: `k8s("deployments", "ns", "foo").Exist()`,
		prepare:    prepareDeploy,
		expectBool: true,
	}, {
		name:       "statefulset",
		expression: `k8s("statefulsets", "ns", "foo").Exist()`,
		prepare:    prepareStatefulset,
		expectBool: true,
	}, {
		name:       "daemonset",
		expression: `k8s("daemonsets", "ns", "foo").Exist()`,
		prepare:    prepareDaemonset,
		expectBool: true,
	}, {
		name:       "no kind",
		expression: `k8s({"foo": "bar"}, "ns", "foo").Exist()`,
		prepare:    emptyPrepare,
		expectErr:  true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			vm, err := expr.Compile(tt.expression, kubernetes.KubernetesValidatorFunc(),
				kubernetes.PodValidatorFunc())
			assert.Nil(t, err)

			result, err := expr.Run(vm, expr.Env(tt))
			assert.Equal(t, tt.expectErr, err != nil)
			if err == nil {
				assert.Equal(t, tt.expectBool, result)
			}
		})
	}
}

func emptyPrepare() {}

func preparePod() {
	gock.New("http://foo").
		Get("/api/v1/namespaces/ns/pods/foo").
		MatchHeader("Authorization", "Bearer token").
		Reply(http.StatusOK).
		JSON(`{"kind":"pod"}`)
}

func prepareDeploy() {
	gock.New("http://foo").
		Get("/apis/apps/v1/namespaces/ns/deployments/foo").
		MatchHeader("Authorization", "Bearer token").
		Reply(http.StatusOK).
		JSON(`{"kind":"deploy"}`)
}

func prepareStatefulset() {
	gock.New("http://foo").
		Get("/apis/apps/v1/namespaces/ns/statefulsets/foo").
		MatchHeader("Authorization", "Bearer token").
		Reply(http.StatusOK).
		JSON(`{"kind":"statefulset"}`)
}

func prepareDaemonset() {
	gock.New("http://foo").
		Get("/apis/apps/v1/namespaces/ns/daemonsets/foo").
		MatchHeader("Authorization", "Bearer token").
		Reply(http.StatusOK).
		JSON(`{"kind":"daemonset"}`)
}

func prepareCRDVM() {
	gock.New("http://foo").
		Get("/apis/bar/v2/namespaces/ns/vms/foo").
		MatchHeader("Authorization", "Bearer token").
		Reply(http.StatusOK).
		JSON(`{"kind":"vm"}`)
}
