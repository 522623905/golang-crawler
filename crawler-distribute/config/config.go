package config

const (
	//Parser names
	ParseCity     = "ParseCity"
	ParseCityList = "ParseCityList"
	//	ParseProfile  = "ParseProfile"
	ParseProfile = "ProfileParser"
	NilParser    = "NilParser"

	//elasticsearch index
	ElasticIndex = "dating_profile"

	//rpc endpoints
	ItemSaverRpc    = "ItemSaverService.Save"
	CrawlServiceRpc = "CrawlService.Process"

	//Rate limiting
	Qps = 20
)
