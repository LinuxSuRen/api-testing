/*
Copyright 2024 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package downloader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/blang/semver/v4"
)

type OCIDownloader interface {
	WithBasicAuth(username string, password string)
	Download(image, tag, file string) (reader io.Reader, err error)
}

type PlatformAwareOCIDownloader interface {
	OCIDownloader
	WithOS(string)
	WithArch(string)
	GetTargetFile() string
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
	latestTag := tag
	if tag == "" {
		latestTag, err = getLatestTag(image, authStr)
		if err != nil {
			return
		}
	}

	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("https://index.docker.io/v2/%s/manifests/%s", image, latestTag), nil); err != nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authStr))
	req.Header.Set("Accept", "application/vnd.oci.image.manifest.v1+json")

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return
	} else {
		var data []byte
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		manifest := &Manifest{}
		if err = json.Unmarshal(data, manifest); err != nil {
			return
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

// getLatestTag returns the latest artifact tag
// we assume the artifact tags do not have the prefix `v`
func getLatestTag(image, authToken string) (tag string, err error) {
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("https://registry-1.docker.io/v2/%s/tags/list", image), nil); err != nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil || resp.StatusCode != http.StatusOK {
		return
	} else {
		defer resp.Body.Close()

		var data []byte
		if data, err = io.ReadAll(resp.Body); err != nil {
			return
		}

		imageTagList := &ImageTagList{}
		if err = json.NewDecoder(bytes.NewBuffer(data)).Decode(imageTagList); err == nil {
			var latestVer semver.Version
			for _, t := range imageTagList.Tags {
				if strings.HasPrefix(t, "v") {
					continue
				}

				if ver, verErr := semver.ParseTolerant(t); verErr == nil {
					if ver.GT(latestVer) {
						tag = t
						latestVer = ver
					}
				}
			}
		}
	}
	return
}

type ImageTagList struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
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
