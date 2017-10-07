package main
import (
	//"errors"
	"fmt"
	//"regexp"
	//"strings"
	//"log"
	"database/sql"
	_ "github.com/lib/pq"
)
const (
	DB_USER     = "postgres"
	DB_PASSWORD = "111111"
	DB_NAME     = "test"
)

var dbhandle *sql.DB

func DBInit() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, _ := sql.Open("postgres", dbinfo)
	return db
}

func GetDBHandle() *sql.DB {
	if dbhandle == nil {
		dbhandle = DBInit()
	}

	return dbhandle
}

func CloseDBHandle() {
	dbhandle.Close()
}

func createInitialTables(){
	
}