package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type stdout struct {
	searches []*search
}

func (out *stdout) collect(s *search) {
	out.searches = append(out.searches, s)
}

func (out *stdout) print() {
	for i, s := range out.searches {
		if i > 0 {
			fmt.Println()
		}

		dash := strings.Repeat("-", (20 - len(s.keyword)))
		color.Red(s.keyword + " " + dash)

		for _, e := range s.entries {
			fmt.Printf("%-13s %s  ", e.time, e.title)
			color.Green(e.price)
		}
	}
}
