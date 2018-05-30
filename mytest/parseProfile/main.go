package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	//	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	//	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

//人物的个人信息数据
type Profile struct {
	Name       string
	Gender     string
	Age        int
	Height     int
	Weight     int
	Income     string
	Marriage   string
	Education  string
	Occupation string
	Hokou      string
	Xinzuo     string
	House      string
	Car        string
}

//个人信息的相关正则表达式,其中[\d]表示匹配数字
var name = regexp.MustCompile(`<p>[^<]+<a class="name fs24">([^<]+)</a>`)
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

//Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0
func main() {
	//增加头，否则爬虫可能会出现403权限不足的错误
	client := &http.Client{}
	url := "http://album.zhenai.com/u/1463502270"
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0")
	if err != nil {
		panic(err)
	}
	resp, _ := client.Do(reqest)
	//	resp, err := http.Get("http://album.zhenai.com/u/1147788337") //http://album.zhenai.com/u/110294731
	//	if err != nil {
	//		panic(err)
	//		//		resp, err = http.Get("http://album.zhenai.com/u/1147788337")

	//	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:status code", resp.StatusCode)
		//		resp, err = http.Get("http://album.zhenai.com/u/1147788337")
		return
	}

	//	2.要把GBK(传回的html中meta charset显示)转换成UTF-8格式(golang默认格式)，否则乱码
	//			utf8Reader := transform.NewReader(
	//				resp.Body, simplifiedchinese.GBK.NewDecoder())
	//	1.直接读取body：中文乱码 all, err := ioutil.ReadAll(resp.Body)

	//3. 最好做法是，自动识别编码格式，然后转换成utf-8
	e := determineEncoding(resp.Body)
	utf8Reader := transform.NewReader(
		resp.Body, e.NewDecoder())
	all, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		panic(err)
	}

	printProfile(all)
}

//自动识别html返回的格式
func determineEncoding(r io.Reader) encoding.Encoding {
	//不可直接传入r给DetermineEncoding()，因为读取1024字节后，
	//	就丢失了，因此先返回1024字节出来给一个bytes,并不移动r的指针
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		panic(err)
	}
	//它会收取一开始的1024 bytes来判断编码的格式
	e, _, _ := charset.DetermineEncoding(
		bytes, "")
	return e
}

//解析对应url下名字为name的人的相关信息
func printProfile(contents []byte) {
	profile := Profile{}
	//	profile.Name = extractString(contents, name)
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

	fmt.Println("Person:\n age: ", age, " Marriage: ", profile.Marriage,
		"  height: ", profile.Height, " Income: ", profile.Income, " Education: ", profile.Education,
		" House: ", profile.House, " Car: ", profile.Car)

}

//从contens中，按照正则表达式提取信息并返回
func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}
