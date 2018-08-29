package main

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	https        = "https:"
	root         = "http://search.smzdm.com/"
	query        = root + "/?v=b&c=home&s="
	base64Prefix = "data:image/jpg;base64,"
)

var officialImgServer = []string{"zdmimg.com", "smzdm.com"}

type entry struct {
	Title, Price, Time, Img, Href string
}

func (e *entry) extract(s *goquery.Selection) *entry {
	a := s.Find(".feed-block-title a").First()
	e.Href, _ = a.Attr("href")
	e.Title = strings.TrimSpace(a.Text())
	e.Price = strings.TrimSpace(s.Find(".feed-block-title a div").First().Text())
	timeBlock := s.Find(".feed-block-extras").First()
	timeBlock.Children().Remove()
	e.Time = strings.TrimSpace(timeBlock.Text())
	img, _ := s.Find("img").First().Attr("src")
	e.keepImg(https + img)
	return e
}

func (e *entry) keepImg(img string) {
	if !needConvert(img) {
		e.Img = img
		return
	}

	resp, err := http.Get(img)
	if err != nil {
		e.Img = img
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	content, _ := ioutil.ReadAll(reader)
	str := base64.StdEncoding.EncodeToString(content)
	e.Img = base64Prefix + str
}

func needConvert(img string) bool {
	for _, server := range officialImgServer {
		if strings.Contains(img, server) {
			return true
		}
	}
	return false
}

type search struct {
	doc     *goquery.Document
	Keyword string
	Entries []*entry
}

func (s *search) ing(k string) *search {
	s.Keyword = k

	key := url.QueryEscape(s.Keyword)
	resp, err := http.Get(query + key)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	s.doc, err = goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		panic(err)
	}

	return s
}

func (s *search) extract() {
	s.doc.Find("#feed-main-list .feed-block").Each(func(i int, selection *goquery.Selection) {
		e := new(entry).extract(selection)
		s.Entries = append(s.Entries, e)
	})
}
