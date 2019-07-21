package streamserver

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func (m *middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.l.getConn() {
		sendErrorResponse(w, http.StatusTooManyRequests, "too many request")
		return
	}

	defer m.l.releaseConn()

	m.r.ServeHTTP(w, r)
}

func newMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := &middleWareHandler{}
	m.r = r
	m.l = newConnLimiter(cc)

	return m
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
	mh := newMiddleWareHandler(r, 2)
	http.ListenAndServe(":9000", mh)
}
