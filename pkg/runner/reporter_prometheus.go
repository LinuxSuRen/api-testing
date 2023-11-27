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

		var responseTime prometheus.Gauge
		responseTime = prometheus.NewGauge(prometheus.GaugeOpts{
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
