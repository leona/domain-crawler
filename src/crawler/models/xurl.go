package models

import (
	"github.com/leona/domain-crawler/src/crawler/utilities"
	"golang.org/x/net/publicsuffix"
	"strings"
	"net/url"
	"path/filepath"
)

type Xurl struct {
	Suffix []string
	FullComponents []string
	Components []string
	RootComponents []string
	Raw string
	Extension string
	Url *url.URL
}

func MakeXurl(raw string) *Xurl {
	xurl := new(Xurl)
	xurl.Raw = raw
	xurl.Build()
	return xurl
}

func (self Xurl) IsAsset() bool {
    switch self.Extension {
    case
        ".jpg",
        ".jpeg",
		".gif",
		".mp4",
		".png",
		".mp3",
		".pdf",
		".css",
		".js",
		".webp",
        ".svg":
        return true
    }
    return false
}

func (self *Xurl) Build() {
	// Parse the url
	if err := self.Parse(); err != nil {
		utilities.Info(2, "Failed to parse url:", self.Raw)
		return
	}

	// Get the suffix
	suffix, _ := publicsuffix.PublicSuffix("http://" + self.Url.Host)
	self.Suffix = strings.Split(suffix, ".")

	// Store each component of the domain before the suffix and the full components
	components := strings.Split(self.Url.Host, ".") 
	self.FullComponents = utilities.Reverse(components)

	if len(self.FullComponents) < len(self.Suffix) + 1 {
		utilities.Info(2, "Invalid URL:", self.Raw, self.FullComponents, self.Suffix)
	} else {
		self.RootComponents = self.FullComponents[0:len(self.Suffix) + 1]
	}

	self.Components = self.FullComponents[len(self.Suffix):]
	
	// Store file extension if any
	path := self.Url.EscapedPath()
	self.Extension = filepath.Ext(path)
}

func (self *Xurl) Parse() error  {
	self.Raw = strings.Replace(self.Raw, "https://", "http://", -1)
	parsedUrl, err := url.Parse(self.Raw)

	if err != nil {
		utilities.Info(3, "Failed to parse:", self.Raw)
		return err
	}

	requestUri := parsedUrl.RequestURI()

	if requestUri == "/" {
		requestUri = ""
	}

	self.Raw = "http://" + parsedUrl.Host + requestUri
	self.Url, err = url.Parse(self.Raw)
	return err
}