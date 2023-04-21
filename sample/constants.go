package sample

import _ "embed"

//go:embed testsuite-gitlab.yaml
var TestSuiteGitLab string

//go:embed api-testing-schema.json
var Schema string
