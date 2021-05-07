package main

import (
	"flag"
	"net/http"
	"os"
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
		runServer(keywords)
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
	searching := make(chan *search)
	searched := make(chan *search)
	var searchWg, resultWg sync.WaitGroup
	searchWg.Add(keysWorkerCount)
	resultWg.Add(1)

	for n := keysWorkerCount; n > 0; n-- {
		go func() {
			for search := range searching {
				searched <- search.ing()
			}
			searchWg.Done()
		}()
	}

	go func() {
		for s := range searched {
			result.collect(s)
		}
		resultWg.Done()
	}()

	for i, k := range keywords {
		s := search{Index: i, Keyword: k}
		searching <- &s
	}

	close(searching)
	searchWg.Wait()
	close(searched)
	resultWg.Wait()
	result.sort()
	return result
}

func runServer(specificKeys ...[]string) {
	var handler func(writer http.ResponseWriter, request *http.Request)
	keys := specificKeys[0]
	if len(keys) > 0 {
		handler = static(keys)
	} else {
		handler = dyna
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port(), nil)
}

func port() string {
	if port := os.Getenv("PORT"); len(port) != 0 {
		return port
	}
	return "80"
}
