package taskrunner

type Runner struct {
	Controller controlChan
	Error controlChan //返回的就是controlChan 中close部分，我们把控制信息和错误信息分开，便于维护
	Data dataChan
	dataSize  int
	longLived bool //我们开了很多channel，这个用来判断channel是否是长期存活；是：我们不回收，不是：我们就回收资源
	Dispatcher fn
	Executor fn
}

func NewRunner(size int, longlived bool, d fn, e fn) *Runner {
	return &Runner{
		Controller: make(chan string, 1), //我们在这里使用的是带buffer的channel，非阻塞式的；null buffer很可能在一开始就卡住
		Error: make(chan string, 1),
		Data: make(chan interface{}, size),
		dataSize: size,
		longLived: longlived,
		Dispatcher: d,
		Executor: e,
	}
}

//常驻任务
func (r *Runner) startDispatch() {
	defer func() {
		if !r.longLived {
			close(r.Controller)
			close(r.Error)
			close(r.Data)
		}
	}()
	
	for {
		select {
		case c := <- r.Controller:
			if c == READY_TO_DISPATCH {
				err := r.Dispatcher(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_EXECUTE
				}
			}

			if c == READY_TO_EXECUTE {
				err := r.Executor(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_DISPATCH
				}
			}
		case e := <- r.Error:
			if e == CLOSE {
				return
			}
		default:

		}
	}
}

func (r *Runner) StartAll() {//往controller中写入一个内容，开始执行
	r.Controller <- READY_TO_DISPATCH
	r.startDispatch()
}