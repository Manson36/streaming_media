package api

import (
	"github.com/streaming_media/api/dbops"
	"github.com/streaming_media/api/defs"
	"github.com/streaming_media/api/utils"
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

func LoadSessionFromDB() {
	r, err := dbops.RetrieveAllSession()
	if err != nil {
		return
	}

	r.Range(func(k, v interface{}) bool {
		ss := v.(defs.SimpleSession)
		sessionMap.Store(k, ss)
		return true
	})
}

func GetNewSessionId(username string) string {
	id, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 30*60*1000

	ss := &SimpleSession{UserName: username, TTL: ttl}
	sessionMap.Store(id, ss)
	dbops.InserSession(id, ttl, username)

	return id
}

func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	ct := nowInMilli()

	if ok {
		if ss.(defs.SimpleSession).TTL < ct {
			deleteExpireSession(sid)
			return "", true
		}

		return ss.(defs.SimpleSession).Username, false
	} else {
		ss, err := dbops.RetrieveSession(sid)
		if err != nil || ss == nil {
			return "", true
		}

		if ss.TTL < ct {
			return "", true
		}

		sessionMap.Store(sid, ss)
		return ss.Username, false
	}

	return "", true
}