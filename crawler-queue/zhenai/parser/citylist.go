package parser

import (
	//	"log"
	"regexp"

	"u2pppw/crawler/crawler-distribute/config"
	"u2pppw/crawler/crawler-queue/engine"
)

//城市列表的正则表达式. 加上()为的是正则匹配后能提取出该处的内容； [^>]* 表示非>的任意字符
//如:<a href="http://www.zhenai.com/zhenghun/guangzhou" class="">广州</a>
const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

//在http://www.zhenai.com/zhenghun中解析出城市列表:城市+url地址
func ParseCityList(contents []byte, url string) engine.ParseResult {
	re, _ := regexp.Compile(cityListRe)
	//查找所有的匹配，能拆分上述正则表达式中()的内容
	matches := re.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	//	limit := 1 //限制城市个数
	for _, m := range matches {
		//		result.Items = append(result.Items, "City "+string(m[2])) //m[2]为城市名
		result.Requests = append(
			result.Requests, engine.Request{
				Url: string(m[1]),
				Parser: engine.NewFuncParser(
					ParseCity, config.ParseCity),
			}) //m[1]为城市对应的url
		//		log.Printf("City: %s, Url:%s\n", m[2], m[1])

		//		limit--
		//		if limit == 0 {
		//			break
		//		}
	}
	return result
}
