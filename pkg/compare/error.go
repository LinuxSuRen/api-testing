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

package compare

import (
	"fmt"
	"strings"
)

type noEqualErr struct {
	Path    []string
	Message string
	NeqErrs *noEqualErrs
}

type noEqualErrs struct {
	errs []error
}

type task struct {
	path []string
	errs []error
}

func (e *noEqualErr) Error() string {
	if e.NeqErrs == nil {
		return fmt.Sprintf("compare: field %s: %s", strings.Join(e.Path, "."), e.Message)
	}

	msg := []string{}
	q := []*task{{
		path: e.Path,
		errs: e.NeqErrs.errs,
	}}
	iter := 0
	end := len(e.NeqErrs.errs)
	for {
		if iter == end {
			iter = 0

			q[0] = nil
			q = q[1:]
			if len(q) == 0 {
				break
			}

			end = len(q[0].errs)
			continue
		}

		v := q[0].errs[iter]
		switch v.(type) {
		case *noEqualErr:
			if v.(*noEqualErr).NeqErrs != nil {
				q = append(q, &task{append(q[0].path, v.(*noEqualErr).Path...), v.(*noEqualErr).NeqErrs.errs})
			} else {
				msg = append(msg,
					fmt.Sprintf("compare: field %s: %s",
						strings.Join(append(q[0].path, v.(*noEqualErr).Path...), "."), v.(*noEqualErr).Message))
			}
		case *noEqualErrs:
			q = append(q, &task{path: q[0].path, errs: v.(*noEqualErrs).errs})
		}
		iter++
	}
	return strings.Join(msg, "\n")
}

// newNoEqualErr returns a NoEqualErr.
func newNoEqualErr(field string, err error) error {
	noEqerr := &noEqualErr{
		Path:    []string{field},
		Message: "",
	}

	if err != nil {
		if v, ok := err.(*noEqualErr); ok {
			noEqerr.Path = append(noEqerr.Path, v.Path...)
			if v.Message != "" {
				noEqerr.Message = v.Message
			}
		} else if v, ok := err.(*noEqualErrs); ok {
			noEqerr.NeqErrs = v
		} else {
			noEqerr.Message = err.Error()
		}
	}
	return noEqerr
}

// JoinErr returns an error that warp the given errors.
func JoinErr(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	e := &noEqualErrs{
		errs: make([]error, 0, n),
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return e
}

func (e *noEqualErrs) Error() string {
	var b []byte
	for i, err := range e.errs {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}
