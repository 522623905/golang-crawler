package client

import (
	"log"

	"u2pppw/crawler/crawler-distribute/config"
	"u2pppw/crawler/crawler-distribute/rpcsupport"
	"u2pppw/crawler/crawler-queue/engine"
)

//host:rpc服务器地址
//goroutine阻塞等待数据到来然后调用RPC存储到ElasticSearch
//返回channel给main,供main提供数据给阻塞的goroutine
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

			//call rpc把item保存到elasticSearch
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
