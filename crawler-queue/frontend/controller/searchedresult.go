package controller

import (
	"context"
	//	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"../model"

	"../../engine"
	"../view"
	"gopkg.in/olivere/elastic.v5"
)

type SearchResultHandler struct {
	view   view.SearchResultView
	client *elastic.Client //从client中把数据给view
}

func CreateSearchResultHandler(
	template string) SearchResultHandler {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	return SearchResultHandler{
		view:   view.CreateSearchResultView(template),
		client: client,
	}
}

//localhost:9200/search?q=男 已购房&from=20
//from代表页
func (h SearchResultHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//拿到q后面的value,并去掉空格
	q := strings.TrimSpace(req.FormValue("q"))
	//拿到from后面的value
	from, err := strconv.Atoi(req.FormValue("from"))
	//如果from后面的参数填错,则忽略为0
	if err != nil {
		from = 0
	}
	//	fmt.Fprintf(w, "q=%s,from=%d", q, from)

	page, err := h.getSearchResult(q, from)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusBadRequest)
	}
	err = h.view.Render(w, page)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusBadRequest)
	}
}

//搜索中输入如:
//	男 已购房 已购车 Age:(<30) Height:(>180)
func (h SearchResultHandler) getSearchResult(q string, from int) (model.SearchResult, error) {
	var result model.SearchResult
	result.Query = q //获取url中"q="后面的内容
	resp, err := h.client.
		Search("dating_profile").
		Query(elastic.NewQueryStringQuery(
			rewriteQueryString(q))).
		From(from).
		Do(context.Background())

	if err != nil {
		return result, err
	}

	result.Hits = resp.TotalHits()
	result.Start = from
	result.Items = resp.Each(
		reflect.TypeOf(engine.Item{}))
	result.PrevFrom = result.Start - len(result.Items)
	result.NextFrom = result.Start + len(result.Items)
	return result, nil
}

//字符串重写
func rewriteQueryString(q string) string {
	re := regexp.MustCompile(`([A-Z][a-z]*):`)
	//将上述匹配的表达式找到的字符替换成Payload.$1($1即上面找到的字符)
	return re.ReplaceAllString(q, "Payload.$1:")
}
