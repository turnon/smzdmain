package main

import (
	"flag"
	"fmt"
)

type output interface {
	collect(*search)
	print()
}

func main() {
	flag.Parse()
	keywords := flag.Args()

	if len(keywords) <= 0 {
		fmt.Println("no keyword given")
		return
	}

	var result output = new(stdout)

	for _, k := range keywords {
		s := new(search).ing(k)
		result.collect(s)
	}

	result.print()
}
