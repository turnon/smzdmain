package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
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

func static(keywords []string) func(writer http.ResponseWriter, request *http.Request) {
	var (
		timer      <-chan time.Time
		m          sync.Mutex
		htmlOutput output
	)

	fetch := func() {
		m.Lock()
		defer m.Unlock()

		timer = time.After(60 * time.Second)
		htmlOutput = new(html)
		process(keywords, htmlOutput)
		log.Println("pre-fetched")
	}

	fetch()

	go func() {
		for {
			<-timer
			fetch()
		}
	}()

	return func(writer http.ResponseWriter, request *http.Request) {
		m.Lock()
		defer m.Unlock()
		htmlOutput.print(writer)
	}
}

func dyna(writer http.ResponseWriter, request *http.Request) {
	q := request.URL.Query()["q"]
	if len(q) == 0 {
		return
	}
	keys := strings.Split(q[0], " ")
	log.Println(request.RemoteAddr, keys)
	htmlOutput := new(html)
	process(keys, htmlOutput).print(writer)
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
	http.ListenAndServe(":80", nil)
}
