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

package downloader

import (
	"fmt"
	"io"
	"os"
)

type ProgressReader interface {
	io.Reader
	Withoutput(io.Writer)
	SetLength(int64)
	SetTitle(string)
}

type defaultProgressReader struct {
	reader      io.Reader
	w           io.Writer
	total       int64
	current     int
	lastPercent int64
	title       string
}

var _ io.Reader = (*defaultProgressReader)(nil)

func NewProgressReader(r io.Reader) ProgressReader {
	return &defaultProgressReader{
		reader: r,
		w:      os.Stdout,
	}
}

func (r *defaultProgressReader) SetTitle(title string) {
	r.title = title
}

func (r *defaultProgressReader) Read(p []byte) (count int, err error) {
	count, err = r.reader.Read(p)
	if r.total > 0 {
		if count > 0 {
			r.current += count
			newPercent := int64(r.current*100) / r.total
			if newPercent != int64(r.lastPercent) {
				r.lastPercent = newPercent
				fmt.Fprintf(r.w, "%s\t%d%%\n", r.title, newPercent)
			}
		}
	} else {
		fmt.Fprintf(r.w, "%d\n", count)
	}
	return
}

func (r *defaultProgressReader) Withoutput(w io.Writer) {
	r.w = w
}

func (r *defaultProgressReader) SetLength(len int64) {
	r.total = len
}
