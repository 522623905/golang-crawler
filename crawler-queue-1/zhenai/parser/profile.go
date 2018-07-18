package parser

import (
	"regexp"
	"strconv"

	"u2pppw/crawler/crawler-queue-1/engine"

	"u2pppw/crawler/crawler-queue-1/model"
)

//个人信息的相关正则表达式,其中[\d]表示匹配数字
var ageRe = regexp.MustCompile(`<td><span class="label">年龄：</span>([\d]+)岁</td>`)
var heightRe = regexp.MustCompile(`<td><span class="label">身高：</span>([\d]+)CM</td>`)
var weightRe = regexp.MustCompile(`<td><span class="label">体重：</span><span field="">([\d]+)KG</span></td>`)
var incomeRe = regexp.MustCompile(`<td><span class="label">月收入：</span>([^<]+)</td>`)
var genderRe = regexp.MustCompile(`<td><span class="label">性别：</span><span field="">([^<]+)</span></td>`)
var xinzuoRe = regexp.MustCompile(`<td><span class="label">星座：</span><span field="">([^<]+)</span></td>`)
var marriageRe = regexp.MustCompile(`<td><span class="label">婚况：</span>([^<]+)</td>`)
var educationRe = regexp.MustCompile(`<td><span class="label">学历：</span>([^<]+)</td>`)
var occupationRe = regexp.MustCompile(`<td><span class="label">职业： </span>([^<]+)</td>`)
var hokouRe = regexp.MustCompile(`<td><span class="label">籍贯：</span>([^<]+)</td>`)
var houseRe = regexp.MustCompile(`<td><span class="label">住房条件：</span><span field="">([^<]+)</span></td>`)
var carRe = regexp.MustCompile(`<td><span class="label">是否购车：</span><span field="">([^<]+)</span></td>`)

//页面出现的"猜你喜欢"的客户
var guessRe = regexp.MustCompile(`<a class="exp-user-name"[^>]*href="(http://albnum.zhenai.com/u/[\d]+)">([^<])</a>`)
var idUrlRe = regexp.MustCompile(`http://albnum.zhenai.com/u/([\d]+)`)

//解析对应url下名字为name的人的相关信息
func ParseProfile(contents []byte, url string, name string) engine.ParseResult {
	profile := model.Profile{}
	profile.Name = name
	age, err := strconv.Atoi(extractString(contents, ageRe))
	if err == nil {
		profile.Age = age
	}

	profile.Marriage = extractString(contents, marriageRe)

	height, err := strconv.Atoi(extractString(contents, heightRe))
	if err == nil {
		profile.Height = height
	}

	weight, err := strconv.Atoi(extractString(contents, weightRe))
	if err == nil {
		profile.Weight = weight
	}

	profile.Income = extractString(contents, incomeRe)
	profile.Gender = extractString(contents, genderRe)
	profile.Xinzuo = extractString(contents, xinzuoRe)
	profile.Marriage = extractString(contents, marriageRe)
	profile.Education = extractString(contents, educationRe)
	profile.Occupation = extractString(contents, occupationRe)
	profile.Hokou = extractString(contents, hokouRe)
	profile.House = extractString(contents, houseRe)
	profile.Car = extractString(contents, carRe)

	result := engine.ParseResult{
		Items: []engine.Item{
			{
				Url: url,
				//用于存储到elasticSearch的相关内容
				Type:    "zhenai",
				Id:      extractString([]byte(url), idUrlRe),
				Payload: profile,
			},
		},
	}

	matches := guessRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests,
			engine.Request{
				Url:       string(m[1]),
				ParseFunc: ProfileParser(string(m[2])),
			})
	}

	return result
}

//从contens中，按照正则表达式提取信息并返回
func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents) //只要返回的第一个匹配的正则表达式
	//match[0]是匹配的正则出来的所有字符串，match[1]是正则表达式中第一个()的内容
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

func ProfileParser(name string) engine.ParseFunc {
	return func(
		c []byte, url string) engine.ParseResult {
		return ParseProfile(c, url, name)
	}
}
