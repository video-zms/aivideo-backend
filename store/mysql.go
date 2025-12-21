package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"axe-backend/config"
)

var MainDB *sqlx.DB


func ConnectMainDB() {
	mainConn := sqlx.MustConnect("mysql", config.GetMainMysqlDsn())
	mainConn.SetMaxOpenConns(500)
	mainConn.SetMaxIdleConns(100)
	MainDB = mainConn.Unsafe()
}