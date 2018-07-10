package config

const (
	//Parser names
	ParseCity     = "ParseCity"
	ParseCityList = "ParseCityList"
	ParseProfile  = "ProfileParser"
	NilParser     = "NilParser"

	//elasticsearch index
	ElasticIndex = "dating_profile"

	//rpc endpoints
	//1.存储至ElasticSearch  2.解析request返回result
	ItemSaverRpc    = "ItemSaverService.Save"
	CrawlServiceRpc = "CrawlService.Process"

	//Rate limiting
	Qps = 200
)
