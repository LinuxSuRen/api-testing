package limit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	limiter := NewDefaultRateLimiter(1, 1)
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
	}
	assert.True(t, num <= 10)
}
