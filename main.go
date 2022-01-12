package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"sync"
)

const (
	root  string = "https://www.mini4k.com"
	oscar int    = 24799
	//golden      int    = 24867
	//cannes      int    = 24800
	//venice      int    = 24801
	//berlin      int    = 24802
	pagesQuery      string = "#block-white-content > div > div > nav > div"
	yearsQuery      string = "#block-white-content > div > div > div > ul > li"
	moviesQuery     string = "#block-white-content > article > div.node-content > div > div > div > div > div > ul > li"
	movieDOUQuery   string = "#block-white-content > article > div > div > div > div.eleven.wide.column > div.node-overview.clearfix > div.node-detail > div.clearfix > div.title-right > div.out-links > a.douban"
	movieIMDBQuery  string = "#block-white-content > article > div > div > div > div.eleven.wide.column > div.node-overview.clearfix > div.node-detail > div.clearfix > div.title-right > div.out-links > a.imdb"
	movieNameQuery  string = "#block-white-content > article > div > div > div > div.eleven.wide.column > div.node-overview.clearfix > div.node-detail > div.clearfix > div.node-title > h1 > span"
	movieSizeQuery  string = "#block-white-content > article > div > div > div > div.eleven.wide.column > div.reference-torrent > div > div > div > table > tbody > tr > td.views-field.views-field-nothing.views-align-center"
	torrentPreQuery string = "td.views-field.views-field-title.views-align-left > a"
	torrentQuery    string = "#block-white-content > article > div > div > div > div.eleven.wide.column > div.download-buttons > a.ui.primary.button"
)

type Movie struct {
	douBan  string
	imdb    string
	torrent string
	name    string
	size    int
}

var (
	parsedYears      []string
	parsedMovies     = make(chan string, 16)
	serializedMovies = make(chan *Movie, 16)
)

func init() {
	InitFlag()
	flag.Parse()
}

func main() {
	connection, err := NewDatabaseConnection()
	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {
			log.Println(err)
		}
	}(connection)
	if err != nil {
		log.Fatalln(err)
	}

	prepare, err := connection.Prepare("insert into oscar values (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatalln(err)
	}

	var pages []*goquery.Document
	dom := getPage(fmt.Sprintf("%s/awards?term=%s&page=0", root, strconv.Itoa(oscar)))
	pages = append(pages, dom)

	length, _ := parseQuery(dom, pagesQuery)
	for i := 1; i < length.(int)-2; i++ {
		dom := getPage(fmt.Sprintf("%s/awards?term=%s&page=%d", root, strconv.Itoa(oscar), i))
		pages = append(pages, dom)
	}
	for _, v := range pages {
		_, years := parseQuery(v, yearsQuery)
		for _, v := range years {
			parsedYears = append(parsedYears, v)
		}
	}
	var wg sync.WaitGroup
	for _, v := range parsedYears {
		wg.Add(1)
		go func(v string, wg *sync.WaitGroup) {
			dom := getPage(fmt.Sprintf("%s%s", root, v))
			_, movies := parseQuery(dom, moviesQuery)
			for _, v := range movies {
				parsedMovies <- fmt.Sprintf("%s%s", root, v)
			}
			wg.Done()
		}(v, &wg)
	}
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			for {
				select {
				case url := <-parsedMovies:
					m := new(Movie)
					dom := getPage(url)
					douBan, _ := parseQuery(dom, movieDOUQuery)
					imdb, _ := parseQuery(dom, movieIMDBQuery)
					name, _ := parseQuery(dom, movieNameQuery)
					size, torrent := parseQuery(dom, movieSizeQuery)
					dom = getPage(root + torrent[0])
					tmp, _ := parseQuery(dom, torrentQuery)
					magnet := strings.Split(tmp.(string), "&")
					m.douBan, m.imdb, m.torrent, m.name, m.size = douBan.(string), imdb.(string), magnet[0], name.(string), size.(int)
					serializedMovies <- m
				}
			}
		}()
	}
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			for {
				select {
				case m := <-serializedMovies:
					fmt.Printf("%+v\n", m)
					_, err := prepare.Exec(m.douBan, m.imdb, m.torrent, m.name, m.size)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}()
	}
	wg.Wait()
}

func getPage(url string) *goquery.Document {
	page, err := fetchUrl(url)
	if err != nil {
		log.Println(err)
	}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(page)))
	if err != nil {
		log.Println(err)
	}
	return dom
}
