package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &UserCredential{}

	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, ErrorRequestBodyParseFailed)
		return
	}


}
