package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

type entry struct {
	s *goquery.Selection
}

func (e *entry) printf() {
	title := strings.TrimSpace(e.s.Find(".feed-block-title a").First().Text())
	price := strings.TrimSpace(e.s.Find(".feed-block-title a div").First().Text())
	timeBlock := e.s.Find(".feed-block-extras").First()
	timeBlock.Children().Remove()
	time := strings.TrimSpace(timeBlock.Text())
	fmt.Printf("%-13s %s  ", time, title)
	color.Green(price)
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

	s.printkeyword()

	doc.Find("#feed-main-list .z-feed-content").Each(func(i int, s *goquery.Selection) {
		(&entry{s}).printf()
	})
}

func (s *search) printkeyword() {
	dash := strings.Repeat("-", (20 - len(s.keyword)))
	color.Red(s.keyword + " " + dash)
}
