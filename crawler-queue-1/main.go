package main

import (
	"u2pppw/crawler/crawler-queue-1/engine"
	"u2pppw/crawler/crawler-queue-1/persist"
	"u2pppw/crawler/crawler-queue-1/scheduler"
	"u2pppw/crawler/crawler-queue-1/zhenai/parser"
)

func main() {
	//用于与elasticSearch通信的channel
	itemChan, err := persist.ItemSaver("dating_profile")
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
		Url:       "http://www.zhenai.com/zhenghun",
		ParseFunc: parser.ParseCityList,
	})
	//	e.Run(engine.Request{
	//		Url:       "http://www.zhenai.com/zhenghun/nanning",
	//		ParseFunc: parser.ParseCity,
	//	})
}
