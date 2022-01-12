package main

import (
	"github.com/go-resty/resty/v2"
	"time"
)

func fetchUrl(url string) ([]byte, error) {
	<-time.Tick(1000 * time.Millisecond)
	client := resty.New()
	get, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	return get.Body(), err
}
