package main

import (
	"fmt"
	"net/url"
	"github.com/mbndr/figlet4go"
	"flag"
	"os"
	"github.com/leona/domain-crawler/src/crawler"
)

func workerContainer(host string, quit chan *crawler.WorkerShare) {
	workerShare := crawler.WorkerShare{
		Paths: make(map[string]*url.URL),
		Host: host,
		MaxDepth: *crawler.InputOptions.Depth,
	}

	workerShare.Init()
	
	worker := crawler.Worker{
		Share: &workerShare,
	}

	worker.Init()
	worker.Start()
	quit <- &workerShare
}

func initWorkers() {
	fmt.Println("Starting workers")
	quit := make(chan *crawler.WorkerShare)

	for id := 0; id < *crawler.InputOptions.Threads; id++ {
        go func() {
			quit <- nil
		}()
	}

	for {
		share := <- quit

		if share != nil {
			hosts := share.Paths.UniqueHosts()
			count := crawler.Store.SaveHosts(hosts)

			crawler.Info(2, "Saving share for:", share.Host, len(hosts), "unique hosts.", count, "stored")
		}

		host := crawler.Store.Pop(crawler.StoreKeyUncrawled, 1)

		if len(host) == 0 {
			fmt.Println("No uncrawled found")
			continue
		}

		go workerContainer(host[0], quit)
	}
}

func main() {
	if *crawler.InputOptions.Help == true {
		ascii := figlet4go.NewAsciiRender()

		options := figlet4go.NewRenderOptions()
		options.FontColor = []figlet4go.Color{
			figlet4go.ColorWhite,
		}
		renderStr, _ := ascii.RenderOpts("Domain Crawler", options)

		fmt.Print(renderStr)
		fmt.Println("github.com/leona/domain-crawler\n\nArguments\n")
		flag.PrintDefaults()
		fmt.Println("\n")
		os.Exit(0)
	}

	if len(*crawler.InputOptions.Seed) > 0 {
		crawler.Info(0, "Seed provided:", crawler.InputOptions.Seed)
		crawler.Store.SaveHosts([]string{*crawler.InputOptions.Seed})
	}

	if len(*crawler.InputOptions.Search) > 0 {
		crawler.Info(0, "Searching DB")

		results := crawler.Store.Search(*crawler.InputOptions.Search)

		for _, item := range results {
			fmt.Println(item)
		}
		
		crawler.Info(0, len(results), "results found")
		os.Exit(0)
	}

	if *crawler.InputOptions.Stat == true {
		allCount, uncrawledCount := crawler.Store.Stat()
		fmt.Println(allCount, "total -", uncrawledCount, "uncrawled.")
		os.Exit(0)
	}

	initWorkers()
}