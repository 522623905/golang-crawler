package model

//import "../../engine"

//搜索elasticSearch后返回的结果格式
type SearchResult struct {
	Hits     int64  //命中,即总共找到多少个
	Start    int    //从第几个数据开始
	Query    string //获取参数内容
	PrevFrom int    //上一页
	NextFrom int    //下一页
	//	Items    []engine.Item //具体的内容
	Items []interface{}
}
