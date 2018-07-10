package client

import (
	"net/rpc"
	"u2pppw/crawler/crawler-distribute/config"
	"u2pppw/crawler/crawler-distribute/worker"
	"u2pppw/crawler/crawler-queue/engine"
)

//创建处理器
func CreateProcessor(
	clients chan *rpc.Client) engine.Processor {

	return func(req engine.Request) (engine.ParseResult, error) {

		sReq := worker.SerializeRequest(req)

		var sResult worker.ParseResult
		//每次从rpc client池中取出一个client
		c := <-clients

		//rpc调用解析ｒｅｑｕｅｓｔ，返回ｒｅｓｕｌｔ
		err := c.Call(config.CrawlServiceRpc,
			sReq, &sResult)

		if err != nil {
			return engine.ParseResult{}, err
		}

		//反序列化成engine.ParseResult
		return worker.DeserializeResult(sResult), nil
	}
}
