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

package runner

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type prometheusReporter struct {
	remote       string
	sync         bool
	errorCount   prometheus.Counter
	successCount prometheus.Counter
}

// NewPrometheusWriter creates a new PrometheusWriter
func NewPrometheusWriter(remote string, sync bool) TestReporter {
	return &prometheusReporter{
		remote: remote,
		sync:   sync,
	}
}

const namespace = "api_testing"

// PutRecord puts the test result into the Prometheum Pushgateway
func (w *prometheusReporter) PutRecord(record *ReportRecord) {
	var wait sync.WaitGroup
	wait.Add(1)

	go func() {
		defer wait.Done()

		responseTime := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			Name:        "response_time",
			ConstLabels: getConstLabels(record),
			Help:        "The response time in milliseconds of the API.",
		})
		responseTime.Set(float64(record.EndTime.Sub(record.BeginTime).Milliseconds()))

		pusher := push.New(w.remote, "api-testing").Collector(responseTime)

		if record.Error != nil {
			if w.errorCount == nil {
				w.errorCount = prometheus.NewCounter(prometheus.CounterOpts{
					Namespace:   namespace,
					Name:        "error_count",
					ConstLabels: getConstLabels(record),
				})
			}
			w.errorCount.Inc()
			pusher.Collector(w.errorCount)
		} else {
			if w.successCount == nil {
				w.successCount = prometheus.NewCounter(prometheus.CounterOpts{
					Namespace:   namespace,
					Name:        "success_count",
					ConstLabels: getConstLabels(record),
				})
			}
			w.successCount.Inc()
			pusher.Collector(w.successCount)
		}

		if err := pusher.Push(); err != nil {
			log.Println("Could not push completion time to Pushgateway:", err)
		}
	}()
	if w.sync {
		wait.Wait()
	}
	return
}

func (r *prometheusReporter) GetResourceUsage() []ResourceUsage {
	return nil
}

func getConstLabels(record *ReportRecord) prometheus.Labels {
	return prometheus.Labels{
		"group":  record.Group,
		"name":   record.Name,
		"api":    record.API,
		"method": record.Method,
	}
}

func (w *prometheusReporter) GetAllRecords() []*ReportRecord {
	// no support
	return nil
}

func (r *prometheusReporter) ExportAllReportResults() (result ReportResultSlice, err error) {
	// no support
	return
}
