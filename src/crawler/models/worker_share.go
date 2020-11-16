package models

import (
	_ "fmt"
	"sync"
	"errors"
)

type WorkerShare struct {
	Paths XurlMap
	Uncrawled XurlMap
	Host string
	MaxDepth int
	Depth int
	mu sync.Mutex
}

func (self *WorkerShare) Init() {
	if len(self.Uncrawled) == 0 {
		self.Uncrawled = make(XurlMap)
		host := "http://" + self.Host
		xurlMap := make(XurlMap)
		xurl := MakeXurl(host)
		xurlMap[host] = xurl
		self.Uncrawled = xurlMap
	}
}

func (self *WorkerShare) AppendExternal(url *Xurl) {
	self.mu.Lock()
	self.Paths[url.Raw] = url
	self.mu.Unlock()
}

func (self *WorkerShare) Append(url *Xurl) {
	self.mu.Lock()
	self.Paths[url.Raw] = url
	self.Uncrawled[url.Raw] = url
	self.mu.Unlock()
}

func (self *WorkerShare) PopItem() (*Xurl, error) {
	self.mu.Lock()
	
	output, err := self.Uncrawled.PopItem()

	if err != nil {
		self.mu.Unlock()
		return output, errors.New("No uncrawled")
	}

	self.Depth += 1	
	self.mu.Unlock()
	return output, nil
}