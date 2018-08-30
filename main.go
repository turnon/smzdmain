package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
	"sync"
)

const keysWorkerCount int = 2

func main() {
	var result output
	toHTML := flag.Bool("h", false, "output html")
	webServer := flag.Bool("w", false, "work as webserver")
	flag.Parse()
	keywords := flag.Args()

	if *webServer {
		runServer()
		return
	}

	if len(keywords) <= 0 {
		panic("no keyword given")
	}

	if *toHTML {
		result = new(html)
	} else {
		result = new(stdout)
	}

	process(keywords, result).print()
}

func process(keywords []string, result output) output {
	keysCh := make(chan string)
	searchesCh := make(chan *search)
	var keysWg, resultWg sync.WaitGroup
	keysWg.Add(keysWorkerCount)
	resultWg.Add(1)

	for n := keysWorkerCount; n > 0; n-- {
		go func() {
			for k := range keysCh {
				searchesCh <- new(search).ing(k)
			}
			keysWg.Done()
		}()
	}

	go func() {
		for s := range searchesCh {
			result.collect(s)
		}
		resultWg.Done()
	}()

	for _, k := range keywords {
		keysCh <- k
	}

	close(keysCh)
	keysWg.Wait()
	close(searchesCh)
	resultWg.Wait()
	return result
}

func runServer() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		q := request.URL.Query()["q"]
		if len(q) == 0 {
			return
		}
		keys := strings.Split(q[0], " ")
		log.Println(request.RemoteAddr, keys)
		htmlOutput := new(html)
		process(keys, htmlOutput).print(writer)
	})
	http.ListenAndServe(":80", nil)
}
