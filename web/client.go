package main

import (
	"bytes"
	"encoding/json"
	"github.com/streaming_media/web/config"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var httpClient *http.Client //首先声明一个全局的httpClient，这样它就可以被复用

func init() {  //然后我们将它初始化出来
	httpClient = &http.Client{}
}

func request(b *ApiBody, w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error

	//后来的config添加：
	u, _ := url.Parse(b.Url)//net/url的包是一个很不错的包
	u.Host = config.GetLBAddr() + ":" + u.Port()//获取真正的host在映射到它的port上
	newUrl := u.String()//这是真正的映射到api上，并且监听在cloud上的url
	//添加完成，将下面的代码中，b.url改为newUrl

	switch b.Method {
	case http.MethodGet:
		req, _ := http.NewRequest("Get", newUrl, nil)
		req.Header = r.Header//???目的是啥？
		resp, err = httpClient.Do(req)//这么做的好处是：前台(java script)处理到的response会和我们api service完全保持一致
		if err != nil {
			log.Printf("%s",err)
			return
		}

		normalResponse(w, resp)
	case http.MethodPost:
		req, _ := http.NewRequest("POST", newUrl, bytes.NewBuffer([]byte(b.ReqBody)))
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf("%s",err)
			return
		}

		normalResponse(w, resp)
	case http.MethodDelete:
		req, _ := http.NewRequest("Delete", newUrl, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf("%s",err)
			return
		}

		normalResponse(w, resp)
	}
}

//处理一下格式错误，或者说发过来的request可能是会有问题的，没有问题我们就透传
func normalResponse(w http.ResponseWriter, r *http.Response) {
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		re, _ := json.Marshal(ErrorInternalFaults)
		w.WriteHeader(500)
		io.WriteString(w, string(re))
		return
	}

	w.WriteHeader(r.StatusCode)
	io.WriteString(w, string(res))
}