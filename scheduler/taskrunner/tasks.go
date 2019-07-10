package taskrunner

import (
	"errors"
	"github.com/streaming_media/scheduler/dbops"
	"log"
	"os"
	"sync"
)

func VideoClearDispatcher(dc dataChan) error {
	res, err := dbops.ReadVideoDeletionRecord(3)
	if err != nil {
		log.Printf("Video clear dispatcher error: %v", err)
		return err
	}

	if len(res) == 0 {
		return errors.New("all task finished")
	}

	for _, id := range res {
		dc <- id
	}

	return nil
}

func VideoClearExecutor(dc dataChan) error {
	//我们要将err带出，但又不可能出现错误就退出；使用sync.Map是线程安全的
	errMap := &sync.Map{}
	var err error

	forloop:
		for {
			select {
			case vid := <-dc:
				go func(id interface{}) { //会出现没有完成删除操作，内容又被读取问题，这里没有影响
				//为什么不直接将调用vid，而是作为参数传进来：在闭包中调用goroutine，实际上会拿到环境参数瞬时的状态
				//而不会将它的状态保存；只有在将参数传入，才会完整的遍历。
					if err :=  deleteVideo(id.(string)); err != nil {//删除文件
						errMap.Store(id, err)
						return
					}

					if err := dbops.DelVideoDeletionRecord(id.(string)); err != nil{
						errMap.Store(id, err)
						return
					}
				}(vid)
			default:
				break forloop
			}
		}

	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			return false //只要有一个error，我们就停止遍历
		}
		return true
	})

	return err
}

//在这里补一下删文件操作
func deleteVideo(vid string) error {
	err := os.Remove(VIDEO_PATH + vid)

	if err != nil && !os.IsNotExist(err) {
		log.Printf("Deleting video error: %v", err)
		return err
	}

	return nil
}
