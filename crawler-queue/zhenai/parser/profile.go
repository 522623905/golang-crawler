package parser

import (
	"fmt"
	//	"log"
	"os"
	"regexp"
	"strconv"

	"../../engine"
	"../../model"
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

//解析对应url下名字为name的人的相关信息
func ParseProfile(contents []byte, name string) engine.ParseResult {
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

	//	log.Printf("name:%s,age:%d,Gender:%s,Marriage:%s,height:%d,weight:%d,Income:%s,Xingzuo:%s,Education:%s,Occupation::%s,Hokou:%s,House:%s,Car:%s\n", name, age, profile.Gender, profile.Marriage, height, weight, profile.Income, profile.Xinzuo, profile.Education, profile.Occupation, profile.Hokou, profile.House, profile.Car)
	//	log.SetPrefix("")

	//记录信息到文件当中
	f, err := os.OpenFile("data.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		sWrite := fmt.Sprintf("name:%s,age:%d,Gender:%s,Marriage:%s,height:%d,weight:%d,Income:%s,Xingzuo:%s,Education:%s,Occupation::%s,Hokou:%s,House:%s,Car:%s\n", name, age, profile.Gender, profile.Marriage, height, weight, profile.Income, profile.Xinzuo, profile.Education, profile.Occupation, profile.Hokou, profile.House, profile.Car)
		f.Write([]byte(sWrite))
	}

	result := engine.ParseResult{
		Items: []interface{}{profile},
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
