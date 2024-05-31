package downloader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OCIDownloader interface {
	WithBasicAuth(username string, password string)
	Download(image, tag, file string) (reader io.Reader, err error)
}

type defaultOCIDownloader struct {
}

func NewDefaultOCIDownloader() OCIDownloader {
	return &defaultOCIDownloader{}
}

func (d *defaultOCIDownloader) WithBasicAuth(username string, password string) {
	fmt.Println("not support yet")
}

func (d *defaultOCIDownloader) Download(image, tag, file string) (reader io.Reader, err error) {
	authStr := auth(image)

	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("https://index.docker.io/v2/%s/manifests/%s", image, tag), nil); err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authStr))
	req.Header.Set("Accept", "application/vnd.oci.image.manifest.v1+json")

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		panic(err)
	} else {
		var data []byte
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		manifest := &Manifest{}
		if err := json.Unmarshal(data, manifest); err != nil {
			panic(err)
		}

		for _, layer := range manifest.Layers {
			if v, ok := layer.Annotations["org.opencontainers.image.title"]; ok && v == file {
				reader = downloadLayer(image, layer.Digest, authStr)
				return
			}
		}
	}

	err = fmt.Errorf("not found %s", file)
	return
}

func downloadLayer(image, digest, authToken string) io.Reader {
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("https://registry-1.docker.io/v2/%s/blobs/%s", image, digest), nil); err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	if resp, err := http.DefaultClient.Do(req); err != nil {
		panic(err)
	} else {
		var data []byte
		if data, err = io.ReadAll(resp.Body); err != nil {
			panic(err)
		}

		defer resp.Body.Close()
		return bytes.NewBuffer(data)
	}
}

func auth(image string) string {
	resp, err := http.Get(fmt.Sprintf("https://auth.docker.io/token?scope=repository:%s:pull&service=registry.docker.io", image))
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	regAuth := &RegistryAuth{}
	if err := json.Unmarshal(data, regAuth); err != nil {
		panic(err)
	}

	return regAuth.Token
}

type RegistryAuth struct {
	Token string `json:"token"`
}

type Manifest struct {
	Layers []Layer `json:"layers"`
}

type Layer struct {
	Annotations map[string]string `json:"annotations"`
	Digest      string            `json:"digest"`
}
