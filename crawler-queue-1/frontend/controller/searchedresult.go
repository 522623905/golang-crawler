package controller

import (
	"context"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"u2pppw/crawler/crawler-queue-1/frontend/model"

	"u2pppw/crawler/crawler-queue-1/engine"

	"u2pppw/crawler/crawler-queue-1/frontend/view"

	"gopkg.in/olivere/elastic.v5"
)

type SearchResultHandler struct {
	view   view.SearchResultView //解析模板文件后的html template对象
	client *elastic.Client       //与elasticSearch通信的对象,从client中把数据给view
}

//用于创建对象,传入的参数为模板文件名
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

//http.Handle必须要实现ServeHTTP()接口
//localhost:9200/search?q=男 已购房&from=20
//from代表从第20条数据开始,即第二页开始(默认一次显示10行数据(0-9行))
func (h SearchResultHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//拿到要查询的参数,即q后面的value,并去掉空格
	q := strings.TrimSpace(req.FormValue("q"))
	//拿到from后面的value
	from, err := strconv.Atoi(req.FormValue("from"))
	//如果from后面的参数,若填错,则忽略为0
	if err != nil {
		from = 0
	}

	//根据q,from参数返回搜索的结果,即一页的内容page(elasticSearch默认10条数据)
	page, err := h.getSearchResult(q, from)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusBadRequest)
	}

	//把结果page写入html模板中并显示
	err = h.view.Render(w, page)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusBadRequest)
	}
}

//搜索中输入如:
//	男 已购房 已购车 Age:(<30) Height:(>180)
//q:要查询的内容  from:页码
//从elasticSearch获取相关内容
func (h SearchResultHandler) getSearchResult(q string, from int) (model.SearchResult, error) {
	var result model.SearchResult
	result.Query = q //获取url中"q="后面的内容

	//elasticSearch默认一次显示10行数据(0-9行),使用From来指定从第几行开始显示
	resp, err := h.client.
		Search("dating_profile").
		Query(elastic.NewQueryStringQuery(
			rewriteQueryString(q))).
		From(from).
		Do(context.Background())

	if err != nil {
		return result, err
	}

	result.Hits = resp.TotalHits() //命中个数
	result.Start = from            //第几条数据开始
	//使用Each函数,遍历匹配的数据
	//使用反射机制,解析成engine.Item{}格式
	result.Items = resp.Each(
		reflect.TypeOf(engine.Item{}))
	result.PrevFrom = result.Start - len(result.Items) //上一页
	result.NextFrom = result.Start + len(result.Items) //下一页
	return result, nil
}

//字符串重写,前面加上Payload. ,
//因为数据存入elasticSearch时数据是以Payload存储的
//否则查询搜索的时候要手动加入 Payload. ,比较麻烦
func rewriteQueryString(q string) string {
	re := regexp.MustCompile(`([A-Z][a-z]*):`)
	//将上述匹配的表达式找到的字符替换成Payload.$1($1即上面找到的字符)
	return re.ReplaceAllString(q, "Payload.$1:")
}
