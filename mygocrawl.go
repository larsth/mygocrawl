package main

import (
	"log"
	"regexp"
	"os"
	"bufio"
	"bytes"
	"io"
)

const envVariablePrefix string = "MYGOCRAWL_"
const fatalErrorSeedsEnvVarFormatString string = "Fatal error: Enviroment variable %s is not set or its content is the empty string"
const errorSeedsFileFormatString string = "An I/O error ocurred during reeding seeds text file: %s, Error %s"


func readLines(fileName string) ([]string, error) {
    var (
	file *os.File
	lines [] string = make([]string, 0)
	err error
        part []byte
        prefix bool
    )
    if file, err = os.Open(fileName); err != nil {
        return lines, err
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    buffer := bytes.NewBuffer(make([]byte, 0))
    for {
        if part, prefix, err = reader.ReadLine(); err != nil {
            break
        }
        buffer.Write(part)
        if !prefix {
            lines = append(lines, buffer.String())
            buffer.Reset()
        }
    }
    if err == io.EOF {
        err = nil
    }
    return lines, err
}

func main() {
	seedsEnvVarName := envVariablePrefix + "SEEDS_FILENAME"
	seedsFileName := os.Getenv(seedsEnvVarName)
	if seedsFileName == "" {
		log.Fatalf(fatalErrorSeedsEnvVarFormatString, seedsEnvVarName)
	}

	lines, err := readLines(seedsFileName)
	if err != nil {
		log.Fatal(err.Error())
	}

	regex := regexp.MustCompile("<a.*?href=[\"'](http.*?)[\"']")

	curl := make(chan []byte)
	csite := make(chan Site)
	death := make(chan string)
	
	if len(lines)== 0 {
		log.Fatal("No seeds!")
	}

	// Give our crawler a place to start.
	go Seed(curl, lines[0])
	
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
