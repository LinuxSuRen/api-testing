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

// Package compare provides basic functions to compare JSON object and array.
// Currently we use gjson to store the JSON data. Please check tidwall/gjson.
package compare

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/tidwall/gjson"
)

// Object compares two JSON object.
func Object(field string, expect, actul map[string]gjson.Result) error {
	var errs error
	for k, ev := range expect {
		av, ok := actul[k]
		if !ok {
			errs = JoinErr(errs, newNoEqualErr(field, fmt.Errorf("field %s is not exist", k)))
			continue
		}

		err := Element(k, ev, av)
		if err != nil {
			errs = JoinErr(errs, newNoEqualErr(field, err))
		}
	}
	return errs
}

// Array compares two JSON Array.
func Array(field string, expect, actul []gjson.Result) error {
	var errs error
	if l1, l2 := len(expect), len(actul); l1 != l2 {
		return newNoEqualErr(field, fmt.Errorf("length is not equal, expect %v fields but got %v", l1, l2))
	}

	for i := range expect {
		err := Element(strconv.Itoa(i), expect[i], actul[i])
		if err != nil {
			errs = JoinErr(errs, newNoEqualErr(field, err))
		}
	}
	return errs
}

// Element compares two JSON Element.
func Element(field string, expect, actul gjson.Result) error {
	if expect.Type != actul.Type {
		return newNoEqualErr(field, fmt.Errorf("expect type %s but got %v", expect.Type.String(), actul.Type.String()))
	}

	if expect.IsObject() {
		return Object(field, expect.Map(), actul.Map())
	}

	if expect.IsArray() {
		return Array(field, expect.Array(), actul.Array())
	}

	if !reflect.DeepEqual(expect.Value(), actul.Value()) {
		return newNoEqualErr(field, fmt.Errorf("expect %v but got %v", expect.Value(), actul.Value()))
	}

	return nil
}
