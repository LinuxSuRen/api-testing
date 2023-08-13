package runner

import "github.com/linuxsuren/api-testing/pkg/testing"

// GetTestSuiteRunner returns a proper runner according to the given test suite.
func GetTestSuiteRunner(suite *testing.TestSuite) TestCaseRunner {
	// TODO: should be refactored to meet more types of runners
	if suite.Spec.GRPC != nil {
		return NewGRPCTestCaseRunner(suite.API, *suite.Spec.GRPC)
	}
	return NewSimpleTestCaseRunner()
}
