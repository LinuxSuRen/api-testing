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

package runner

import (
	"testing"

	testWriter "github.com/linuxsuren/api-testing/pkg/runner/writer_templates"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func TestGRPCResultWriter(t *testing.T) {
	t.Run("test request", func(t *testing.T) {
		s := grpc.NewServer()
		testServer := &testWriter.ReportServer{}
		testWriter.RegisterReportWriterServer(s, testServer)
		reflection.RegisterV1(s)
		l := runServer(t, s)
		api := "/writer_templates.ReportWriter/SendReportResult"
		host := l.Addr().String()
		url := host + api
		writer := NewGRPCResultWriter(url)
		err := writer.Output([]ReportResult{{
			Name:    "test",
			API:     "/api",
			Max:     1,
			Average: 2,
			Error:   3,
			Count:   1,
		}})
		assert.NoError(t, err)
		s.Stop()
	})
	t.Run("test reflect unsupported on server", func(t *testing.T) {
		s := grpc.NewServer()
		testServer := &testWriter.ReportServer{}
		testWriter.RegisterReportWriterServer(s, testServer)
		l := runServer(t, s)
		api := "/writer_templates.ReportWriter/SendReportResult"
		host := l.Addr().String()
		url := host + api
		writer := NewGRPCResultWriter(url)
		err := writer.Output([]ReportResult{{
			Name:    "test",
			API:     "/api",
			Max:     1,
			Average: 2,
			Error:   3,
			Count:   1,
		}})
		assert.NotNil(t, err)
	})
}
