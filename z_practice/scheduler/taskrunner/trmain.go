package taskrunner

import "time"

type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

