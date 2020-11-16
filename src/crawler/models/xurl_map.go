package models

import (
	"errors"
)

type XurlMap map[string]*Xurl

func (self XurlMap) PopItem() (*Xurl, error) {
	for key, item := range self {
		delete(self, key)
		return item, nil
	}

	return nil, errors.New("No items to pop")
}

func (self XurlMap) UniqueHosts() ([]*Xurl) {
	unique := map[string]*Xurl{}
	output := []*Xurl{}
	
	for _, item := range self {
		_, exists := unique[item.Url.Host]

		if !exists {
			unique[item.Url.Host] = item
		}
	}

	for _, item := range unique {
		output = append(output, item)
	}

	return output
}