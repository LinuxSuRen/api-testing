/*
Copyright 2025 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language 24 permissions and
limitations under the License.
*/
package render

import (
	"errors"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

type UserDefinedTemplates struct {
	Items []UserDefinedTemplate
}

type UserDefinedTemplate struct {
	Name   string
	Render string
}

func (t *UserDefinedTemplates) Validate() (err error) {
	for _, item := range t.Items {
		_, rErr := Render(item.Name, item.Render, nil)
		err = errors.Join(err, rErr)
	}
	return
}

func (t *UserDefinedTemplates) ConflictWith(funcMap template.FuncMap) (conflict error) {
	for _, item := range t.Items {
		if _, ok := funcMap[item.Name]; ok {
			conflict = errors.New("conflict with existing function: " + item.Name)
			break
		}
	}
	return
}

func ParseUserDefinedTemplates(data []byte) (templates *UserDefinedTemplates, err error) {
	templates = &UserDefinedTemplates{}
	err = yaml.Unmarshal(data, templates)
	return
}

func ParseUserDefinedTemplatesFromFile(filePath string) (templates *UserDefinedTemplates, err error) {
	var data []byte
	data, err = os.ReadFile(filePath)
	if err == nil {
		templates, err = ParseUserDefinedTemplates(data)
		if err == nil {
			err = templates.Validate()
		}
	}
	return
}
