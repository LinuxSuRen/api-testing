package kubernetes

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/antonmedv/expr"
)

// Reader represents a reader interface
type Reader interface {
	GetPod(namespace, name string) map[string]interface{}
}

type defualtReader struct {
	server string
	token  string
}

// NewDefualtReader returns a reader implement
func NewDefualtReader(server, token string) Reader {
	return &defualtReader{
		server: server,
		token:  token,
	}
}

// GetPod gets a pod by namespace and name
func (r *defualtReader) GetPod(namespace, name string) (result map[string]interface{}) {
	client := GetClient()

	api := fmt.Sprintf("%s/api/v1/namespaces/%s/pods/%s", r.server, namespace, name)
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

func podExist(params ...interface{}) (exist interface{}, err error) {
	server := os.Getenv("KUBERNETES_SERVER")
	token := os.Getenv("KUBERNETES_TOKEN")
	reader := NewDefualtReader(server, token)
	exist = (reader.GetPod(params[0].(string), params[1].(string)) != nil)
	return
}

// PodExistFunc returns a expr for checking pod existing
func PodExistFunc() expr.Option {
	return expr.Function("podExist", podExist, new(func(...string) bool))
}
