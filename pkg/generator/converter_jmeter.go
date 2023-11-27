/**
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

package generator

import (
	"encoding/xml"
	"log"
	"net/url"

	"github.com/linuxsuren/api-testing/pkg/testing"
)

type jmeterConverter struct {
}

func init() {
	RegisterTestSuiteConverter("jmeter", &jmeterConverter{})
}

func (c *jmeterConverter) Convert(testSuite *testing.TestSuite) (result string, err error) {
	var jmeterTestPlan *JmeterTestPlan
	if jmeterTestPlan, err = c.buildJmeterTestPlan(testSuite); err == nil {
		var data []byte
		if data, err = xml.MarshalIndent(jmeterTestPlan, "  ", "    "); err == nil {
			result = string(data)
		}
	}
	return
}

func (c *jmeterConverter) buildJmeterTestPlan(testSuite *testing.TestSuite) (result *JmeterTestPlan, err error) {
	emptyCtx := make(map[string]interface{})
	if err = testSuite.Render(emptyCtx); err != nil {
		return
	}

	requestItems := []interface{}{}
	for _, item := range testSuite.Items {
		item.Request.RenderAPI(testSuite.API)
		if reqRenderErr := item.Request.Render(emptyCtx, ""); reqRenderErr != nil {
			log.Println("Error rendering request: ", reqRenderErr)
		}

		api, err := url.Parse(item.Request.API)
		if err != nil {
			continue
		}

		requestItem := &HTTPSamplerProxy{
			GUIClass:  "HttpTestSampleGui",
			TestClass: "HTTPSamplerProxy",
			Enabled:   true,
			Name:      item.Name,
			StringProp: []StringProp{{
				Name:  "HTTPSampler.domain",
				Value: api.Hostname(),
			}, {
				Name:  "HTTPSampler.port",
				Value: api.Port(),
			}, {
				Name:  "HTTPSampler.path",
				Value: api.Path,
			}, {
				Name:  "HTTPSampler.method",
				Value: item.Request.Method,
			}},
		}
		if item.Request.Body != "" {
			requestItem.BoolProp = append(requestItem.BoolProp, BoolProp{
				Name:  "HTTPSampler.postBodyRaw",
				Value: "true",
			})
			requestItem.ElementProp = append(requestItem.ElementProp, ElementProp{
				Name: "HTTPsampler.Arguments",
				Type: "Arguments",
				CollectionProp: []CollectionProp{{
					Name: "Arguments.arguments",
					ElementProp: []ElementProp{{
						Name: "",
						Type: "HTTPArgument",
						BoolProp: []BoolProp{{
							Name:  "HTTPArgument.always_encode",
							Value: "false",
						}},
						StringProp: []StringProp{{
							Name:  "Argument.value",
							Value: item.Request.Body,
						}, {
							Name:  "Argument.metadata",
							Value: "=",
						}},
					}},
				}},
			})
		}
		requestItems = append(requestItems, requestItem)
		requestItems = append(requestItems, HashTree{})
	}
	requestItems = append(requestItems, &ResultCollector{
		Enabled:   true,
		GUIClass:  "SummaryReport",
		TestClass: "ResultCollector",
		Name:      "Summary Report",
	})

	result = &JmeterTestPlan{
		Version:    "1.2",
		Properties: "5.0",
		JMeter:     "5.0",
		HashTree: HashTree{
			Items: []interface{}{
				&TestPlan{
					StringProp: []StringProp{{
						Name:  "TestPlan.comments",
						Value: "comment",
					}},
					Name:      testSuite.Name,
					GUIClass:  "TestPlanGui",
					TestClass: "TestPlan",
					Enabled:   true,
				},
				HashTree{
					Items: []interface{}{
						&ThreadGroup{
							StringProp: []StringProp{{
								Name:  "ThreadGroup.num_threads",
								Value: "1",
							}},
							GUIClass:  "ThreadGroupGui",
							TestClass: "ThreadGroup",
							Enabled:   true,
							Name:      "Thread Group",
							ElementProp: ElementProp{
								Name:      "ThreadGroup.main_controller",
								Type:      "LoopController",
								GUIClass:  "LoopControlPanel",
								TestClass: "LoopController",
								BoolProp: []BoolProp{{
									Name:  "LoopController.continue_forever",
									Value: "false",
								}},
								StringProp: []StringProp{{
									Name:  "LoopController.loops",
									Value: "1",
								}},
							},
						},
						HashTree{
							Items: requestItems,
						},
					},
				},
			},
		},
	}
	return
}

type JmeterTestPlan struct {
	XMLName    xml.Name `xml:"jmeterTestPlan"`
	Version    string   `xml:"version,attr"`
	Properties string   `xml:"properties,attr"`
	JMeter     string   `xml:"jmeter,attr"`
	HashTree   HashTree `xml:"hashTree"`
}

type HashTree struct {
	XMLName xml.Name      `xml:"hashTree"`
	Items   []interface{} `xml:"items"`
}

type TestPlan struct {
	XMLName    xml.Name     `xml:"TestPlan"`
	Name       string       `xml:"testname,attr"`
	GUIClass   string       `xml:"guiclass,attr"`
	TestClass  string       `xml:"testclass,attr"`
	Enabled    bool         `xml:"enabled,attr"`
	StringProp []StringProp `xml:"stringProp"`
}

type ThreadGroup struct {
	XMLName     xml.Name     `xml:"ThreadGroup"`
	GUIClass    string       `xml:"guiclass,attr"`
	TestClass   string       `xml:"testclass,attr"`
	Enabled     bool         `xml:"enabled,attr"`
	Name        string       `xml:"testname,attr"`
	StringProp  []StringProp `xml:"stringProp"`
	ElementProp ElementProp  `xml:"elementProp"`
}

type HTTPSamplerProxy struct {
	XMLName     xml.Name      `xml:"HTTPSamplerProxy"`
	Name        string        `xml:"testname,attr"`
	GUIClass    string        `xml:"guiclass,attr"`
	TestClass   string        `xml:"testclass,attr"`
	Enabled     bool          `xml:"enabled,attr"`
	StringProp  []StringProp  `xml:"stringProp"`
	BoolProp    []BoolProp    `xml:"boolProp"`
	ElementProp []ElementProp `xml:"elementProp"`
}

type ResultCollector struct {
	XMLName   xml.Name `xml:"ResultCollector"`
	Enabled   bool     `xml:"enabled,attr"`
	GUIClass  string   `xml:"guiclass,attr"`
	TestClass string   `xml:"testclass,attr"`
	Name      string   `xml:"testname,attr"`
}

type ElementProp struct {
	Name           string           `xml:"name,attr"`
	Type           string           `xml:"elementType,attr"`
	GUIClass       string           `xml:"guiclass,attr"`
	TestClass      string           `xml:"testclass,attr"`
	Enabled        bool             `xml:"enabled,attr"`
	StringProp     []StringProp     `xml:"stringProp"`
	BoolProp       []BoolProp       `xml:"boolProp"`
	CollectionProp []CollectionProp `xml:"collectionProp"`
}

type CollectionProp struct {
	Name        string        `xml:"name,attr"`
	ElementProp []ElementProp `xml:"elementProp"`
}

type StringProp struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type BoolProp struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}
