package kubernetes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/logging"
	"github.com/linuxsuren/api-testing/pkg/util"

	unstructured "github.com/linuxsuren/unstructured/pkg"
)

var (
	writeLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("write")
)

// Reader represents a reader interface
type Reader interface {
	GetResource(group, kind, version, namespace, name string) (map[string]interface{}, error)
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

// GetResource returns the resource
func (r *defualtReader) GetResource(group, kind, version, namespace, name string) (map[string]interface{}, error) {
	api := fmt.Sprintf("%s/api/%s/%s/namespaces/%s/%s/%s", r.server, group, version, namespace, kind, name)
	api = strings.ReplaceAll(api, "api//", "api/")
	if !strings.Contains(api, "api/v1") {
		api = strings.ReplaceAll(api, "api/", "apis/")
	}
	return r.request(api)
}

func (r *defualtReader) request(api string) (result map[string]interface{}, err error) {
	client := GetClient()
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, api, nil); err == nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.token))
		var resp *http.Response
		if resp, err = client.Do(req); err == nil && resp.StatusCode == http.StatusOK {
			var data []byte
			if data, err = io.ReadAll(resp.Body); err == nil {
				result = make(map[string]interface{})

				err = json.Unmarshal(data, &result)
			}
		}
	}
	return
}

var client *http.Client

// GetClient returns a default client
func GetClient() *http.Client {
	if client == nil {
		client = util.TlsAwareHTTPClient(true) // TODO should have a way to change it
	}
	return client
}

// ResourceValidator represents a generic resource validator
type ResourceValidator interface {
	Exist() bool
	Count() int
	ExpectCount(int) bool
	ExpectField(value interface{}, fields ...string) bool
}

type defaultResourceValidator struct {
	data map[string]interface{}
	err  error
}

// Exist returns true if the specific resource exist
func (v *defaultResourceValidator) Exist() bool {
	if v.err != nil {
		writeLogger.Info(v.err.Error())
		return false
	}
	return v.data != nil && len(v.data) > 0
}

// Count returns the count of the specific resources
func (v *defaultResourceValidator) Count() int {
	val, ok, err := unstructured.NestedField(v.data, "items")
	if ok && err == nil {
		items := val.([]interface{})
		return len(items)
	}
	return -1
}

// ExpectCount returns true if there are specific number of the target k8s resources
func (v *defaultResourceValidator) ExpectCount(expect int) bool {
	return expect == v.Count()
}

// ExpectField verify the specific field and value
func (v *defaultResourceValidator) ExpectField(value interface{}, fields ...string) (result bool) {
	val, ok, err := unstructured.NestedField(v.data, fields...)
	if !ok || err != nil {
		fmt.Printf("cannot find '%v',error: %v\n", fields, err)
		return
	}
	if result = fmt.Sprintf("%v", val) == fmt.Sprintf("%v", value); !result {
		fmt.Printf("expect: '%v', actual: '%v'\n", value, val)
	}
	return
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
	data, err := reader.GetResource(group, kind, version, params[1].(string), params[2].(string))
	validator = &defaultResourceValidator{
		data: data,
		err:  err,
	}
	return
}
