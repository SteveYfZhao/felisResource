package main

import (
	//"errors"
	"fmt"
	"strconv"
	//"regexp"
	//"strings"
	"database/sql"
	"log"

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

func createInitialTables() {

}

func RowsToStringSlice(rows *sql.Rows) []string {
	var result []string
	for rows.Next() {
		var row string
		err := rows.Scan(&row)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, row)
	}
	return result
}

func queryDBTable(columns string, tablename string, order string, join string, where string, pageSize int, offset int) (interface{}, error) {

	db := GetDBHandle()
	if pageSize <= 0 {
		pageSize = 20
	}
	//const columns = "username, email, created, lastlogin, enabled, createdby"

	//cols := strings.Split(columns, ",")
	qstr := "SELECT " + columns + " FROM " + tablename

	if join != "" {
		log.Print("has join")
		qstr = qstr + where
	}

	if where != "" {
		log.Print("has where")
		qstr = qstr + " WHERE " + where
	}
	if order != "" {
		log.Print("has order")
		qstr = qstr + " ORDER BY " + order
	}
	qstr = qstr + " OFFSET " + strconv.Itoa(offset) + " ROWS LIMIT " + strconv.Itoa(pageSize)

	log.Print("qstr = " + qstr)
	rows, err := db.Query(qstr)
	if err != nil {
		log.Fatal(err)
	}

	/*
		colCount, _ := rows.Columns()

		count := len(colCount)
		//rt := make([][]interface{}, pageSize)
		rt := make([]map[string]interface{}, pageSize)
		i := 0
		rtprt := &rt
		valuePtrs := make([]interface{}, count)
		for rows.Next() {
			//rt[i] = make([]interface{}, count)
			rt[i] = make(map[string]interface{})
			tempval := make([]interface{}, count)
			for j := 0; j < count; j++ {
				//valuePtrs[j] = &rt[i][j]
				valuePtrs[j] = &tempval[j]
			}

			err := rows.Scan(valuePtrs...)
			if err != nil {
				log.Print(err)
			}
			for k := 0; k < count; k++ {
				rt[i][colCount[k]] = tempval[k]
			}

			log.Print((*rtprt)[i])
			i++
			log.Print("loop " + strconv.Itoa(i))
		}
		log.Print(rt[:i])
		log.Print(len(rt))
		return rt[:i], nil
	*/

	return convDBRowToDict(rows)
}

func queryDBTableAdv(qstr string) ([]map[string]string, error) {
	db := GetDBHandle()
	rows, err := db.Query(qstr)
	if err != nil {
		log.Fatal(err)
	}
	return convDBRowToDict(rows)
}

func convDBRowToDict(rows *sql.Rows) ([]map[string]string, error) {
	colCount, _ := rows.Columns()

	count := len(colCount)
	//rt := make([][]interface{}, pageSize)
	rt := make([]map[string]string, 0)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {

		tempMap := make(map[string]string)
		tempval := make([]interface{}, count)
		for j := 0; j < count; j++ {

			valuePtrs[j] = &tempval[j]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Print(err)
		}
		for k := 0; k < count; k++ {
			tempMap[colCount[k]] = fmt.Sprintf("%v", tempval[k])
		}
		rt = append(rt, tempMap)
	}
	return rt, nil
}
