package fetcher

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"u2pppw/crawler/crawler-distribute/config"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

//设置个定时器，防止爬虫太快，网站有保护措施
var rateLimiter = time.Tick(
	time.Second / config.Qps)

//提取url内容
func Fetch(url string) ([]byte, error) {
	<-rateLimiter //定时接收
	log.Printf("Fetching url %s", url)
	//要增加User-Agent头，否则爬虫的时候可能会返回403权限错误
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0")
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(reqest)

	//	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}
	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	defer resp.Body.Close() //注意放置的位置,否则可能出现invalid memory
	return ioutil.ReadAll(utf8Reader)
}

//自动识别编码格式
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
