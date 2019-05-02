package mysql

import (
	"database/sql"
	"fmt"
	"os"

	cfg "../../config"
	_ "github.com/go-sql-driver/mysql" // annonymous import
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", cfg.MySQLSource)
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err: " + err.Error())
		os.Exit(1)
	}
}

// DBConn: return database instance
func DBConn() *sql.DB {
	return db
}
