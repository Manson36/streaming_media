package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user/", CreateUser)
	router.POST("/user/user_name", Login)
	router.GET("/user/user_name", GetUserInfo)

	return router
}

func main() {
	r := RegisterHandlers()
	http.ListenAndServe(":8000", r)
}


