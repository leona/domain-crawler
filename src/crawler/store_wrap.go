package crawler

import (
	_"fmt"
	"github.com/leona/domain-crawler/src/crawler/models"
	"github.com/leona/domain-crawler/src/crawler/utilities"
)

type StoreWrap struct {
	StoreType
}

var Store StoreWrap

func init() {
	Store = StoreWrap{}
	Store.init()
}

func (self StoreWrap) Search(components []string) ([]string) {
	results := self.getNestedBuckets(StoreKeyAll, components, 0)
	return results
}

func (self StoreWrap) Stat() (int, int) {
	allCount := len(self.getNestedBuckets(StoreKeyAll, []string{}, 0))
	uncrawledCount := len(self.getNestedBuckets(StoreKeyUncrawled, []string{}, 0))
	return allCount, uncrawledCount
}

func (self StoreWrap) searchLimit(components []string, limit int) ([]string) {
	results := self.getNestedBuckets(StoreKeyAll, components, limit)
	return results
}

func (self StoreWrap) ShouldCrawl(components []string) bool {
	searchResults := self.searchLimit(components, *utilities.InputOptions.Limit)

	if len(searchResults) >= *utilities.InputOptions.Limit {
		utilities.Info(3, "Skipping:", components, "Has:", len(searchResults), "previous results")
		return false
	} else {
		utilities.Info(3, "Not skipping:", components, "Has:", len(searchResults), "limit:", *utilities.InputOptions.Limit)
	}

	return true
}

func (self StoreWrap) SaveHosts(hosts []*models.Xurl) int {
	counter := 0

	for _, host := range hosts {
		current := self.getNestedBucket(StoreKeyAll, host.FullComponents)

		if current == nil {
			if err := self.createNestedBucket(StoreKeyAll, host.FullComponents); err != nil {
				continue
			}

			counter += 1

			if self.ShouldCrawl(host.RootComponents) == false {
				continue
			}

			utilities.Info(3, "Adding to crawl list:", host.Url.Host)
			self.put(StoreKeyUncrawled, host.Url.Host)
		} else {
			utilities.Info(3, host.Url.Host, "already exists")
		}
	}

	return counter
}