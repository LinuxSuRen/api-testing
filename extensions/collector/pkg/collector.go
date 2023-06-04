package pkg

import (
	"fmt"
	"net/http"
	"sync"
)

type Collects struct {
	once       sync.Once
	signal     chan string
	stopSignal chan struct{}
	keys       map[string]*http.Request
	requests   []*http.Request
	events     []Event
}

func NewCollects() *Collects {
	return &Collects{
		once:       sync.Once{},
		signal:     make(chan string, 5),
		stopSignal: make(chan struct{}),
		keys:       make(map[string]*http.Request),
	}
}

func (c *Collects) Add(req *http.Request) {
	key := fmt.Sprintf("%s-%s", req.Method, req.URL.String())
	if _, ok := c.keys[key]; !ok {
		c.keys[key] = req
		c.requests = append(c.requests, req)
		c.signal <- key
	}
}

type Event func(r *http.Request)

func (c *Collects) AddEvent(e Event) {
	c.events = append(c.events, e)
	c.handleEvents()
}

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
