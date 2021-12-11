package crawler

import (
	_ "fmt"
	"time"
	"net/http"
	"io/ioutil"
	"regexp"
	"net/url"
	"log"
	"path/filepath"
	"strings"
)

type Worker struct {
	urlPattern *regexp.Regexp
	Share *WorkerShare
}

var client *http.Client
var transport *http.Transport

func init() {
	transport = &http.Transport{
		MaxIdleConns:          100000,
		IdleConnTimeout:       3 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 3 * time.Second,
		DisableCompression: false,
	}
	
	client = &http.Client{Transport: transport}	
}

func (self * Worker) Init() {
	Info(3, "Starting working for:", self.Share.Host)
	urlPattern, _ := regexp.Compile(URL_REGEX)
	self.urlPattern = urlPattern
}

func (self * Worker) getBodyUrls(body []byte) []string {
	if match := self.urlPattern.FindAllString(string(body), -1); len(match) > 0 {
		return match
	}

	return []string{}
}

func (self * Worker) getBody(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", USER_AGENT)
	resp, err := client.Do(req)

	if err != nil {
		Info(3, "Errror making request to:", url)
		return []byte{}, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (self *Worker) addUrl(urlString string) {
	urlString = strings.Replace(urlString, "https://", "http://", -1)
	url, err := url.Parse(urlString)

	if err != nil {
		log.Print(err)
		return
	}

	requestUri := url.RequestURI()

	if requestUri == "/" {
		requestUri = ""
	}

	urlString = "http://" + url.Host + requestUri
	url, _ = url.Parse(urlString)

	if url.Host != self.Share.Host {
		self.Share.appendExternal(urlString, url)
		return
	}

	self.Share.append(urlString, url)
}

func (self *Worker) Start() {
	for {
		if self.Share.MaxDepth > 0 && self.Share.Depth > self.Share.MaxDepth {
			Info(3, "Reached max depth for worker:", self.Share.Host)
			return
		}

		url, err := self.Share.popItem()

		if err != nil {
			Info(3, "No more work for:", self.Share.Host)
			return
		}

		path := url.EscapedPath()
		extension := filepath.Ext(path)

		if isInvalidExtension(extension) {
			Info(3, "Invalid extension for url:", url.String())
			continue
		}

		Info(3, "Crawling:", url.String(), "For worker:", self.Share.Host)
		body, err := self.getBody(url.String())
		
		if err != nil {
			Info(3, "Errror getting response body for:", url.String())
			continue
		}

		urls := self.getBodyUrls(body) 

		if len(urls) == 0 {
			Info(3, "No urls for:", url.String())
		}
		
		for _, url := range urls {
			self.addUrl(url)
		}
	}
}