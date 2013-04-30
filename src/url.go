package url

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

type Site struct {
	Url, Body []byte
}

const DATA = "url.data"

// GetUrl will make an HTTP GET request, build a site object and put it on a channel.
// It will send a message on the stop channel after the function finishes.
func GetUrl(url []byte, csite chan Site, death chan string) {

	resource := string(url)
	defer func() {
		death <- resource
	}()
	resp, err := http.Get(resource)
	if err != nil {
		fmt.Println("We have an error!: ", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Getting %v\n", resource)
	SaveData(resource)
	body, _ := ioutil.ReadAll(resp.Body)
	site := Site{url, body}
	csite <- site
}

func SaveData(url string) error {

	file, err := os.Open(DATA)
	if err != nil {

		return err
	}
	//defer file.Close()
	buf := bytes.NewBuffer(nil)
	buf.WriteString(url)
	buf.WriteTo(file)
	//println("-------------")
	//file.WriteString(url)
	file.Close()
	return nil
}

// ThrottledCrawl will limit the number of goroutines making requests.
// It will listen to a URL channel and spawn a goroutine for each URL.
// It manages the number of goroutines using a stop channel.
// This function does not return and should be used as a goroutine.
func ThrottledCrawl(curl chan []byte, csite chan Site, death chan string, visited map[string]int) {
	maxGos := 10
	numGos := 0
	for {
		if numGos > maxGos {
			<-death
			numGos -= 1
		}
		url := string(<-curl)
		if _, ok := visited[url]; !ok {
			go GetUrl([]byte(url), csite, death)
			numGos += 1
		}
		visited[url] += 1
	}
}

// Seed starts the crawling process by feeding the URL channel a URL.
func Seed(curl chan []byte) {
	curl <- []byte("http://www.huawei.com/cn/")
}

// GetUrls parses a site object and looks for links to sites.
func GetUrls(curl chan []byte, site Site, regex *regexp.Regexp) {
	matches := regex.FindAllSubmatch(site.Body, -1)
	for _, match := range matches {
		curl <- match[1]
	}
}
