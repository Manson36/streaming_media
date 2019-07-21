package streamserver

import (
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func streamHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	vid := ps.ByName("vid-id")
	vl := "./videos/" + vid

	video, err := os.Open(vl)
	if err != nil {
		log.Printf("error when open file: %s", err)
		sendErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}

	defer video.Close()

	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)
}

func uploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//首先做一些静态检查，最大缓冲区是多少
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "file is too big")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}

	filename := ps.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR + filename, data, 0666)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}

	//写入成功，给前端返回
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Upload Successfully")
}
