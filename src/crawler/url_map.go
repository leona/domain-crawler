package crawler

import (
	"net/url"
	"errors"

)

type UrlMap map[string]*url.URL

func (self UrlMap) popItem() (*url.URL, error) {
	for key, item := range self {
		delete(self, key)
		return item, nil
	}

	return nil, errors.New("No items to pop")
}

func (self UrlMap) UniqueHosts() ([]string) {
	unique := map[string]*url.URL{}
	output := []string{}
	
	for _, item := range self {
		_, exists := unique[item.Host]

		if !exists {
			unique[item.Host] = item
		}
	}

	for key, _ := range unique {
		output = append(output, key)
	}

	return output
}