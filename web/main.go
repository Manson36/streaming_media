package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", homeHandler)//登陆主界面

	router.POST("/", homeHandler)

	router.GET("/userhome", userHomeHandler)//用户登录页面

	router.POST("/userhome", userHomeHandler)

	router.POST("/api", apiHandler) //api透传

	router.GET("/videos/:vid-id", proxyVideoHandler)

	router.POST("/upload/:vid-id", proxyUploadHandler)//proxy透传

	router.ServeFiles("/statics/*filepath", http.Dir("./templates"))

	return router
}

func main() {
	r := RegisterHandler()
	http.ListenAndServe(":8080", r)
}