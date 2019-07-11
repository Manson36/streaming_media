package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/streaming_media/scheduler/taskrunner"
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
