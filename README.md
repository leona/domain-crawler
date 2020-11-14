# Domain Crawler

This application crawls websites for URLs and stores the domains in a nested key value database for reverse subdomain searches.

The database consists of 2 root buckets, uncrawled and all. See below for an example.

```
UNCRAWLED
    google.com
    facebook.com
    docs.google.com
ALL
    com
        google
            sheets
            docs
        facebook
    co
        uk
            amazon
```

![image](https://github.com/leona/domain-crawler/blob/master/screenshot.png?raw=true)

## Usage
Download the binary [here](./bin)
Start crawling using default options
```./crawler```

Use flag -h to show all arguments
```./crawler -h```

Example
```./crawler -seed bbc.com -threads 200```

While the crawler is stopped you can get stats or make a query using
```./crawler -search google.com```
```./crawler -stat```

## Build from source

Requirements
* Golang 1.15

```make build-linux```