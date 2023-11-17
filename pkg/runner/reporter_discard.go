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

type discardTestReporter struct {
}

// NewDiscardTestReporter creates a test reporter which discard everything
func NewDiscardTestReporter() TestReporter {
	return &discardTestReporter{}
}

// PutRecord does nothing
func (r *discardTestReporter) PutRecord(*ReportRecord) {
	// Do nothing which is the design purpose
}

// GetAllRecords does nothing
func (r *discardTestReporter) GetAllRecords() []*ReportRecord {
	return nil
}

// ExportAllReportResults does nothing
func (r *discardTestReporter) ExportAllReportResults() (ReportResultSlice, error) {
	return nil, nil
}
