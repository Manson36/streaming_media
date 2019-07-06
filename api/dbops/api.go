package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streaming_media/api/defs"
	"github.com/streaming_media/api/utils"
	"log"
	"time"
)

//api的操作就是对数据库的增删改查
func AddUserCredential(loginName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("insert into users (login_name, pwd) values (?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmtIns.Close()

	//执行预处理的mysql指令
	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}

	return nil
}

func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare("select pwd from users where login_name = ?")
	if err != nil {
		log.Printf(	"%s", err)
		return "", err
	}
	defer stmtOut.Close()

	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	//注意：如果loginName 在库中不存在，Scan会将err返回，为ErrNoRows
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return pwd, nil
}

func DeleteUser(loginName string, pwd string) error {
	stmtDle, err := dbConn.Prepare("DELETE from users where login_name=? and pwd=?")
	if err != nil {
		log.Printf("DeleteUser error: %s", err)
		return err
	}
	defer stmtDle.Close()

	_, err = stmtDle.Exec(loginName, pwd)
	if err != nil {
		return err
	}

	return nil
}

//Video 的实现
func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {//返回整个的video info
	//create uuid
	vid ,err := utils.NewUUID()
	if err != nil {
		return nil, err
	}

	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05")//这里面的格式的内容是固定的

	stmtIns, err := dbConn.Prepare(`insert into video_info 
					(id, author_id, name, display_ctime) values (?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}

	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: ctime}

	return res, nil
}

func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(
		"select author_id, name, display_time from video_info where vid=?")
	if err != nil {
		log.Printf("GetVideoInfo err: %v", err)
		return nil, err
	}

	defer stmtOut.Close()

	var aid int
	var dct string
	var name string

	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: dct}

	return res, nil
}

func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare("delete from video_info where id=?")
	if err != nil {
		log.Printf("deleteVideoInfo err: %v", err)
		return err
	}

	defer stmtDel.Close()

	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}

	return nil
}