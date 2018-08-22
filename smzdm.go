package main

import (
	"fmt"
	"net/http"
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

func main() {
	resp, err := http.Get("http://search.smzdm.com/?c=home&s=%E5%BE%B7%E8%8A%99&v=b")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		panic(err)
	}

	doc.Find("#feed-main-list .z-feed-content").Each(func(i int, s *goquery.Selection) {
		record := (&entry{s}).toStr()
		fmt.Println(record)
	})
}
