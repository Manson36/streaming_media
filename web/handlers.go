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

type UserPage struct {
	Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err1 := r.Cookie("username") //通过cookie的key查找内容
	sid, err2 := r.Cookie("session")

	if err1!= nil || err2 != nil {
		p := &HomePage{Name: "avenssi"}

		//非登陆过的用户，普通的visitor
		t, e := template.ParseFiles("./templates/home.html")//将模板parse成能够理解的二进制文件，而不是一个html文件
		if e != nil {
			log.Printf("Parsing templates home.html error: %s", e)
			return
		}

		t.Execute(w, p)//将模板和需要的变量一起执行，通过w返回给前端

		return
	}

	//登陆过的用户
	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound)
		return
	}
}

func userHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//如果是visitor跳转到homePage，如果是user就把它的一些内容展现出来
	cname, err1 := r.Cookie("username") //通过cookie的key查找内容
	_, err2 := r.Cookie("session")

	if err1!= nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fname := r.FormValue("username")//从表单中获取信息

	var p *UserPage
	if len(cname.Value) != 0 {
		p = &UserPage{Name: cname.Value}
	} else if len(fname) != 0 {
		p = &UserPage{Name: fname}
	}

	t, e := template.ParseFiles("./templates/userhome.html")
	if e != nil {
		log.Printf("Parsing userhome.html error：%s", e)
		return
	}

	t.Execute(w, t)
}

