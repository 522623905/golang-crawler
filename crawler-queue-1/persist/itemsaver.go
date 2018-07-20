package persist

import (
	"context"
	"errors"
	//	"fmt"
	"log"

	"u2pppw/crawler/crawler-queue-1/engine"

	"gopkg.in/olivere/elastic.v5"
)

/*
**该文件用于把爬取到的数据存储到elasticSearch中
 */

//index为elasticSearch的database name
func ItemSaver(index string) (chan engine.Item, error) {
	//没有指定url,默认去寻找localhost:9200
	//Sniff给客户端维护值权的状态,但是由于elastic运行在docker上,
	//我们外面看不到内网的状态,因此必须false
	client, err := elastic.NewClient(
		elastic.SetSniff(false))

	if err != nil {
		return nil, err
	}

	out := make(chan engine.Item)
	//注意这里的操作,开启goroutine来阻塞等待out channel数据来写入elasticSearch
	//而out channel被return回去给main,由main来进行数据的提取
	go func() {
		//记录数据个数
		itemCount := 0
		for {
			//item阻塞等待out channel有数据
			//见引擎的Run()函数
			item := <-out
			log.Printf("Item Saver:got item "+
				"#%d: %v", itemCount, item)
			itemCount++

			//把item保存到elasticSearch的index(database name)中
			err := Save(client, index, item)
			if err != nil {
				log.Printf("Item Saver: error "+
					"saving item %v:%v", item, err)
			}
		}
	}()
	return out, nil
}

//把item保存到elasticSearch的index(database name)中
func Save(client *elastic.Client, index string, item engine.Item) error {

	//如果ElasticSearch没有指定Type(table name),则出错
	if item.Type == "" {
		return errors.New("must supply Type")
	}

	//以json格式保存数据到ElasticSearch
	//client.Index()则为索引文档,即存数据的意思
	//index -> database ; Type -> table,并以json格式存储
	indexService := client.Index().
		Index(index).
		Type(item.Type).
		BodyJson(item)
	if item.Id != "" {
		indexService.Id(item.Id)
	}

	//执行完后,则存储进elasticSearch
	_, err := indexService.
		Do(context.Background())

	if err != nil {
		return err
	}
	//%+v会把结构体的字段名也打出来
	//	fmt.Printf("%+v", resp)

	return nil
}
