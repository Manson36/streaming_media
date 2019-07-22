package taskrunner

type controlChan chan string
type dataChan chan interface{}
type fn func(dc dataChan) error

const (
	READY_TO_DISPATCH = "d"
	READY_TO_EXECUTE = "e"
	CLOSE = "c"
)