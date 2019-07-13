package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
)

type HomePage struct {
	Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p := &HomePage{Name: "avenssi"}

	t, e := template.ParseFiles("./template/home.html")//将模板parse成能够理解的二进制文件，而不是一个html文件
	if e != nil {
		log.Printf("Parsing template home.html error: %s", e)
		return
	}

	t.Execute(w, p)//将模板和需要的变量一起执行，通过w返回给前端

	return
}
