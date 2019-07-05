package dbops

import ("database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err error
)

//在这里使用初始化定义mysql连接，而不创建返回sql.DB的函数，是为了防止多次执行sql.Open, 创建过多连接。
func init() {
	dbConn, err := sql.Open("mysql",
		"root:admin!@#@(localhost:3306)/video_server?charset=utf8")
	if err != nil {
		panic(err)
	}
}

