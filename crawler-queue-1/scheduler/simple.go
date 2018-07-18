package scheduler

import "u2pppw/crawler/crawler-queue-1/engine"

type SimpleScheduler struct {
	workerChan chan engine.Request
}

func (s *SimpleScheduler) Submit(r engine.Request) {
	// send request down to worker chan
	//	s.workerChan <- r
	go func() { s.workerChan <- r }() //使用goroutine，否则可能会卡住
}

func (s *SimpleScheduler) WorkerReady(chan engine.Request) {
}

func (s *SimpleScheduler) WorkerChan() chan engine.Request {
	return s.workerChan
}

func (s *SimpleScheduler) Run() {
	s.workerChan = make(chan engine.Request)
}
