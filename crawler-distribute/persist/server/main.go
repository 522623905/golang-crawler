package main

import (
	"fmt"
	"u2pppw/crawler/crawler-distribute/config"
	"u2pppw/crawler/crawler-distribute/persist"
	"u2pppw/crawler/crawler-distribute/rpcsupport"

	"gopkg.in/olivere/elastic.v5"
)

func main() {
	err := serveRpc(
		fmt.Sprintf(":%d", config.ItemSaverPort),
		config.ElasticIndex)
	if err != nil {
		panic(err)
	}
}

//host: rpc服务的ip+端口
//index: elasticsearch index
func serveRpc(host, index string) error {
	client, err := elastic.NewClient(
		elastic.SetSniff(false))
	if err != nil {
		return err
	}

	//因为ItemSaverService提供的方法是指针接受者,因此加上&地址符合
	return rpcsupport.ServeRpc(host,
		&persist.ItemSaverService{
			Client: client,
			Index:  index,
		})
}
