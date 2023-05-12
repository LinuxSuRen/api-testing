// Package runner provides the common expr style functions
package runner

import (
	"fmt"
	"time"
)

// ExprFuncSleep is a expr function for sleeping
func ExprFuncSleep(params ...interface{}) (res interface{}, err error) {
	if len(params) < 1 {
		err = fmt.Errorf("the duration param is required")
		return
	}

	switch duration := params[0].(type) {
	case int:
		time.Sleep(time.Duration(duration) * time.Second)
	case string:
		var dur time.Duration
		if dur, err = time.ParseDuration(duration); err == nil {
			time.Sleep(dur)
		}
	}
	return
}
