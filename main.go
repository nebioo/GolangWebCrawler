package main

import (
	"fmt"
	"net/http"
	"os"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strings"
	log "github.com/llimllib/loglevel"
	"golang.org/x/text/unicode/rangetable"
)

type link struct{
	url string
	text string
	depth int
}

type HttpError struct {
	original string
}

func LinkReader(resp *http.Response, depth int) []link {
	page := html.NewTokenizer(resp.Body)
	links := []Link{}

	var start *html.Token
	var text string

	for {
		_ = page.Next()
		token := page.Token()
		if token.Type == html.ErrorToken {
			break
		}

		if start != nil && token.Type == html.TextToken {
			text = fmt.Sprintf("%s%s", text, token.Data)
		}

		if token.DataAtom == atom.A {
			switch token.Type {
			case html.StartTagToken:
				if len(token.Attr) > 0 {
					start = &token
				}
			case html.EndTagToken:
				if start == nil {
					log.Warnf("Link End found without Start: &s",text)
					continue
				}
				link := NewLink(*start, text, depth)
				if link.Valid() {
					links = append(links, link)
					log.Debugf("Link Found %v", link)
				}

				start = nil
				text = ""

			}
		}
	}

	log.Debug(links)
	return links
}

func NewLink(tag html.Token, text string, depth int) Link  {
	link := Link{text: strings.TrimSpace(text), depth: depth}

	for i:= range tag.Attr {
		if tag.Attr[i].Key == "href" {
			link.url = strings.TrimSpace(tag.Attr[i].Val)
		}
	}
}