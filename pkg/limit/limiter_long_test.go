//go:build longtest
// +build longtest

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

package limit_test

import (
	"testing"
	"time"

	"github.com/linuxsuren/api-testing/pkg/limit"
	"github.com/stretchr/testify/assert"
)

func TestLimiterWithLongTime(t *testing.T) {
	for i := 0; i < 10; i++ {
		testLimiter(t, int32(8+i*2))
	}
}

func testLimiter(t *testing.T, count int32) {
	t.Log("test limit with count", count)
	limiter := limit.NewDefaultRateLimiter(count, 1)
	num := int32(0)

	loop := true
	go func(l limit.RateLimiter) {
		for loop {
			l.Accept()
			num += 1
		}
	}(limiter)

	select {
	case <-time.After(time.Second):
		loop = false
	}
	assert.True(t, num <= count+1, num)
}
