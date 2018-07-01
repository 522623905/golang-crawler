package client

import (
	"net/rpc"
	"u2pppw/crawler/crawler-distribute/config"
	"u2pppw/crawler/crawler-distribute/worker"
	"u2pppw/crawler/crawler-queue/engine"
)

func CreateProcessor(
	clients chan *rpc.Client) engine.Processor {

	return func(req engine.Request) (engine.ParseResult, error) {

		sReq := worker.SerializeRequest(req)

		var sResult worker.ParseResult
		c := <-clients //每次从rpc client池中取出一个client

		err := c.Call(config.CrawlServiceRpc,
			sReq, &sResult)

		if err != nil {
			return engine.ParseResult{}, err
		}
		return worker.DeserializeResult(sResult), nil
	}
}
