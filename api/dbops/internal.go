package dbops

import (
	"database/sql"
	"github.com/streaming_media/api/defs"
	"log"
	"strconv"
	"sync"
)

//1.往DB的session中写入一个session
func InserSession(sid string, ttl int64, uname string) error {
	ttlstr := strconv.FormatInt(ttl, 10)

	stmtIns, err := dbConn.Prepare(
		"insert into sessions (session_id, TTl, login_name) values (?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(sid, ttlstr, uname)
	if err != nil {
		return err
	}

	return nil
}

//2.从DB中获取session信息
func RetrieveSession(sid string) (*defs.SimpleSession, error) {
	ss := &defs.SimpleSession{}
	stmtOut, err := dbConn.Prepare(
		"select TTl, login_name from sessions where session_id = ?")
	if err != nil {
		return nil, err
	}

	defer stmtOut.Close()

	var ttl string
	var uname string

	err = stmtOut.QueryRow(sid).Scan(&ttl, &uname)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if res, err := strconv.ParseInt(ttl, 10, 64); err != nil {
		ss.TTL = res
		ss.Username = uname
	}else {
		return nil, err
	}

	return ss, err
}

func RetrieveAllSession() (*sync.Map, error) {
	m := &sync.Map{}

	stmtOut, err := dbConn.Prepare("select * from sessions")
	if err != nil {
		return nil, err
	}

	defer stmtOut.Close()

	rows, err := stmtOut.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id string
		var ttlstr string
		var login_name string

		if err := rows.Scan(&id, &ttlstr, &login_name); err != nil {
			log.Printf("retrieve sessions error: %v", err)
			break
		}

		if ttl, err := strconv.ParseInt(ttlstr, 10, 64); err != nil {
			ss := &defs.SimpleSession{Username: login_name, TTL: ttl}
			m.Store(id, ss)
			log.Printf("session id: %s, ttl: %d", id, ss.TTL)
		}
	}

	return m, nil 
}