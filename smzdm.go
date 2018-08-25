package main

import (
	"flag"
	"fmt"
)

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

	for _, k := range keywords {
		s := new(search).ing(k)
		result.collect(s)
	}

	result.print()
}
