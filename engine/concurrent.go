package engine

//import "log"

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int //工作协程个数
}

type Scheduler interface {
	Submit(Request)
	ConfigureMasterWorkerChan(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	in := make(chan Request)
	out := make(chan ParseResult)
	e.Scheduler.ConfigureMasterWorkerChan(in)

	//创建工作协程负责解析in请求，并返回解析结果给out
	for i := 0; i < e.WorkerCount; i++ {
		createWorker(in, out)
	}

	//循环遍历传入的request，并递交给in chan,则createWorker()阻塞等待的请求得以运行
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	for {
		result := <-out //out的结果由createWorker()解析后返回
		//打印信息
		//		for _, item := range result.Items {
		//			log.Printf("Got item %d: %v", item)
		//		}

		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}
}

//工作协程,等待请求，并解析返回结果
func createWorker(in chan Request, out chan ParseResult) {
	go func() {
		for {
			request := <-in                //阻塞等待有请求到来
			result, err := worker(request) //解析请求，返回结果集
			if err != nil {
				continue
			}
			out <- result //阻塞等待结果集result给至out
		}
	}()
}
