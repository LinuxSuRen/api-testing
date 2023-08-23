/*
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
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
