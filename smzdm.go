package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type entry struct {
	s *goquery.Selection
}

func (e *entry) toStr() string {
	title := strings.TrimSpace(e.s.Find(".feed-block-title a").First().Text())
	price := strings.TrimSpace(e.s.Find(".feed-block-title a div").First().Text())
	timeBlock := e.s.Find(".feed-block-extras").First()
	timeBlock.Children().Remove()
	time := strings.TrimSpace(timeBlock.Text())
	return time + " " + title + " " + price
}

type search struct {
	keyword string
}

func (s *search) process() {
	keyword := url.QueryEscape(s.keyword)

	resp, err := http.Get("http://search.smzdm.com/?c=home&s=" + keyword + "&v=b")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		panic(err)
	}

	s.printKeyword()

	doc.Find("#feed-main-list .z-feed-content").Each(func(i int, s *goquery.Selection) {
		record := (&entry{s}).toStr()
		fmt.Println(record)
	})
}

func (s *search) printKeyword() {
	dash := strings.Repeat("-", (20 - len(s.keyword)))
	fmt.Println(s.keyword + " " + dash)
}

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
