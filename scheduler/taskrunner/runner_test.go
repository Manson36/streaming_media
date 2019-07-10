package taskrunner

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	d := func(dc dataChan) error {
		for i := 0; i < 30; i++ {
			dc <- i
			log.Printf("Dispatcher send: %d", i)
		}

		return nil
	}

	e := func(dc dataChan) error {
		forloop:
		for {
			select {
			case d := <-dc:
				log.Printf("Executor received: %v", d)
			default:
				break forloop
			}
		}

		return errors.New("executor") //如果这里是nil，会循环发送data，接收data内容
	}

	runner := NewRunner(30, false, d, e)
	go runner.StartAll() //这里为什么使用go；是因为startDispatch中有for会阻塞，不会执行下面的代码
	time.Sleep(3 * time.Second)
}
