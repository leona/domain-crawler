package crawler

import (
	"github.com/leona/domain-crawler/src/crawler/utilities"
	"github.com/leona/domain-crawler/src/crawler/models"
	"time"
	"net/http"
	"io/ioutil"
	"regexp"
	"net"
)

type Worker struct {
	urlPattern *regexp.Regexp
	Share *models.WorkerShare
}

var client *http.Client
var transport *http.Transport

func init() {
	utilities.Info(0, "Creating HTTP Client")
	
	transport = &http.Transport{
		MaxIdleConns:          *utilities.InputOptions.Threads,
		MaxIdleConnsPerHost: *utilities.InputOptions.Threads,
		MaxConnsPerHost: 0,
		IdleConnTimeout:       3 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 3 * time.Second,
		ResponseHeaderTimeout: 3 * time.Second,
		DisableCompression: false,
		ForceAttemptHTTP2:     false,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
			DualStack: true,
		}).DialContext,
	}
	
	client = &http.Client{Transport: transport}	
}

func (self * Worker) Init() {
	utilities.Info(3, "Starting working for:", self.Share.Host)
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
		utilities.Info(2, "Errror making request to:", url, err)
		return []byte{}, err
	}
	
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (self *Worker) addUrl(urlString string) {
	xurl := models.MakeXurl(urlString)

	if xurl.Url == nil {
		return
	}

	if xurl.Url.Host != self.Share.Host {
		self.Share.AppendExternal(xurl)
		return
	}

	self.Share.Append(xurl)
}

func (self *Worker) Start() {
	for {
		if self.Share.MaxDepth > 0 && self.Share.Depth > self.Share.MaxDepth {
			utilities.Info(3, "Reached max depth for worker:", self.Share.Host)
			return
		}

		xurl, err := self.Share.PopItem()

		if err != nil {
			utilities.Info(3, "No more work for:", self.Share.Host, err)
			return
		}

		if xurl.IsAsset() {
			utilities.Info(3, "Invalid extension for url:", xurl.Raw)
			continue
		}

		utilities.Info(3, "Crawling:", xurl.Raw, "For worker:", self.Share.Host)
		body, err := self.getBody(xurl.Raw)
		
		if err != nil {
			utilities.Info(3, "Errror getting response body for:", xurl.Raw)
			continue
		}

		urls := self.getBodyUrls(body) 

		if len(urls) == 0 {
			utilities.Info(3, "No urls for:", xurl.Raw)
		}
		
		for _, url := range urls {
			self.addUrl(url)
		}
	}
}