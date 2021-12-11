package crawler

import (
	"flag"
)

type InputOptionsType struct {
	Threads *int
	Db *string
	Seed *string
	Search *string
	Stat *bool
	Depth *int
	Help *bool
	Verbose *int
	Limit *int
}

var InputOptions InputOptionsType

func init() {
	InputOptions = InputOptionsType{
		Threads: flag.Int("threads", 5, "number of threads"),
		Db: flag.String("db", "crawler", "filename of database"),
		Seed: flag.String("seed", "", "initial seed host"),
		Search: flag.String("search", "", "perform a domain search on the db"),
		Stat: flag.Bool("stat", false, "return stats about the db"),
		Depth: flag.Int("depth", 5, "max number of pages to crawl on a single host"),
		Help: flag.Bool("h", false, "show help and exit"),
		Verbose: flag.Int("v", 2, "Log verbosity level 0-3"),
		Limit: flag.Int("limit", 500, "top level subdomain crawl limit"),
	}

	flag.Parse()
}