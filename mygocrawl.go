package main

import (
	"mygocrawl/conf"
	. "mygocrawl/crawl"
	"mygocrawl/log"
	"regexp"
)

func main() {

	var cc crawlCof
	var err error

	cf, _ := conf.ReadConfigFile("conf/config.cfg")

	if cc.host, err = cf.GetString("default", "host"); err != nil {
		log.Error("read host error...")
	}
	//println(cc.host)

	regex := regexp.MustCompile("<a.*?href=[\"'](http.*?)[\"']")

	curl := make(chan []byte)
	csite := make(chan Site)
	death := make(chan string)

	// Give our crawler a place to start.
	go Seed(curl, cc.host)

	// Keeps track of which urls we have visted.
	visited := make(map[string]int)

	// Start the throttled crawling.
	go ThrottledCrawl(curl, csite, death, visited)

	// Main loop that never exits and blocks on the data of a page.
	for {
		site := <-csite
		go GetUrls(curl, site, regex)
	}
}
