package persist

import (
	"context"
	"errors"
	//	"fmt"
	"log"

	"../engine"

	"gopkg.in/olivere/elastic.v5"
)

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
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver:got item "+
				"#%d: %v", itemCount, item)
			itemCount++

			err := save(client, index, item)
			if err != nil {
				log.Printf("Item Saver: error "+
					"saving item %v:%v", item, err)
			}
		}
	}()
	return out, nil
}

func save(client *elastic.Client, index string, item engine.Item) error {

	if item.Type == "" {
		return errors.New("must supply Type")
	}

	//开始存数据
	//index -> database ; Type -> table
	indexService := client.Index().
		Index(index).
		Type(item.Type).
		BodyJson(item)
	if item.Id != "" {
		indexService.Id(item.Id)
	}

	_, err := indexService.
		Do(context.Background())

	if err != nil {
		return err
	}
	//%+v会把结构体的字段名也打出来
	//	fmt.Printf("%+v", resp)

	return nil
}
