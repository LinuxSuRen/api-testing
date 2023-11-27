package pkg

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Collects is a HTTP request collector
type Collects struct {
	once       sync.Once
	signal     chan string
	stopSignal chan struct{}
	keys       map[string]*RequestAndResponse
	requests   []*http.Request
	events     []EventHandle
}

type SimpleResponse struct {
	StatusCode int
	Body       string
}

type RequestAndResponse struct {
	Request  *http.Request
	Response *SimpleResponse
}

// NewCollects creates an instance of Collector
func NewCollects() *Collects {
	return &Collects{
		once:       sync.Once{},
		signal:     make(chan string, 5),
		stopSignal: make(chan struct{}, 1),
		keys:       make(map[string]*RequestAndResponse),
	}
}

// Add adds a HTTP request
func (c *Collects) Add(req *http.Request, resp *SimpleResponse) {
	key := fmt.Sprintf("%s-%s", req.Method, req.URL.String())
	if _, ok := c.keys[key]; !ok {
		c.keys[key] = &RequestAndResponse{
			Request:  req,
			Response: resp,
		}
		c.requests = append(c.requests, req)
		c.signal <- key
	}
}

// EventHandle is the collect event handle
type EventHandle func(r *RequestAndResponse)

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
	log.Println("handle events")
	c.once.Do(func() {
		go func() {
			log.Println("start handle events")
			for {
				select {
				case key := <-c.signal:
					log.Println("receive signal", key)
					for _, e := range c.events {
						e(c.keys[key])
					}
				case <-c.stopSignal:
					log.Println("stop")
					return
				}
			}
		}()
	})
}
