package main

import (
	"fmt"
	"u2pppw/crawler/crawler-distribute/config"
	"u2pppw/crawler/crawler-distribute/persist/client"
	"u2pppw/crawler/crawler-queue/engine"
	"u2pppw/crawler/crawler-queue/scheduler"
	"u2pppw/crawler/crawler-queue/zhenai/parser"
)

func main() {
	//通过rpc用于与elasticSearch通信
	itemChan, err := client.ItemSaver(
		fmt.Sprintf(":%d", config.ItemSaverPort))
	if err != nil {
		panic(err)
	}

	//创建引擎
	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 100,
		ItemChan:    itemChan, //用于与elasticSearch通信的channel
	}

	//引擎启动,并到该URl中爬取数据
	e.Run(engine.Request{
		Url: "http://www.zhenai.com/zhenghun",
		Parser: engine.NewFuncParser(
			parser.ParseCityList, "ParseCityList"),
	})
	//	e.Run(engine.Request{
	//		Url:       "http://www.zhenai.com/zhenghun/nanning",
	//		ParseFunc: parser.ParseCity,
	//	})
}
