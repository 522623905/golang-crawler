package engine

import (
	"fmt"
	"log"
	"strings"

	"../fetcher"
)

type SimpleEngine struct{}

//运行
func (e SimpleEngine) Run(seeds ...Request) {
	//把传进来的请求，添加到Request slice中
	var requests []Request
	for _, r := range seeds {
		requests = append(requests, r)
	}

	//循环取出每个请求，并work
	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]

		//解析请求r
		parseResult, err := worker(r)
		if err != nil {
			continue
		}

		//解析后的返回结果集又有新的请求，则继续添加到请求队列
		requests = append(requests, parseResult.Requests...)
		for _, item := range parseResult.Items {
			log.Printf("Got item %v", item) //%v表示数据的默认格式
		}
	}
}

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
	return r.ParseFunc(body), nil
}
