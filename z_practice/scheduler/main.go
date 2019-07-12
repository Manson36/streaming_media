package scheduler

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/streaming_media/scheduler/dbops"
	"github.com/streaming_media/scheduler/taskrunner"
	"io"
	"log"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/video-delete-record/:vid-id", vidDelRecHandler)

	return router
}

func main() {
	go taskrunner.Start()

	r := RegisterHandlers()
	http.ListenAndServe(":9001", r)//这里本身就是一个阻塞的模式，可以直接使用goroutine
}

func vidDelRecHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")

	if len(vid) == 0 {
		sendResponse(w, 400, "video id should not be empty")
		return
	}

	err := dbops.AddVideoDeletionRecord(vid)
	if err != nil {
		sendResponse(w, 500, "Internal server error")
		return
	}

	sendResponse(w, 200, "")
	return
}

func sendResponse(w http.ResponseWriter, sc int, resp string) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}

type ConnLimiter struct {
	concurrentConn int
	bucket chan int
}

func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket: make(chan int, cc),
	}
}

func (cl *ConnLimiter) GetConn() bool {
	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("Reached the rate limitation")
		return false
	}

	cl.bucket <- 1

	return true
}

func (cl *ConnLimiter) ReleaseConn() {
	c := <- cl.bucket
	log.Printf("New connection coming: %d", c)
}

func AddVideoDeletionRecord(vid string) error {
	stmtInt, err := dbConn.Prepare("insert into vidoe_del_rec (video_id) values(?)")
	if err != nil {
		return err
	}

	defer stmtInt.Close()

	_, err =stmtInt.Exec(vid)
	if err != nil {
		log.Printf("AddVideoDeletionRecord error: %v", err)
		return err
	}

	return nil
}

var (
	dbConn *sql.DB
	err error
)

func init() {
	dbConn, err = sql.Open("mysql",
		"root:admin@tcp(localhost:3306)/video_server?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
}

func ReadVideoDeletionRecord(count int) ([]string, error) {
	stmtOut, err := dbConn.Prepare("select video_id from video_del_rec limit ?")

	var ids []string

	if err != nil {
		return ids, err
	}

	defer stmtOut.Close()

	rows, err := stmtOut.Query(count)
	if err != nil {
		log.Printf("Query VideoDeletionRecord error: %v", err)
		return ids, err
	}

	for rows.Next(){
		var id string
		if err := rows.Scan(id); err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func DelVideoDeletionRecord(vid string) error {
	stmtDel, err := dbConn.Prepare("delete from video_del_rec where video_id = ?")
	if err != nil {
		return err
	}

	defer stmtDel.Close()

	_, err = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("Deleting VideoDeletionRecord error: %v", err)
		return err
	}

	return nil
}

type controlChan chan string

type dataChan chan interface{}

type fn func(dc dataChan) error

//定义controlChan中的消息
const (
	READY_TO_DISPATCH = "d"
	READY_TO_EXECUTE = "e"
	CLOSE = "c"

	VIDEO_PATH = "./videos/"
)

type Runner struct {
	Controller controlChan
	Error controlChan
	Data dataChan
	dataSize  int
	longLived bool 
	Dispatcher fn
	Executor fn
}

func NewRunner(size int, longlived bool, d fn, e fn) *Runner {
	return &Runner{
		Controller: make(chan string, 1),
		Error: make(chan string, 1),
		Data: make(chan interface{}, size),
		dataSize: size,
		longLived: longlived,
		Dispatcher: d,
		Executor: e,
	}
}

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

func (r *Runner) StartAll() {
	r.Controller <- READY_TO_DISPATCH
	r.startDispatch()
}