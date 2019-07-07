package session

import (
	"github.com/streaming_media/api/dbops"
	"github.com/streaming_media/api/defs"
	"github.com/streaming_media/api/utils"
	"sync"
	"time"
)

//这个simplesession在别的地方也会用到，我们放到defs中
//type SimpleSession struct {
//	Username string //login name
//	TTL int64 //用来检查用户是否过期
//}

var sessionMap *sync.Map

func init() {
	sessionMap := &sync.Map{} //？？？
}

func nowInMilli() int64 {
	return time.Now().UnixNano()/1000000//将纳秒转换为毫秒
}

func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

//从DB中拉取session，把已经存在的session全部load到api的节点中
func LoadSessionFromDB() {
	//我们为什么在这里没有返回值，是因为这个操作是在内部操作的，外面也不会有人来调用这个函数，或者调用对函数没有显式的需求
	r, err := dbops.RetrieveAllSession()
	if err != nil {
		return
	}

	//sessionMap = r ???
	r.Range(func(k, v interface{}) bool {
		ss := v.(*defs.SimpleSession)
		sessionMap.Store(k, ss)
		return true
	})
}

//为新用户注册新的id
func GenerateNewSessionId(username string) string {
	id, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 30 * 60 * 1000 //severside session valid time: 30 min

	ss := &defs.SimpleSession{Username:username, TTL:ttl}
	sessionMap.Store(id, ss)
	dbops.InserSession(id, ttl, username)

	return id
}

//根据sessionId判断session是否过期，返回username和结果
func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	ct := nowInMilli()

	if ok {
		if ss.(defs.SimpleSession).TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}

		return ss.(defs.SimpleSession).Username, false
	} else {
		ss, err := dbops.RetrieveSession(sid)

		if err != nil || ss == nil {
			return "", true
		}

		if ss.TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}

		sessionMap.Store(sid, ss)
		return ss.Username, false
	}

	return "", true
}