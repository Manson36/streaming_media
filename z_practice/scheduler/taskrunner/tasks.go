package taskrunner

import (
	"errors"
	"github.com/streaming_media/scheduler/dbops"
	"github.com/streaming_media/z_practice/scheduler"
	"log"
	"os"
	"sync"
)

func VideoClearDispatcher(dc dataChan) error {
	res, err := scheduler.ReadVideoDeletionRecord(3)
	if err != nil {
		log.Printf("video clear dispatcher error: %s", err)
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
	errMap:= &sync.Map{}
	var err error

	forloop :
		for {
			select {
			case vid := <- dc:
				go func(id interface{}) {
					if err = deleteVideo(id.(string)); err != nil {
						errMap.Store(id, err)
						return
					}

					if err = dbops.DelVideoDeletionRecord(id.(string)); err != nil {
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
			return false
		}
		return true
	})

	return err
}

func deleteVideo(vid string) error {
	err := os.Remove("./video/" + vid)
	if err != nil && !os.IsNotExist(err){
		log.Printf("video delete error: %s", err)
		return err
	}

	return nil
}