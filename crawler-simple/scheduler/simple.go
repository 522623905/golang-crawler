package scheduler

import "../engine"

type SimpleScheduler struct {
	workerChan chan engine.Request
}

//把request请求递交给调度器的channel
func (s *SimpleScheduler) Submit(r engine.Request) {
	// send request down to worker chan
	//	s.workerChan <- r
	go func() { s.workerChan <- r }() //使用goroutine，否则可能会卡住
}

//设置调度器的channel
func (s *SimpleScheduler) ConfigureMasterWorkerChan(c chan engine.Request) {
	s.workerChan = c
}
