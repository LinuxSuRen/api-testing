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
package mock

import (
	"errors"
	"github.com/linuxsuren/api-testing/docs"
	"os"

	"gopkg.in/yaml.v3"
)

type Reader interface {
	Parse() (*Server, error)
	GetData() []byte
}

type Writer interface {
	Write([]byte)
}

type localFileReader struct {
	file string
	data []byte
}

func NewLocalFileReader(file string) Reader {
	return &localFileReader{file: file}
}

func (r *localFileReader) Parse() (server *Server, err error) {
	if r.data, err = os.ReadFile(r.file); err == nil {
		server, err = validateAndParse(r.data)
	}
	return
}

func (r *localFileReader) GetData() []byte {
	return r.data
}

type inMemoryReader struct {
	data []byte
}

type ReaderAndWriter interface {
	Reader
	Writer
}

func NewInMemoryReader(config string) ReaderAndWriter {
	return &inMemoryReader{
		data: []byte(config),
	}
}

func (r *inMemoryReader) Parse() (server *Server, err error) {
	server, err = validateAndParse(r.data)
	return
}

func (r *inMemoryReader) GetData() []byte {
	return r.data
}

func (r *inMemoryReader) Write(data []byte) {
	r.data = data
}

func validateAndParse(data []byte) (server *Server, err error) {
	server = &Server{}
	err = yaml.Unmarshal(data, server)
	err = errors.Join(err, docs.Validate(data, docs.MockSchema))
	return
}
