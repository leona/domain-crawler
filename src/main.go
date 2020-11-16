package main

import (
	"fmt"
	"github.com/mbndr/figlet4go"
	"flag"
	"os"
	"github.com/leona/domain-crawler/src/crawler"
	"github.com/leona/domain-crawler/src/crawler/models"
	"github.com/leona/domain-crawler/src/crawler/utilities"
)

func workerContainer(host string, quit chan *models.WorkerShare) {
	workerShare := models.WorkerShare{
		Paths: make(map[string]*models.Xurl),
		Host: host,
		MaxDepth: *utilities.InputOptions.Depth,
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
	quit := make(chan *models.WorkerShare)

	for id := 0; id < *utilities.InputOptions.Threads; id++ {
        go func() {
			quit <- nil
		}()
	}

	for {
		share := <- quit

		if share != nil {
			hosts := share.Paths.UniqueHosts()
			count := crawler.Store.SaveHosts(hosts)

			utilities.Info(2, "Saving share for:", share.Host, len(hosts), "unique hosts.", count, "stored")
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
	if *utilities.InputOptions.Test == true {
		fmt.Println("Running tests")
		os.Exit(0)
	}

	if *utilities.InputOptions.Help == true {
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

	if len(*utilities.InputOptions.Seed) > 0 {
		utilities.Info(0, "Seed provided:", utilities.InputOptions.Seed)
		url := models.MakeXurl("http://" + *utilities.InputOptions.Seed)
		crawler.Store.SaveHosts([]*models.Xurl{url})
	}

	if len(*utilities.InputOptions.Search) > 0 {
		utilities.Info(0, "Searching DB")
		url := models.MakeXurl("http://" + *utilities.InputOptions.Search)

		results := crawler.Store.Search(url.FullComponents)

		for _, item := range results {
			fmt.Println(item)
		}
		
		utilities.Info(0, len(results), "results found")
		os.Exit(0)
	}

	if *utilities.InputOptions.Stat == true {
		allCount, uncrawledCount := crawler.Store.Stat()
		fmt.Println(allCount, "total -", uncrawledCount, "uncrawled.")
		os.Exit(0)
	}

	initWorkers()
}