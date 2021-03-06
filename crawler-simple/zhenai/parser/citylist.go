package parser

import (
	"log"
	"regexp"

	"../../engine"
)

//城市列表的正则表达式. 加上()为的是正则匹配后能提取出该处的内容； [^>]* 表示非>的任意字符
const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

//在http://www.zhenai.com/zhenghun中解析出城市列表:城市+url地址
func ParseCityList(contents []byte) engine.ParseResult {
	re, _ := regexp.Compile(cityListRe)
	//查找所有的匹配，能拆分上述正则表达式中()的内容
	matches := re.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	//	limit := 1 //限制城市个数
	//遍历满足正则表达式的所有结果，并记录到result集中
	for _, m := range matches {
		result.Items = append(result.Items, "City "+string(m[2])) //m[2]为城市名
		//记录下每一个链接，每一个链接，实际是另一个要爬取的新请求（某城市对应的url）
		result.Requests = append(
			result.Requests,
			engine.Request{Url: string(m[1]), ParseFunc: ParseCity}) //m[1]为城市对应的url
		log.Printf("City: %s, Url:%s\n", m[2], m[1])

		//		limit--
		//		if limit == 0 {
		//			break
		//		}
	}
	return result
}
