package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streaming_media/api/defs"
	"log"
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

}