package main

import (
	"testing"
	"time"

	"u2pppw/crawler/crawler-queue/engine"

	"u2pppw/crawler/crawler-distribute/config"
	"u2pppw/crawler/crawler-distribute/rpcsupport"
	"u2pppw/crawler/crawler-queue/model"
)

//get localhost:9200/test1/_search查看是否插入成功
func TestItemSaver(t *testing.T) {
	const host = ":1234"

	//start ItemSaverServer
	go serveRpc(host, "test1")
	time.Sleep(time.Second)

	//start ItemSaverClient
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		panic(err)
	}

	//call save
	item := engine.Item{
		Url:  "http://albnum.zhenai.com/u/108906739",
		Type: "zhenai",
		Id:   "108906739",
		Payload: model.Profile{
			Name:       "安静的雪",
			Age:        34,
			Height:     162,
			Weight:     57,
			Income:     "3001-5000元",
			Gender:     "女",
			Xinzuo:     "牡羊座",
			Marriage:   "离异",
			Education:  "大学本科",
			Occupation: "人事/行政",
			Hokou:      "山东菏泽",
			House:      "已购房",
			Car:        "未购车",
		},
	}

	result := ""
	err = client.Call(config.ItemSaverRpc, item, &result) //调用ItemSaverService.Save方法
	if err != nil || result != "ok" {
		t.Errorf("result: %s; err:%s", result, err)
	}
}
