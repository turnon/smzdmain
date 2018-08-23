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

	for i, k := range keywords {
		if i > 0 {
			fmt.Println()
		}
		(&search{k}).process()
	}
}
