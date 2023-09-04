package limit

import (
	"sync"
	"time"
)

type RateLimiter interface {
	Accept() bool
	Stop()
	Burst() int32
}

type defaultRateLimiter struct {
	qps       int32
	burst     int32
	lastToken time.Time
	singal    chan struct{}
	mu        sync.Mutex
}

func NewDefaultRateLimiter(qps, burst int32) RateLimiter {
	if qps <= 0 {
		qps = 5
	}
	if burst <= 0 {
		burst = 5
	}
	limiter := &defaultRateLimiter{
		qps:       qps,
		burst:     burst,
		singal:    make(chan struct{}, 1),
		lastToken: time.Now(),
	}
	go limiter.updateBurst()
	return limiter
}

func (r *defaultRateLimiter) resver() (delay time.Duration, ok bool) {
	delay = time.Now().Sub(r.lastToken) / time.Millisecond
	delayRequire := time.Second / time.Duration(r.qps)
	r.lastToken = time.Now()
	if delay >= delayRequire {
		ok = true
	} else if r.Burst() > 0 {
		r.Setburst(r.Burst() - 1)
		ok = true
	} else {
		delay = delayRequire - delay
	}
	return
}

func (r *defaultRateLimiter) Accept() bool {
	delay, ok := r.resver()
	if ok {
		return ok
	}

	if delay > 0 {
		time.Sleep(delay)
	}
	return ok
}

func (r *defaultRateLimiter) Setburst(burst int32) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.burst = burst
}

func (r *defaultRateLimiter) Burst() int32 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.burst
}

func (r *defaultRateLimiter) Stop() {
	r.singal <- struct{}{}
}

func (r *defaultRateLimiter) updateBurst() {
	for {
		select {
		case <-time.After(time.Second):
			r.Setburst(r.Burst() + r.qps)
		case <-r.singal:
			return
		}
	}
}
