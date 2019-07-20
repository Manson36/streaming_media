package dbops

import (
	"database/sql"
	_ "database/sql/driver"
	"github.com/streaming_media/z_practice/api"
	"log"
)

func AddUserCredential(loginName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("insert into users (login_name, pwd) values (?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmtIns.Close()

	_, err =stmtIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}

	return nil
}

func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare("select pwd from users where login_name = ?" )
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}
	defer stmtOut.Close()

	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows{
		return "", err
	}

	return pwd, nil
}

func DeleteUser(loginName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("delete from users where login_name = ? and pwd = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}
	defer stmtDel.Close()

	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(loginName string) (*api.User, error) {
	stmtOut, err := dbConn.Prepare("select id, pwd from users where login_name = ?")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	defer stmtOut.Close()

	var id int
	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&id, &pwd)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	res := &api.User{Id: id, LoginName:loginName, Pwd: pwd}

	return res, nil 
}
