package limit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimiter(t *testing.T) {
	t.Log("run rate limit test")
	limiter := NewDefaultRateLimiter(0, 0)
	num := 0

	loop := true
	go func(l RateLimiter) {
		for loop {
			l.Accept()
			num += 1
		}
	}(limiter)

	select {
	case <-time.After(time.Second):
		loop = false
		limiter.Stop()
	}
	assert.True(t, num <= 10, num)
}
