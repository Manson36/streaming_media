package session

import "sync"

//这个simplesession在别的地方也会用到，我们放到defs中
//type SimpleSession struct {
//	Username string //login name
//	TTL int64 //用来检查用户是否过期
//}

var sessionMap *sync.Map

func init() {
	sessionMap := &sync.Map{}
}

//从DB中拉取session，把已经存在的session全部load到api的节点中
func LoadSessionFromDB() {

}

//为新用户注册新的id
func GenerateNewSessionId(username string) string {

}

//根据sessionId判断session是否过期，返回username和结果
func IsSessionExpired(sid string) (string, bool) {

}