package engine

//import "log"

//定义引擎的结构体
type ConcurrentEngine struct {
	Scheduler   Scheduler //调度器
	WorkerCount int       //工作协程个数
}

//定义了调度器需要实现的接口
type Scheduler interface {
	Submit(Request)
	ConfigureMasterWorkerChan(chan Request)
}

//执行
func (e *ConcurrentEngine) Run(seeds ...Request) {
	in := make(chan Request)                  //负责传输request的channel
	out := make(chan ParseResult)             //负责传输解析request成result后的channel
	e.Scheduler.ConfigureMasterWorkerChan(in) //把in与调度器的channel绑定

	//创建工作协程负责解析in请求，并返回解析结果给out
	for i := 0; i < e.WorkerCount; i++ {
		createWorker(in, out)
	}

	//循环遍历传入的request(主函数中传入的url)，并通过调度器递交给in chan,则createWorker()阻塞等待的请求得以运行
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	for {
		result := <-out //out的结果由createWorker()解析后返回
		//打印信息
		//		for _, item := range result.Items {
		//			log.Printf("Got item %d: %v", item)
		//		}

		//再从result结果集中把新的请求递交给调度器
		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}
}

//工作协程,等待请求，并解析返回结果
func createWorker(in chan Request, out chan ParseResult) {
	go func() {
		for {
			//tell scheduler i'm ready

			request := <-in                //阻塞等待有请求到来
			result, err := worker(request) //解析请求，返回结果集
			if err != nil {
				continue
			}
			out <- result //阻塞等待结果集result给至out
		}
	}()
}
