package configs

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	)

func KoneksiDB()(db *sql.DB, er error){
	db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/bank")
	if err != nil {
		panic(err)
	}
	return
}