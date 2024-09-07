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
}

type defaultProgressReader struct {
	reader  io.Reader
	w       io.Writer
	total   int64
	current int
}

var _ io.Reader = (*defaultProgressReader)(nil)

func NewProgressReader(r io.Reader) ProgressReader {
	return &defaultProgressReader{
		reader: r,
		w:      os.Stdout,
	}
}

func (r *defaultProgressReader) Read(p []byte) (count int, err error) {
	count, err = r.reader.Read(p)
	if r.total > 0 {
		if count > 0 {
			r.current += count
			fmt.Println(count, "==", r.total, "==", r.current)
			fmt.Fprintf(r.w, "%d\n", int64(r.current*100)/r.total)
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
