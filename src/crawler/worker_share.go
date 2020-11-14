package crawler

import (
	_ "fmt"
	"net/url"
	"sync"
	"errors"
)

type WorkerShare struct {
	Paths UrlMap
	Uncrawled UrlMap
	Host string
	MaxDepth int
	Depth int
	mu sync.Mutex
}

func (self *WorkerShare) Init() {
	if len(self.Uncrawled) == 0 {
		self.Uncrawled = make(UrlMap)
		host := "http://" + self.Host
		url, _ := url.Parse(host)
		self.Uncrawled[host] = url
	}
}

func (self *WorkerShare) appendExternal(urlString string, url *url.URL) {
	self.mu.Lock()
	self.Paths[urlString] = url
	self.mu.Unlock()
}

func (self *WorkerShare) append(urlString string, url *url.URL) {
	self.mu.Lock()
	self.Paths[urlString] = url
	self.Uncrawled[urlString] = url
	self.mu.Unlock()
}

func (self *WorkerShare) popItem() (*url.URL, error) {
	self.mu.Lock()
	
	output, err := self.Uncrawled.popItem()

	if err != nil {
		self.mu.Unlock()
		return output, errors.New("No uncrawled")
	}

	self.Depth += 1	
	self.mu.Unlock()
	return output, nil
}