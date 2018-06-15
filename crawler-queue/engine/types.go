package engine

type ParseFunc func(contents []byte, url string) ParseResult

//解析请求struct
type Request struct {
	Url       string    //要解析的url
	ParseFunc ParseFunc //该url对应的解析函数
}

//解析后返回的结果集struct
type ParseResult struct {
	Requests []Request //要请求的内容
	Items    []Item    //结果集的具体内容
}

type Item struct {
	Url     string //人物的url
	Id      string //人物的ID,去重时用的
	Type    string
	Payload interface{}
}

func NilParser([]byte) ParseResult {
	return ParseResult{}
}
