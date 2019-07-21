package scheduler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func vidDelRecHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	vid := ps.ByName("vid-id")

	if len(vid) == 0 {
		sendResponse(w, 400, "video should not be empty")
		return
	}

	err := addVideoDeletionRecord(vid)
	if err != nil {
		sendResponse(w, 500, "internal error")
		return
	}

	sendResponse(w, 200, "")
	return
}
