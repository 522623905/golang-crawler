package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	//	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	//	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func main() {
	resp, err := http.Get("http://www.zhenai.com/zhenghun")
	//	resp, err := http.Get("http://www.zhenai.com/zhenghun/zibo")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:status code", resp.StatusCode)
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
	printCityList(all)
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

//从contens中，按照正则表达式提取信息并返回
func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

func printCityList(contents []byte) {
	//[^>]*:表示非>的任意字符
	re := regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`)
	//查找所有的匹配，能拆分上述正则表达式中()的内容
	matches := re.FindAllSubmatch(contents, -1)
	//打印出城市，ｕｒｌ
	for _, m := range matches {
		fmt.Printf("Ｃｉｔｙ:%s,URL:%s\n", m[2], m[1])
	}
	fmt.Printf("Matches found: %d\n", len(matches))
}
