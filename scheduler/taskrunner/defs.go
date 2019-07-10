package taskrunner

//预定义runner中对象
type controlChan chan string

type dataChan chan interface{}

type fn func(dc dataChan) error //我们的dispatcher和executor

//定义controlChan中的消息
const (
	READY_TO_DISPATCH = "d"
	READY_TO_EXECUTE = "e"
	CLOSE = "c"
)
