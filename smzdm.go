package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	keywords := flag.Args()

	if len(keywords) <= 0 {
		fmt.Println("no keyword given")
		return
	}

	var result output = new(html)

	for _, k := range keywords {
		s := new(search).ing(k)
		result.collect(s)
	}

	result.print()
}
