package kubernetes_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/antonmedv/expr"
	"github.com/h2non/gock"
	"github.com/linuxsuren/api-testing/pkg/runner/kubernetes"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestKubernetesValidatorFunc(t *testing.T) {
	os.Setenv("KUBERNETES_SERVER", urlFoo)
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
		name:       "daemonset count",
		expression: `k8s("daemonsets", "ns", "foo").ExpectCount(0)`,
		prepare:    prepareDaemonset,
		expectBool: true,
	}, {
		name:       "no kind",
		expression: `k8s({"foo": "bar"}, "ns", "foo").Exist()`,
		expectErr:  true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare = util.MakeSureNotNil(tt.prepare)
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

func preparePod() {
	gock.New(urlFoo).
		Get("/api/v1/namespaces/ns/pods/foo").
		MatchHeader("Authorization", defaultToken).
		Reply(http.StatusOK).
		JSON(`{"kind":"pod"}`)
}

func prepareDeploy() {
	gock.New(urlFoo).
		Get("/apis/apps/v1/namespaces/ns/deployments/foo").
		MatchHeader("Authorization", defaultToken).
		Reply(http.StatusOK).
		JSON(`{"kind":"deploy"}`)
}

func prepareStatefulset() {
	gock.New(urlFoo).
		Get("/apis/apps/v1/namespaces/ns/statefulsets/foo").
		MatchHeader("Authorization", defaultToken).
		Reply(http.StatusOK).
		JSON(`{"kind":"statefulset"}`)
}

func prepareDaemonset() {
	gock.New(urlFoo).
		Get("/apis/apps/v1/namespaces/ns/daemonsets/foo").
		MatchHeader("Authorization", defaultToken).
		Reply(http.StatusOK).
		JSON(`{"kind":"daemonset","items":[]}`)
}

func prepareCRDVM() {
	gock.New(urlFoo).
		Get("/apis/bar/v2/namespaces/ns/vms/foo").
		MatchHeader("Authorization", defaultToken).
		Reply(http.StatusOK).
		JSON(`{"kind":"vm"}`)
}

const urlFoo = "http://foo"
const defaultToken = "Bearer token"
