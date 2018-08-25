package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type entry struct {
	title, price, time string
}

func (e *entry) extract(s *goquery.Selection) *entry {
	e.title = strings.TrimSpace(s.Find(".feed-block-title a").First().Text())
	e.price = strings.TrimSpace(s.Find(".feed-block-title a div").First().Text())
	timeBlock := s.Find(".feed-block-extras").First()
	timeBlock.Children().Remove()
	e.time = strings.TrimSpace(timeBlock.Text())
	return e
}

type search struct {
	keyword string
	entries []*entry
}

func (s *search) ing(k string) *search {
	s.keyword = k

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

	doc.Find("#feed-main-list .z-feed-content").Each(func(i int, selection *goquery.Selection) {
		e := new(entry).extract(selection)
		s.entries = append(s.entries, e)
	})

	return s
}
