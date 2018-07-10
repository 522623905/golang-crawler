package main

import (
	"flag"
	"log"
	"net/rpc"
	"strings"
	"u2pppw/crawler/crawler-distribute/config"
	"u2pppw/crawler/crawler-distribute/persist/client"
	"u2pppw/crawler/crawler-distribute/rpcsupport"
	"u2pppw/crawler/crawler-queue/engine"
	"u2pppw/crawler/crawler-queue/scheduler"
	"u2pppw/crawler/crawler-queue/zhenai/parser"

	client2 "u2pppw/crawler/crawler-distribute/worker/client"
)

var (
	itemSaverHost = flag.String(
		"itemsaver_host", "", "itemsaver host")
	workerHosts = flag.String(
		"worker_hosts", "", "worker hosts (comma separated)") //逗号分割
)

//go run main.go -itemsaver_host=":1234" --worker_hosts=":9000,:9001"
func main() {
	flag.Parse()
	//连接itemSaverHost的rpc服务器,通过rpc用于与elasticSearch通信
	itemChan, err := client.ItemSaver(
		*itemSaverHost)
	if err != nil {
		panic(err)
	}
	log.Printf("Connect to %s", *itemSaverHost)

	//创建连接池
	pool := createClientPool(
		strings.Split(*workerHosts, ",")) //逗号分割出来

	//创建处理器
	processor := client2.CreateProcessor(pool)

	//创建引擎
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      100,
		ItemChan:         itemChan, //用于与elasticSearch通信的channel
		RequestProcessor: processor,
	}

	//引擎启动,并到该URl中爬取数据
	e.Run(engine.Request{
		Url: "http://www.zhenai.com/zhenghun",
		Parser: engine.NewFuncParser(
			parser.ParseCityList,
			config.ParseCityList),
	})
}

//创建rpc client连接池
func createClientPool(hosts []string) chan *rpc.Client {
	var clients []*rpc.Client
	for _, h := range hosts {
		client, err := rpcsupport.NewClient(h)
		if err == nil {
			clients = append(clients, client)
			log.Printf("Connected to %s", h)
		} else {
			log.Printf("error connecting to %s: %v", h, err)
		}
	}

	//创建一个channel,负责分发池中的client
	out := make(chan *rpc.Client)
	go func() {
		//该层for循环负责一直分发,否则发完了就完了
		for {
			//该层for循环负责轮流分发
			for _, client := range clients {
				out <- client
			}
		}
	}()
	return out
}
