package engine

//解析请求struct
type Request struct {
	Url       string                   //要解析的url
	ParseFunc func([]byte) ParseResult //该url对应的解析函数
}

//解析后返回的结果集struct
type ParseResult struct {
	Requests []Request     //要请求的内容
	Items    []interface{} //结果集的具体内容
}

func NilParser([]byte) ParseResult {
	return ParseResult{}
}
