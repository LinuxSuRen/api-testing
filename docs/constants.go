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
package docs

import (
	_ "embed"
	"fmt"

	yamlconv "github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed api-testing-schema.json
var Schema string

//go:embed api-testing-mock-schema.json
var MockSchema string

func Validate(data []byte, schema string) (err error) {
	// convert YAML to JSON
	var jsonData []byte
	if jsonData, err = yamlconv.YAMLToJSON(data); err == nil {
		schemaLoader := gojsonschema.NewStringLoader(schema)
		documentLoader := gojsonschema.NewBytesLoader(jsonData)

		var result *gojsonschema.Result
		if result, err = gojsonschema.Validate(schemaLoader, documentLoader); err == nil {
			if !result.Valid() {
				err = fmt.Errorf("%v", result.Errors())
			}
		}
	}
	return
}
