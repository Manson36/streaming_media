package dbops

import "database/sql"

//api的操作就是对数据库的增删改查

func openConn() *sql.DB {

}

func AddUserCredential(loginName string, pwd string) error {

}

func GetUsrCredential(loginName string) (string, error) {

}