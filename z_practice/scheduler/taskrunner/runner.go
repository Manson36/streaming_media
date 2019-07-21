package taskrunner

type Runner struct {
	Controller controlChan
	Error controlChan
	data dataChan
	dataSize int
	longLived bool
	Dispatcher fn
	Executor fn 
}
