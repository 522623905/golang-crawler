package engine

import (
	"fmt"
	"log"
	"strings"

	"u2pppw/crawler/crawler-queue/fetcher"
)

//解析请求，返回结果
func worker(r Request) (ParseResult, error) {
	//	log.Printf("Fetching %s", r.Url)
	if strings.Contains(r.Url, "qishi") {
		return ParseResult{}, fmt.Errorf("parse %s is wrong,so continue", r.Url)
	}
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("Fetcher: error fetching url %s: %v", r.Url, err)
		return ParseResult{}, err
	}
	return r.Parser.Parse(body, r.Url), nil
}
