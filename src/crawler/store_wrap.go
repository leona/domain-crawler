package crawler

import (
	"strings"
	_"fmt"
)

type StoreWrap struct {
	StoreType
}

var Store StoreWrap

func init() {
	Store = StoreWrap{}
	Store.init()
}

func (self StoreWrap) splitHost(host string) ([]string) {
	components := strings.Split(host, ".") 
	reversed := reverse(components)
	return reversed
}

func (self StoreWrap) Search(host string) ([]string) {
	components := self.splitHost(host)
	results := self.getNestedBuckets(StoreKeyAll, components, 0)
	return results
}

func (self StoreWrap) Stat() (int, int) {
	allCount := len(self.getNestedBuckets(StoreKeyAll, []string{}, 0))
	uncrawledCount := len(self.getNestedBuckets(StoreKeyUncrawled, []string{}, 0))
	return allCount, uncrawledCount
}

func (self StoreWrap) searchLimit(host string, limit int) ([]string) {
	components := self.splitHost(host)
	results := self.getNestedBuckets(StoreKeyAll, components, limit)
	return results
}

func (self StoreWrap) SaveHosts(hosts []string) int {
	counter := 0

	for _, host := range hosts {
		components := self.splitHost(host)
		current := self.getNestedBucket(StoreKeyAll, components)

		if current == nil {
			counter += 1
			self.createNestedBucket("all", components)

			if isTopLevel(components[0]) {
				searchQuery := strings.Join(reverse(components[0:2]), ".")
				searchResults := self.searchLimit(searchQuery, *InputOptions.Limit)

				if len(searchResults) >= *InputOptions.Limit {
					Info(3, "Skipping top level:", searchQuery, "Has:", len(searchResults), "previous results")
					continue
				}
			}

			if !isTopLevel(components[0]) && len(components) >= 4 {
				searchQuery := strings.Join(reverse(components[0:3]), ".")
				searchResults := self.searchLimit(searchQuery, *InputOptions.Limit)

				if len(searchResults) >= *InputOptions.Limit {
					Info(3, "Skipping 3rd level:", searchQuery, "Has:", len(searchResults), "previous results")
					continue
				}
			}

			self.put("uncrawled", host)
		}
	}

	return counter
}