package taskrunner

import "time"

type Worker struct {
	ticker *time.Ticker //初始化一个ticker，不断接收从系统发过来的时间，达到我们想要的时间间隔的时候触发任务
	runner *Runner
}

func NewWorker(interval time.Duration, r *Runner) *Worker {
	return &Worker{
		ticker: time.NewTicker(interval * time.Second),
		runner: r,
	}
}

func (w *Worker) startWorker() {
	for {
		select {
		case <-w.ticker.C:
			go w.runner.StartAll()
		}
	}
}

func Start() {
	//Start video file cleaning
	r := NewRunner(3, true, VideoClearDispatcher, VideoClearExecutor)
	w := NewWorker(3, r)
	go w.startWorker()

	//somethings else: 可以创建多个r，w， go w.startWorker 同时执行多个定时器
}