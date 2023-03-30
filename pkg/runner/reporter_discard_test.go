package runner_test

import (
	"testing"

	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestDiscardTestReporter(t *testing.T) {
	reporter := runner.NewDiscardTestReporter()
	assert.NotNil(t, reporter)
	assert.Nil(t, reporter.GetAllRecords())

	result, err := reporter.ExportAllReportResults()
	assert.Nil(t, result)
	assert.Nil(t, err)

	reporter.PutRecord(&runner.ReportRecord{})
}
