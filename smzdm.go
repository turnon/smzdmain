package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

	doc.Find(".feed-block-title").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("a").First().Text())
		price := strings.TrimSpace(s.Find("a div").First().Text())
		fmt.Println(title + " " + price)
	})
}
