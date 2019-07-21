package review

import (
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func streamHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	vid := ps.ByName("vid-id")
	vl := "./videos/" + vid

	video, err := os.Open(vl)
	if err != nil {
		return
	}
	defer video.Close()

	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)
}

func uploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, 50 * 1024 * 1024)
	err := r.ParseMultipartForm(50 * 1024 * 1024)
	if err != nil {
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return
	}

	data, _ := ioutil.ReadAll(file)

	filename := ps.ByName("vid-id")
	err = ioutil.WriteFile("./video/" + filename, data, 0666)

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "upload successfully")
}
