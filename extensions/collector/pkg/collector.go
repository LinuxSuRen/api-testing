package pkg

import (
	"fmt"
	"net/http"
	"sync"
)

// Collects is a HTTP request collector
type Collects struct {
	once       sync.Once
	signal     chan string
	stopSignal chan struct{}
	keys       map[string]*http.Request
	requests   []*http.Request
	events     []EventHandle
}

// NewCollects creates an instance of Collector
func NewCollects() *Collects {
	return &Collects{
		once:       sync.Once{},
		signal:     make(chan string, 5),
		stopSignal: make(chan struct{}, 1),
		keys:       make(map[string]*http.Request),
	}
}

// Add adds a HTTP request
func (c *Collects) Add(req *http.Request) {
	key := fmt.Sprintf("%s-%s", req.Method, req.URL.String())
	if _, ok := c.keys[key]; !ok {
		c.keys[key] = req
		c.requests = append(c.requests, req)
		c.signal <- key
	}
}

// EventHandle is the collect event handle
type EventHandle func(r *http.Request)

// AddEvent adds new event handle
func (c *Collects) AddEvent(e EventHandle) {
	c.events = append(c.events, e)
	c.handleEvents()
}

// Stop stops the collector
func (c *Collects) Stop() {
	c.stopSignal <- struct{}{}
}

func (c *Collects) handleEvents() {
	fmt.Println("handle events")
	c.once.Do(func() {
		go func() {
			fmt.Println("start handle events")
			for {
				select {
				case key := <-c.signal:
					fmt.Println("receive signal", key)
					for _, e := range c.events {
						fmt.Println("handle event", key, e)
						e(c.keys[key])
					}
				case <-c.stopSignal:
					fmt.Println("stop")
					return
				}
			}
		}()
	})
}
