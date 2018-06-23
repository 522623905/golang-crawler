package engine

type ConcurrentEngine struct {
	Scheduler   Scheduler //调度器
	WorkerCount int       //工作协程个数
	ItemChan    chan Item //用于与elasticSearch通信的channel
}

//调度器需要实现的接口
type Scheduler interface {
	ReadyNotifier             //一个准备完毕并通知的接口
	Submit(Request)           //递交请求给调度器
	WorkerChan() chan Request //调度器返回request channel
	Run()                     //运行调度器
}

//接口，准备并通知
type ReadyNotifier interface {
	WorkerReady(w chan Request) //用于由外边通知调度器的worker，请求request准备完毕
}

//执行引擎
func (e *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParseResult) //用于传递解析请求的结果的channel
	e.Scheduler.Run()             //运行调度器

	//创建工作协程负责解析in请求，并返回解析结果给out
	for i := 0; i < e.WorkerCount; i++ {
		createWorker(e.Scheduler.WorkerChan(),
			out, e.Scheduler)
	}

	//循环遍历传入的request，并递交给调度器,则createWorker()阻塞等待的请求得以运行
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	for {
		result := <-out //out的结果由createWorker()解析后返回
		//存储信息
		for _, item := range result.Items {
			go func() { e.ItemChan <- item }()
		}

		for _, request := range result.Requests {
			//Url去重
			if isDuplicate(request.Url) {
				continue
			}
			e.Scheduler.Submit(request)
		}
	}
}

//工作协程,等待请求，并解析返回结果
func createWorker(in chan Request,
	out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			//tell scheduler i'm ready
			ready.WorkerReady(in)
			request := <-in                //阻塞等待有请求到来
			result, err := worker(request) //解析请求，返回结果集
			if err != nil {
				continue
			}
			out <- result //阻塞等待结果集result给至out
		}
	}()
}

//用于url去重
var visitedUrls = make(map[string]bool)

//用于url去重
func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}
	visitedUrls[url] = true
	return false
}
