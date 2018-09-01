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
