package main

import (
	"flag"
	"fmt"
	"u2pppw/crawler/crawler-distribute/config"
	"u2pppw/crawler/crawler-distribute/persist"
	"u2pppw/crawler/crawler-distribute/rpcsupport"

	"gopkg.in/olivere/elastic.v5"
)

var port = flag.Int("port", 0, "the port for me to listen on")

//启动rpc服务,用于支持itemsaver服务(将item存储到ElasticSearch)
//go run itemsaver.go --port=1234
func main() {
	//解析port参数,用于指定监听的端口
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}
	err := serveRpc(
		fmt.Sprintf(":%d", *port),
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

	//因为ItemSaverService提供的方法是指针接受者,因此加上&地址符号
	return rpcsupport.ServeRpc(host,
		&persist.ItemSaverService{
			Client: client,
			Index:  index,
		})
}
