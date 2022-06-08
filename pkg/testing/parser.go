package testing

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Parse(configFile string) (testCase *TestCase, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(configFile); err != nil {
		return
	}

	testCase = &TestCase{}
	if err = yaml.Unmarshal(data, testCase); err != nil {
		return
	}
	return
}
