package main

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

func parseQuery(dom *goquery.Document, selector string) (interface{}, []string) {
	switch selector {
	case pagesQuery:
		return dom.Find(selector).Children().Length(), nil
	case yearsQuery:
		result := dom.Find(selector).Contents().Map(func(i int, selection *goquery.Selection) string {
			attr, _ := selection.Attr("href")
			return attr
		})
		return nil, result
	case moviesQuery:
		result := dom.Find(selector).Contents().Children().Find("a").Map(func(i int, selection *goquery.Selection) string {
			attr, _ := selection.Attr("href")
			return attr
		})
		return nil, result
	case movieNameQuery:
		return dom.Find(selector).Text(), nil
	case movieDOUQuery, movieIMDBQuery, torrentQuery:
		attr, _ := dom.Find(selector).Attr("href")
		return attr, nil
	case movieSizeQuery:
		d := dom.Find(selector)
		resultSet := d.Contents().Map(func(i int, selection *goquery.Selection) string {
			return selection.Text()
		})
		biggest, index := findBiggest(resultSet)
		torrent := make([]string, 1)
		d.Parent().Find(torrentPreQuery).Contents().Each(func(i int, selection *goquery.Selection) {
			if i == index {
				attr, _ := selection.Parent().Attr("href")
				torrent[0] = attr
			}
		})
		return biggest, torrent
	default:
		return nil, nil
	}
}

func findBiggest(resultSet []string) (biggest, index int) {
	if resultSet != nil {
		biggest = transfer(resultSet[0])
		index = 0
		for i, v := range resultSet {
			if biggest <= transfer(v) {
				biggest = transfer(v)
				index = i
			}
		}
	}
	return biggest, index
}

func transfer(size string) int {
	switch {
	case strings.Contains(size, "MB "):
		f, _ := strconv.ParseFloat(strings.TrimRight(size, "MB "), 64)
		return int(f * 1024)
	default:
		f, _ := strconv.ParseFloat(strings.TrimRight(size, "GB "), 64)
		return int(f * 1024 * 1024)
	}
}
