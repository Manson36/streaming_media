package api

import (
	"github.com/streaming_media/api/dbops"
	"strconv"
	"sync"
	"time"
)
var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func nowInMilli() int64 {
	return time.Now().UnixNano()/1000000
}

func deleteExpireSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

