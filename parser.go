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
	case movieDOUQuery, movieIMDBQuery:
		attr, _ := dom.Find(selector).Attr("href")
		return attr, nil
	case movieSizeQuery:
		return nil, dom.Find(selector).Contents().Map(func(i int, selection *goquery.Selection) string {
			resultSet := selection.Find("div > div > table > tbody > tr").Map(func(i int, selection *goquery.Selection) string {
				return selection.Find("td.views-field.views-field-nothing.views-align-center").Text()
			})
			return strconv.Itoa(findBiggest(resultSet))
			//result := ""
			//for _, v := range resultSet {
			//	result = result + v
			//}
			//return result
		})
	default:
		return nil, nil
	}
}

func findBiggest(resultSet []string) (biggest int) {
	if resultSet != nil {
		biggest = transfer(resultSet[0])
		for _, v := range resultSet {
			if biggest <= transfer(v) {
				biggest = transfer(v)
			}
		}
	}
	return biggest
}

func transfer(size string) int {
	switch {
	case strings.Contains(size, "GB "):
		f, _ := strconv.ParseFloat(strings.TrimRight(size, "GB "), 64)
		return int(f * 1024 * 1024)
	case strings.Contains(size, "MB "):
		f, _ := strconv.ParseFloat(strings.TrimRight(size, "MB "), 64)
		return int(f * 1024)
	default:
		return 0
	}
}
