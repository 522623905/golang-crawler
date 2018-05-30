package scheduler

import "../engine"

type SimpleScheduler struct {
	workerChan chan engine.Request
}

func (s *SimpleScheduler) Submit(r engine.Request) {
	// send request down to worker chan
	//	s.workerChan <- r
	go func() { s.workerChan <- r }() //使用goroutine，否则可能会卡住
}

func (s *SimpleScheduler) ConfigureMasterWorkerChan(c chan engine.Request) {
	s.workerChan = c
}
