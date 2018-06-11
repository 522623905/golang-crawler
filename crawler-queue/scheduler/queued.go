package scheduler

import "../engine"

type QueuedScheduler struct {
	requestChan chan engine.Request //用于接收request的channel
	workerChan  chan chan engine.Request
}

//把request递交给调度器中的channel
func (s *QueuedScheduler) Submit(r engine.Request) {
	s.requestChan <- r
}

//由外边通知调度器的worker，请求request准备完毕
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
			//分别从总的请求队列、工作队列中取出一个队列元素给当前的request、worker
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeWorker = workerQ[0]
				activeRequest = requestQ[0]
			}
			select {
			//从调度器把新来的请求加入请求队列
			//submit()后可执行
			case r := <-s.requestChan:
				requestQ = append(requestQ, r)
				//从调度器把新来的worker加入worker队列
				//WorkerReady()后可执行
			case w := <-s.workerChan:
				workerQ = append(workerQ, w)
				//把当前的请求activeRequest给activeWorker，并更新总的worker/request队列
			case activeWorker <- activeRequest:
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}

		}
	}()
}
