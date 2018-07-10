package main

import (
	"flag"
	"fmt"
	"log"
	"u2pppw/crawler/crawler-distribute/rpcsupport"

	"u2pppw/crawler/crawler-distribute/worker"
)

//命令行参数
var port = flag.Int("port", 0, "the port for me to listen on")

//启动rpc服务器,用于支持解析请求和返回结果的服务
//go run worker.go --port=9000
//go run worker.go --port=9001
func main() {
	//解析port参数,用于指定监听的端口
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}
	log.Fatal(rpcsupport.ServeRpc(
		fmt.Sprintf(":%d", *port),
		worker.CrawlService{}))
}
