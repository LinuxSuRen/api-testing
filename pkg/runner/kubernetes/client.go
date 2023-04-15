package kubernetes

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/antonmedv/expr"
)

// Reader represents a reader interface
type Reader interface {
	GetPod(namespace, name string) map[string]interface{}
	GetDeploy(namespace, name string) map[string]interface{}
	GetResource(group, kind, version, namespace, name string) map[string]interface{}
}

type defualtReader struct {
	server string
	token  string
}

// NewDefaultReader returns a reader implement
func NewDefaultReader(server, token string) Reader {
	return &defualtReader{
		server: server,
		token:  token,
	}
}

// GetPod gets a pod by namespace and name
func (r *defualtReader) GetPod(namespace, name string) (result map[string]interface{}) {
	api := fmt.Sprintf("%s/api/v1/namespaces/%s/pods/%s", r.server, namespace, name)
	return r.request(api)
}

// GetDeploy gets a pod by namespace and name
func (r *defualtReader) GetDeploy(namespace, name string) (result map[string]interface{}) {
	api := fmt.Sprintf("%s/api/v1/namespaces/%s/deployments/%s", r.server, namespace, name)
	return r.request(api)
}

func (r *defualtReader) GetResource(group, kind, version, namespace, name string) (result map[string]interface{}) {
	api := fmt.Sprintf("%s/api/%s/%s/namespaces/%s/%s/%s", r.server, group, version, namespace, kind, name)
	api = strings.ReplaceAll(api, "api//", "api/")
	if !strings.Contains(api, "api/v1") {
		api = strings.ReplaceAll(api, "api/", "apis/")
	}
	return r.request(api)
}

func (r *defualtReader) request(api string) (result map[string]interface{}) {
	client := GetClient()

	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.token))
	if resp, err := client.Do(req); err == nil && resp.StatusCode == http.StatusOK {
		if data, err := io.ReadAll(resp.Body); err == nil {
			result = make(map[string]interface{})

			_ = json.Unmarshal(data, &result)
		}
	}
	return
}

var client *http.Client

// GetClient returns a default client
func GetClient() *http.Client {
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}
	return client
}

type ResourceValidator interface {
	Exist() bool
}

type defaultResourceValidator struct {
	data map[string]interface{}
}

func (v *defaultResourceValidator) Exist() bool {
	return v.data != nil && len(v.data) > 0
}

func podValidator(params ...interface{}) (validator interface{}, err error) {
	return resourceValidator(append([]interface{}{"pods"}, params...)...)
}

func resourceValidator(params ...interface{}) (validator interface{}, err error) {
	if len(params) < 3 {
		err = errors.New("there are three params at least")
		return
	}

	var kind string
	version := "v1"
	group := ""
	switch obj := params[0].(type) {
	case string:
		kind = obj
	case map[string]interface{}:
		if obj["kind"] != nil {
			kind = obj["kind"].(string)
		}
		if obj["version"] != nil {
			version = obj["version"].(string)
		}
		if obj["group"] != nil {
			group = obj["group"].(string)
		}
	}

	if kind == "" {
		err = errors.New("kind is required")
		return
	}

	switch kind {
	case "deployments", "statefulsets", "daemonsets":
		group = "apps"
	}

	server := os.Getenv("KUBERNETES_SERVER")
	token := os.Getenv("KUBERNETES_TOKEN")
	if server == "" || token == "" {
		err = errors.New("KUBERNETES_SERVER and KUBERNETES_TOKEN are required")
		return
	}
	reader := NewDefaultReader(server, token)
	validator = &defaultResourceValidator{
		data: reader.GetResource(group, kind, version, params[1].(string), params[2].(string)),
	}
	return
}

// PodValidatorFunc returns a expr for checking pod existing
func PodValidatorFunc() expr.Option {
	return expr.Function("pod", podValidator, new(func(...string) ResourceValidator))
}

func KubernetesValidatorFunc() expr.Option {
	return expr.Function("k8s", resourceValidator, new(func(interface{}, ...string) ResourceValidator))
}
