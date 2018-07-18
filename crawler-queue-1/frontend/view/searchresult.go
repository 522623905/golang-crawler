package view

import (
	"html/template"
	"io"

	"u2pppw/crawler/crawler-queue-1/frontend/model"
)

//模板html的template
type SearchResultView struct {
	template *template.Template
}

//解析filename的html模板文件
func CreateSearchResultView(filename string) SearchResultView {
	return SearchResultView{
		template: template.Must(template.ParseFiles(filename)), //加载模板文件
	}
}

//使用模板对象执行模板，数据对象data注入w中
func (s SearchResultView) Render(w io.Writer, data model.SearchResult) error {
	return s.template.Execute(w, data)
}
