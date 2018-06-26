package client

import (
	"log"

	"u2pppw/crawler/crawler-distribute/config"
	"u2pppw/crawler/crawler-distribute/rpcsupport"
	"u2pppw/crawler/crawler-queue/engine"
)

//index为elasticSearch的database name
func ItemSaver(host string) (chan engine.Item, error) {
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		return nil, err
	}

	out := make(chan engine.Item)
	//注意这里的操作,开启goroutine来阻塞等待out channel数据来写入elasticSearch
	//而out channel被return回去给main,由main来进行数据的提取
	go func() {
		itemCount := 0
		for {
			//item阻塞等待out channel有数据
			//见引擎的Run()函数
			item := <-out
			log.Printf("Item Saver:got item "+
				"#%d: %v", itemCount, item)
			itemCount++

			//call rpc把item保存到elasticSearch的index(database name)中
			result := ""
			err := client.Call(config.ItemSaverRpc, item, &result)
			if err != nil {
				log.Printf("Item Saver: error "+
					"saving item %v:%v", item, err)
			}
		}
	}()
	return out, nil
}
