package engine

type ParseFunc func(contents []byte, url string) ParseResult

//解析的接口
type Parser interface {
	Parse(contents []byte, url string) ParseResult //解析url的函数
	Serialize() (name string, args interface{})    //用于rpc的封装,序列化后返回函数名和传递给函数的参数
}

//解析请求struct
type Request struct {
	Url    string //要解析的url
	Parser Parser //该url对应的解析接口
}

type SerializedParser struct {
	Name string      //函数名
	Args interface{} //参数
}

//{"ParserCityList",nil},{"ProfileParser",userName}

//解析后返回的结果集struct
type ParseResult struct {
	Requests []Request //要请求的内容
	Items    []Item    //结果集的具体内容
}

type Item struct {
	Url     string      //人物的url
	Id      string      //人物的ID,去重时用的,也用作ElasticSearch的Id
	Type    string      //ElasticSearch的table name
	Payload interface{} //具体数据
}

//NilParser
type NilParser struct{}

//实现Parser接口
func (NilParser) Parse(
	_ []byte, _ string) ParseResult {
	return ParseResult{}
}

//实现Parser接口
func (NilParser) Serialize() (
	name string, args interface{}) {
	return "NilParser", nil
}

//FuncParser
type FuncParser struct {
	parser ParseFunc //parse函数
	name   string    //parse函数名
}

//实现Parser接口
func (f *FuncParser) Parse(contents []byte, url string) ParseResult {
	return f.parser(contents, url)
}

//实现Parser接口
func (f *FuncParser) Serialize() (
	name string, args interface{}) {
	return f.name, nil
}

//new一个FuncParser结构
func NewFuncParser(p ParseFunc, name string) *FuncParser {
	return &FuncParser{
		parser: p,
		name:   name,
	}
}
