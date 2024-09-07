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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/blang/semver/v4"
)

type OCIDownloader interface {
	WithBasicAuth(username string, password string)
	WithRegistry(string)
	WithRoundTripper(http.RoundTripper)
	WithInsecure(bool)
	WithTimeout(time.Duration)
	WithContext(context.Context)
	Download(image, tag, file string) (reader io.Reader, err error)
}

type PlatformAwareOCIDownloader interface {
	OCIDownloader
	WithOS(string)
	WithArch(string)
	GetTargetFile() string
}

type defaultOCIDownloader struct {
	ctx          context.Context
	timeout      time.Duration
	serviceURL   string
	registry     string
	rawImage     string
	protocol     string
	roundTripper http.RoundTripper
}

func NewDefaultOCIDownloader() OCIDownloader {
	return &defaultOCIDownloader{
		protocol: "https",
		timeout:  time.Minute,
		ctx:      context.Background(),
	}
}

func (d *defaultOCIDownloader) WithBasicAuth(username string, password string) {
	fmt.Println("not support yet")
}

func (d *defaultOCIDownloader) Download(image, tag, file string) (reader io.Reader, err error) {
	if d.registry == "" {
		d.registry = getRegistry(image)
	}
	d.rawImage = strings.TrimPrefix(image, d.registry)
	var authStr string
	if authStr, err = d.auth(d.rawImage); err != nil {
		err = fmt.Errorf("failed to get auth token: %w", err)
		return
	}

	latestTag := tag
	if tag == "" {
		latestTag, err = d.getLatestTag(d.rawImage, authStr)
		if err != nil {
			err = fmt.Errorf("failed to get latest tag: %w", err)
			return
		}
	}

	var req *http.Request
	api := fmt.Sprintf("%s://%s/v2/%s/manifests/%s", d.protocol, d.registry, d.rawImage, latestTag)
	if req, err = http.NewRequestWithContext(d.ctx, http.MethodGet, api, nil); err != nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authStr))
	req.Header.Set("Accept", "application/vnd.oci.image.manifest.v1+json")

	var resp *http.Response
	if resp, err = d.getHTTPClient().Do(req); err != nil {
		err = fmt.Errorf("failed to get manifest from %q, error: %v", api, err)
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to get manifest from %q, status code: %d", api, resp.StatusCode)
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
				reader, err = d.downloadLayer(d.rawImage, layer.Digest, authStr)
				return
			}
		}
	}

	err = fmt.Errorf("not found %s", file)
	return
}

func (d *defaultOCIDownloader) WithRegistry(registry string) {
	d.registry = registry
}

func (d *defaultOCIDownloader) WithInsecure(insecure bool) {
	if insecure {
		d.protocol = "http"
	} else {
		d.protocol = "https"
	}
}

func (d *defaultOCIDownloader) WithTimeout(timeout time.Duration) {
	d.timeout = timeout
}

func (d *defaultOCIDownloader) WithContext(ctx context.Context) {
	d.ctx = ctx
}

func (d *defaultOCIDownloader) WithRoundTripper(rt http.RoundTripper) {
	d.roundTripper = rt
}

// getLatestTag returns the latest artifact tag
// we assume the artifact tags do not have the prefix `v`
func (d *defaultOCIDownloader) getLatestTag(image, authToken string) (tag string, err error) {
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s://%s/v2/%s/tags/list", d.protocol, d.registry, image), nil); err != nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	var resp *http.Response
	if resp, err = d.getHTTPClient().Do(req); err != nil {
		err = fmt.Errorf("failed to get image tags from %q, error: %v", req.URL, err)
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to get image tags from %q, status code: %d", req.URL, resp.StatusCode)
	} else {
		defer resp.Body.Close()

		progressReader := NewProgressReader(resp.Body)
		progressReader.SetLength(resp.ContentLength)

		var data []byte
		if data, err = io.ReadAll(progressReader); err != nil {
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

func (d *defaultOCIDownloader) getHTTPClient() (client *http.Client) {
	client = &http.Client{
		Timeout:   d.timeout,
		Transport: d.roundTripper,
	}
	return
}

type ImageTagList struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func (d *defaultOCIDownloader) downloadLayer(image, digest, authToken string) (reader io.Reader, err error) {
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s://%s/v2/%s/blobs/%s", d.protocol, d.registry, image, digest), nil); err != nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	var resp *http.Response
	if resp, err = d.getHTTPClient().Do(req); err != nil {
		err = fmt.Errorf("failed to get layer from %q, error: %v", req.URL.String(), err)
	} else {
		var data []byte
		if data, err = io.ReadAll(resp.Body); err != nil {
			return
		}

		defer resp.Body.Close()
		reader = bytes.NewBuffer(data)
	}
	return
}

// getRegistry returns the registry of the image
// e.g. registry-1.docker.io, ghcr.io, quay.io
func getRegistry(image string) string {
	segs := strings.Split(image, "/")
	if len(segs) == 1 || len(segs) == 2 || segs[0] == "docker.io" {
		return DockerHubRegistry
	}
	return segs[0]
}

const DockerHubRegistry = "registry-1.docker.io"

func detectAuthURL(protocol, image string) (authURL string, service string, err error) {
	registry := getRegistry(image)
	detectURL := fmt.Sprintf("%s://%s/v2/", protocol, registry)

	var resp *http.Response
	resp, err = http.Get(detectURL)
	if err != nil {
		return
	}

	authHeader := resp.Header.Get(HeaderWWWAuthenticate)
	authHeader = strings.ReplaceAll(authHeader, "Bearer ", "")
	for _, v := range strings.Split(authHeader, ",") {
		if strings.HasPrefix(v, "realm=") {
			v = strings.ReplaceAll(v, "realm=", "")
			v = strings.TrimPrefix(v, `"`)
			v = strings.TrimSuffix(v, `"`)
			authURL = v
		} else if strings.HasPrefix(v, "service=") {
			v = strings.ReplaceAll(v, "service=", "")
			v = strings.TrimPrefix(v, `"`)
			v = strings.TrimSuffix(v, `"`)
			service = v
		}
	}
	return
}

const (
	HeaderWWWAuthenticate = "www-authenticate"
)

func (d *defaultOCIDownloader) auth(image string) (authToken string, err error) {
	var authURL string
	if authURL, d.serviceURL, err = detectAuthURL(d.protocol, fmt.Sprintf("%s/%s", d.registry, d.rawImage)); err != nil {
		return
	}

	var resp *http.Response
	resp, err = http.Get(fmt.Sprintf("%s?scope=repository:%s:pull&service=%s", authURL, image, d.serviceURL))
	if err != nil {
		err = fmt.Errorf("failed to get auth token from %q, error: %v", authURL, err)
		return
	}

	var data []byte
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	regAuth := &RegistryAuth{}
	if err = json.Unmarshal(data, regAuth); err != nil {
		return
	}

	authToken = regAuth.Token
	return
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
