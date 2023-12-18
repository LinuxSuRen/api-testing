/*
Copyright 2023 API Testing Authors.

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

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/linuxsuren/api-testing/pkg/generator"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/spf13/cobra"
)

func createConvertCommand() (c *cobra.Command) {
	opt := &convertOption{}
	c = &cobra.Command{
		Use:     "convert",
		Short:   "Convert the API testing file to other format",
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
	}

	converters := generator.GetTestSuiteConverters()

	flags := c.Flags()
	flags.StringVarP(&opt.pattern, "pattern", "p", "test-suite-*.yaml",
		"The file pattern which try to execute the test cases. Brace expansion is supported, such as: test-suite-{1,2}.yaml")
	flags.StringVarP(&opt.converter, "converter", "", "",
		fmt.Sprintf("The converter format, supported: %s", util.Keys(converters)))
	flags.StringVarP(&opt.source, "source", "", "", "The source format, supported: postman")
	flags.StringVarP(&opt.target, "target", "t", "", "The target file path")

	_ = c.MarkFlagRequired("pattern")
	_ = c.MarkFlagRequired("converter")
	return
}

type convertOption struct {
	pattern   string
	converter string
	source    string
	target    string
}

func (o *convertOption) preRunE(c *cobra.Command, args []string) (err error) {
	switch o.source {
	case "postman":
		o.target = util.EmptyThenDefault(o.target, "sample.yaml")
		o.converter = "raw"
	case "":
		o.target = util.EmptyThenDefault(o.target, "sample.jmx")
	default:
		err = errors.New("only postman supported")
	}

	return
}

func (o *convertOption) runE(c *cobra.Command, args []string) (err error) {
	loader := testing.NewFileWriter("")
	if err = loader.Put(o.pattern); err != nil {
		return
	}

	var suite *testing.TestSuite
	if o.source == "" {
		suite, err = getSuiteFromFile(o.pattern)
	} else {
		suite, err = generator.NewPostmanImporter().ConvertFromFile(o.pattern)
	}

	if err != nil {
		return
	}

	converter := generator.GetTestSuiteConverter(o.converter)
	if converter == nil {
		err = fmt.Errorf("no converter found")
	} else {
		var output string
		output, err = converter.Convert(suite)
		if output != "" {
			err = os.WriteFile(o.target, []byte(output), 0644)
		}
	}
	return
}

func getSuiteFromFile(pattern string) (suite *testing.TestSuite, err error) {
	loader := testing.NewFileWriter("")
	if err = loader.Put(pattern); err == nil {
		var suites []testing.TestSuite
		if suites, err = loader.ListTestSuite(); err == nil {
			if len(suites) > 0 {
				suite = &suites[0]
			} else {
				err = errors.New("no suites found")
			}
		}
	}
	return
}
