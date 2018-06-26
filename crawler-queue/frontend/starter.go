package main

import (
	"net/http"

	"u2pppw/crawler/crawler-queue/frontend/controller"
)

/*
**用于显示elasticSearch数据库中的内容**
 */

func main() {
	//防止CSS等内容没有展示出来
	//因此使用http fileServer提供静态内容
	http.Handle("/", http.FileServer(
		http.Dir("./view/")))

	//显示静态图片
	http.HandleFunc("/images/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	//当访问到/search,则创建对象,解析模板
	//注意,SearchResultHandler使用了http.Handle必须要实现ServeHTTP()
	http.Handle("/search",
		controller.CreateSearchResultHandler(
			"./view/template.html"))

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
