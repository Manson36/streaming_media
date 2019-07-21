package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type MiddleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	//m := MiddleWareHandler{
	//	r: r,
	//	l: NewConnLimiter(cc),
	//}
	m := MiddleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(cc)

	return m
}

func (m MiddleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.l.GetConn() {
		sendErrorResponse(w, http.StatusTooManyRequests, "too many requests")
		return
	}

	defer m.l.ReleaseConn()

	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/videos/:vid-id", streamHandler)

	router.POST("/upload/:vid-id", uploadHandler)

	router.GET("/testpage", testPageHandler)

	return router
}

func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 2)
	http.ListenAndServe(":9000", mh)
}
