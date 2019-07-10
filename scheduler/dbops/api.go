package dbops

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func AddVideoDeletionRecord(vid string) error {
	stmtInt, err := dbConn.Prepare("insert into vidoe_del_rec (video_id) values(?)")
	if err != nil {
		return err
	}

	defer stmtInt.Close()

	_, err =stmtInt.Exec(vid)
	if err != nil {
		log.Printf("AddVideoDeletionRecord error: %v", err)
		return err
	}

	return nil
}
