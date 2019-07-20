package dbops

import (
	"database/sql"
	"github.com/streaming_media/api/defs"
	"log"
	"strconv"
	"sync"
)

//session 的DB操作：
func InsertSession(sid string, ttl int64, uname string) error {
	ttlstr := strconv.FormatInt(ttl, 10)

	stmtIns, err := dbConn.Prepare(
		"insert into sessions (session_id, TTL, login_name) values (?, ?, ?)")
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

func RetrieveSession(sid string) (*defs.SimpleSession, error) {
	ss := &defs.SimpleSession{}
	stmtOut, err := dbConn.Prepare(
		"select TTL, login_name from sessions where session_id = ?")
	if err != nil {
		return nil, err
	}

	defer stmtOut.Close()

	var ttl string
	var uname string

	err = stmtOut.QueryRow(sid).Scan(ttl, uname)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	res, err := strconv.ParseInt(ttl, 10, 64)
	if err == nil {
		ss.TTL = res
		ss.Username = uname
	} else {
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

		err := rows.Scan(&id, &ttlstr, &login_name)
		if err != nil {
			log.Printf("retrieve sessions error: %s", err)
			break
		}

		ttl, err := strconv.ParseInt(ttlstr, 10, 64)
		if err == nil {
			ss := defs.SimpleSession{Username:login_name, TTL:ttl}
			m.Store(id, ss)
			log.Printf("session id: %s, ttl: %s", id, ttl)
		}
	}

	return m, err
}

func DeleteSession(sid string) error {
	stmtOut, err := dbConn.Prepare(
		"delete from sessions where session_id = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	defer stmtOut.Close()

	_, err = stmtOut.Exec(sid)
	if err != nil {
		return err
	}

	return nil
}