package taskrunner

import "time"

type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

func newWorker(internal time.Duration, r *Runner) *Worker {
	return &Worker{
		ticker: time.NewTicker(internal * time.Second),
		runner: r,
	}
}

func (w *Worker)startWorker() {
	for {
		select {
		case <-w.ticker.C:
			go w.runner.StartAll()
		}
	}
}

func start() {
	r := newRunner(3, true, VideoClearDispatcher, VideoClearExecutor)
	w := newWorker(3,r)
	go w.startWorker()
}