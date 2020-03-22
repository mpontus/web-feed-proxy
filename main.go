package main

import (
	"net/http"
	"net/url"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"github.com/gorilla/feeds"
)

type Options struct {
	// Source URL
	url *url.URL
	// Item Selector
	is *xpath.Expr
	// Title Selector
	ts *xpath.Expr
	// Link Selector
	ls *xpath.Expr
}

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		o, err := parseOptions(req.URL.Query())

		if err != nil {
			return nil, err
		}

		resp, err := http.Get(o.url.String())

		doc, err := htmlquery.Parse(resp.Body)

		feed := &feeds.Feed{}
		for _, node := range htmlquery.Find(doc, is) {
			item := models.FeedItem{
				Title: htmlquery.InnerText(htmlquery.FindOne(node, ts)),
				Link:  htmlquery.InnerText(htmlquery.FindOne(node, ts)),
			}

			feed.Items = append(feed.Items, item)
		}

		res.Write(feed.ToRss())
	})
}

func handleRequest(req *http.Request) ([]byte, error) {
	o, err := parseOptions(req.URL.Query())

	if err != nil {
		return nil, err
	}

	resp, err := http.Get(o.url.String())

	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(resp.Body)

	if doc != nil {
		return nil, err
	}

	n := htmlquery.CreateXPathNavigator(doc)

	o.is.Select(n)

	feed, err = assembleFeed(func(item) {
		for _, node := range htmlquery.Find(doc, is) {
			item()
			item := models.FeedItem{
				Title: htmlquery.InnerText(htmlquery.FindOne(node, ts)),
				Link:  htmlquery.InnerText(htmlquery.FindOne(node, ts)),
			}

			feed.Items = append(feed.Items, item)
		}
	})
}

func parseOptions(v url.Values) (*Options, error) {
	u, err := url.Parse(v.Get("url"))

	if err != nil {
		return nil, err
	}

	is, err := xpath.Compile(v.Get("is"))

	if err != nil {
		return nil, err
	}

	ts, err := xpath.Compile(v.Get("ts"))

	if err != nil {
		return nil, err
	}

	ls, err := xpath.Compile(v.Get("ls"))

	if err != nil {
		return nil, err
	}

	return &Options{u, is, ts, ls}, nil
}
