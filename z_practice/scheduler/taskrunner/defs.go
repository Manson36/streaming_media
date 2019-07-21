package taskrunner

type controlChan chan string
type dataChan chan interface{}
type fn func(dc dataChan) error
