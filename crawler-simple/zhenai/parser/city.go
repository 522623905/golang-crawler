package parser

import (
	"log"

	"regexp"

	"../../engine"
)

//获取某城市内的人的（url+名字）表达式，如：
//<a href="http://album.zhenai.com/u/1188668898" target="_blank">交友征婚</a>
const cityRe = `<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`

//在某城市中，解析出人的url和name
func ParseCity(contents []byte) engine.ParseResult {
	re, _ := regexp.Compile(cityRe)
	matches := re.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	//	limit := 10 //限制人的个数
	for _, m := range matches {
		name := string(m[2])
		//记录用户名字
		result.Items = append(result.Items, "User "+string(m[2]))
		//每一个用户对应的链接，又是要爬取的内容
		result.Requests = append(
			result.Requests,
			engine.Request{
				Url: string(m[1]),
				ParseFunc: func(bytes []byte) engine.ParseResult {
					return ParseProfile(bytes, name)
				},
			})
		log.Printf("name:%s,url:%s\n", name, string(m[1]))
		//		limit--
		//		if limit == 0 {
		//			break
		//		}
	}
	return result
}
