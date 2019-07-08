package main

import (
	"github.com/streaming_media/api/defs"
	"github.com/streaming_media/api/session"
	"net/http"
)

var HEADER_FIELD_SESSION= "X-Session-Id"
var HEADER_FIELD_UNAME= "X-user-Name"


//Check if the current user has the permission
//Use session id to do the check
func validateUserSession(r *http.Request) bool {
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}

	uname, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}

	r.Header.Add(HEADER_FIELD_UNAME, uname)
	return true
}

func ValidateUser(w http.ResponseWriter, r *http.Request) bool {
	uname := r.Header.Get(HEADER_FIELD_UNAME)
	if len(uname) ==0 {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return false
	}
	return true
}