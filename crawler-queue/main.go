package main

import (
	"./engine"
	"./persist"
	"./scheduler"
	"./zhenai/parser"
)

func main() {
	itemChan, err := persist.ItemSaver("dating_profile")
	if err != nil {
		panic(err)
	}
	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 100,
		ItemChan:    itemChan,
	}
	e.Run(engine.Request{
		Url:       "http://www.zhenai.com/zhenghun",
		ParseFunc: parser.ParseCityList,
	})
	//	e.Run(engine.Request{
	//		Url:       "http://www.zhenai.com/zhenghun/nanning",
	//		ParseFunc: parser.ParseCity,
	//	})
}
