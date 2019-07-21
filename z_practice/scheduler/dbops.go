package scheduler

import (
	"database/sql"
	"log"
)

var (
	dbConn *sql.DB
	err error
)

func init() {
	dbConn, err = sql.Open("mysql",
		"root:admin@(localhost:3306)/video_server?charset= utf8")
	if err !=nil {
		panic(err)
	}
}

func addVideoDeletionRecord(vid string) error {
	stmtIns, err := dbConn.Prepare("insert into video_del_rec (video_id) values (?)")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(vid)
	if err != nil {
		log.Printf("add vid error: %s", err)
		return err
	}
	return nil
}
