package scheduler

import "u2pppw/crawler/crawler-queue-1/engine"

type QueuedScheduler struct {
	requestChan chan engine.Request      //用于接收request的channel
	workerChan  chan chan engine.Request //表示workerChan是channel,对外传输的是chan engine.Request
}

//把request递交给调度器中的channel
func (s *QueuedScheduler) Submit(r engine.Request) {
	go func() { s.requestChan <- r }()
}

//由外边通知调度器worker准备好了，可以负责接收请求request了
//在concurrent.go中,实际是由WorkerChan()函数返回的channel传给s.workerChan
func (s *QueuedScheduler) WorkerReady(w chan engine.Request) {
	s.workerChan <- w
}

//由调度器来生成request channel
func (s *QueuedScheduler) WorkerChan() chan engine.Request {
	return make(chan engine.Request)
}

//运行调度器
func (s *QueuedScheduler) Run() {
	s.workerChan = make(chan chan engine.Request)
	s.requestChan = make(chan engine.Request)
	go func() {
		var requestQ []engine.Request     //所有请求队列
		var workerQ []chan engine.Request //所有工作队列
		for {
			var activeRequest engine.Request     //当前请求
			var activeWorker chan engine.Request //当前worker，用于接收activeRequest

			//分别从总的请求队列、工作队列中取出第一个队列的元素作为当前活跃的request、worker
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeWorker = workerQ[0]
				activeRequest = requestQ[0]
			}
			select {

			//从调度器把新来的请求加入请求队列
			//submit()后,s.requestChan有内容,可执行
			case r := <-s.requestChan:
				requestQ = append(requestQ, r)

			//从调度器把新来的worker加入worker队列
			//WorkerReady()后,s.workerChan有内容,可执行
			case w := <-s.workerChan:
				workerQ = append(workerQ, w)

			//如果是当前的请求activeRequest给activeWorker，则此时阻塞在createWorker()中的channel得以运行
			//			则要更新总的worker/request队列
			case activeWorker <- activeRequest:
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}

		}
	}()
}
