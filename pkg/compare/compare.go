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
