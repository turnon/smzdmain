package main

import (
	"flag"
	"fmt"
	"sync"
)

const keysWorkerCount int = 2

func main() {
	toHTML := flag.Bool("h", false, "output html")
	flag.Parse()
	keywords := flag.Args()

	if len(keywords) <= 0 {
		fmt.Println("no keyword given")
		return
	}

	var result output
	if *toHTML {
		result = new(html)
	} else {
		result = new(stdout)
	}

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
	result.print()
}
