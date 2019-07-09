package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//check session
	validateUserSession(r)
	m.r.ServeHTTP(w, r)
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", CreateUser)

	router.POST("/user/:user_name", Login)//:之后的是参数，看Login内容

	router.GET("/user/:user_name", GetUserInfo)

	router.POST("user/:username/videos", AddNewVideo)

	router.GET("user/:user_name/videos", ListAllVideos)

	router.DELETE("/user/:user_name/videos/:vid-id", DeleteVideo)

	router.POST("/video/:video-id/comments", PostComment)

	router.GET("/video/:vid-id/comments", ShowComments)

	return router
}


func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r)
	http.ListenAndServe(":8000", mh)

}