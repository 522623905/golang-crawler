package main

import (
	"u2pppw/crawler/crawler-queue/engine"
	"u2pppw/crawler/crawler-queue/persist"
	"u2pppw/crawler/crawler-queue/scheduler"
	"u2pppw/crawler/crawler-queue/zhenai/parser"
)

func main() {
	//用于与elasticSearch Saver通信的channel
	itemChan, err := persist.ItemSaver("dating_profile")
	if err != nil {
		panic(err)
	}

	//创建引擎
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      100,
		ItemChan:         itemChan, //把上述创建的用于与elasticSearch通信的channel传入进来
		RequestProcessor: engine.Worker,
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
