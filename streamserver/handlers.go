package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles("./videos/upload.html")

	t.Execute(w, nil)
}

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//vid := p.ByName("vid-id")//获取vid-id，在main的Register函数中video参数是vid-id，作为streaming的唯一名称
	//vl := VIDEO_DIR + vid //videos 所存在的directory
	//
	//video, err := os.Open(vl)
	//if err != nil {
	//	log.Printf("Error when try to open file: %v", err)//要返回错误，却不知道是什么错误，最好将错误打印出来
	//	sendErrorResponse(w, http.StatusInternalServerError, "internal server")
	//	return
	//}
	//
	//defer video.Close()
	//
	////如果成功，使用SET，我们要把返回的response加一个Header的强制提醒
	////因为我们在传这个文件的时候，我们这个文件可能是没有它的扩展名的，然而他的格式里面真正文件里面的二进制码实际上一样
	////是它的视频的mp4的格式，那么我们在response里面把它的header先加上并且将它的文件Content—Type强制写成Video/mp4
	////在client端也就是browser，浏览器会自动按照mp4来解析，解析完之后，会组成真正的视频来播放
	//w.Header().Set("Content-Type", "video/mp4")
	//http.ServeContent(w, r, "", time.Now(), video)//将文件转换为二进制流传送给browser端

	log.Println("Entered the streamHandler")
	//url是云的公网地址 +上传的文件夹的path(/videos/) + 名字
	targetUrl := "http://avenssi-videos2.oss-cn-qingdao.aliyuncs.com/videos/" + p.ByName("vid-id")
	http.Redirect(w, r, targetUrl, 301)
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//首先做一些静态检查
	//设置我们能读的最大缓冲区是多少或者最大文件是多少
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "File is too big")
		return
	}

	file, _, err := r.FormFile("file")//在form这个tag中有一个叫做name=“”, 在这里我们设成file
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal error")
	}

	filename := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR + filename, data, 0666)
	if err != nil {
		log.Printf("Write file err: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}

	//将文件upload到本地之哦户，upload到oss
	ossfn := "video/" + filename
	path := "./video/" + filename //本地文件的位置
	bn := "avenssi-video2" //oss bucket name

	ret := UploadToOss(ossfn, path, bn)
	if !ret  {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal error")
		return
	}

	os.Remove(path)

	//写入成功了，给前端返回状态码
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Uploaded successfully")
}
