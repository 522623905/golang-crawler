package parser

import (
	//	"log"

	"regexp"

	"u2pppw/crawler/crawler-distribute/config"
	"u2pppw/crawler/crawler-queue/engine"
)

//获取某城市内的人的（url+名字）表达式，如：
//<a href="http://album.zhenai.com/u/1188668898" target="_blank">交友征婚</a>
//第3页的信息，如："http://www.zhenai.com/zhenghun/shanghai/3"
var (
	//解析出客户的url
	profileRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`)
	//解析出city内容下一页的url,便于继续爬去此city的客户url
	cityUrlRe = regexp.MustCompile(`href="(http://www.zhenai.com/zhenghun/[^"]+)"`)
)

//在某城市中，解析出人的url和name
func ParseCity(contents []byte, url string) engine.ParseResult {
	matches := profileRe.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	//	limit := 10 //限制人的个数
	for _, m := range matches {
		//		result.Items = append(result.Items, "User "+string(m[2]))
		result.Requests = append(
			result.Requests,
			engine.Request{
				Url: string(m[1]),
				Parser: NewProfileParser(
					string(m[2])),
			})
		//		log.Printf("name:%s,url:%s\n", name, string(m[1]))
		//		limit--
		//		if limit == 0 {
		//			break
		//		}
	}

	//看到更多的城市信息
	matches = cityUrlRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests, engine.Request{
			Url: string(m[1]),
			Parser: engine.NewFuncParser(
				ParseCity, config.ParseCity),
		})
	}
	return result
}
