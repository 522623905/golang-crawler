package worker

import "u2pppw/crawler/crawler-queue/engine"

//rpc服务,负责解析请求和返回结果

type CrawlService struct {
}

//解析请求和返回结果
func (CrawlService) Process(
	req Request, result *ParseResult) error {
	engineReq, err := DeserializeRequest(req)
	if err != nil {
		return err
	}

	//engine的worker解析请求并返回engineResult结果
	engineResult, err := engine.Worker(engineReq)
	if err != nil {
		return err
	}

	//把engineResult序列化成ParseResult
	*result = SerializeResult(engineResult)
	return nil
}
