//go:build longtest
// +build longtest

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
