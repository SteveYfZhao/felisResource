package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserInfo struct {
	UserID      string
	Email       string
	Created     time.Time
	Lastlogin   time.Time
	Disabled    bool
	Createdby   string
	Roles       []string
	Permissions []string
}

const MaxPageSize = 200

func ListUsers(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	pageSize, offset := CalcPageSizeAndOffset(r.Form.Get("pageSize"), r.Form.Get("offset"))

	rt := make([]UserInfo, pageSize)
	db := GetDBHandle()
	//const columns = "username, email, created, lastlogin, disabled, createdby"
	const columns = "username, email, disabled"
	rows, err := db.Query("SELECT "+columns+" FROM useraccount ORDER BY username OFFSET $1 ROWS LIMIT $2", offset, pageSize)
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	rtprt := &rt
	for rows.Next() {
		err := rows.Scan(&(rt[i].UserID), &(rt[i].Email), &(rt[i].Disabled))
		if err != nil {
			log.Print(err)
		}
		log.Print((*rtprt)[i])
		i++
	}
	log.Print(rt[:i])
	return rt[:i], nil
}

func FindUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	params := []string{"userid", "email", "pageSize", "offset"}
	paraMap := MapAllPostParams(r, params)
	pageSize, offset := CalcPageSizeAndOffset(paraMap["pageSize"], paraMap["offset"])

	if !IsEmptyStr(paraMap["userid"]) || !IsEmptyStr(paraMap["email"]) {
		rt := make([]UserInfo, pageSize)
		db := GetDBHandle()
		const columns = "username, email, disabled"
		//const columns = "username, email, created, lastlogin, disabled, createdby"
		rows, err := db.Query("SELECT "+columns+" FROM useraccount ORDER BY username OFFSET $1 ROWS LIMIT $2 WHERE username=$3 OR email=$4", offset, pageSize, paraMap["userid"], paraMap["email"])
		if err != nil {
			log.Fatal(err)
		}

		i := 0
		for rows.Next() {
			err := rows.Scan(&(rt[i].UserID), &(rt[i].Email), &(rt[i].Disabled))
			if err != nil {
				log.Fatal(err)
			}
		}
		return rt, nil
	}
	return nil, nil
}

func ListAllRolls(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	pageSize, offset := CalcPageSizeAndOffset(r.Form.Get("pageSize"), r.Form.Get("offset"))
	rt := make([]string, pageSize)
	db := GetDBHandle()
	const columns = "rolename"
	rows, err := db.Query("SELECT "+columns+" FROM rolelist ORDER BY rolename OFFSET $1 ROWS  LIMIT $2", offset, pageSize)
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for rows.Next() {
		err := rows.Scan(&rt[i])
		if err != nil {
			log.Fatal(err)
		}
	}
	return rt, nil
}

func ListAllPerms(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	pageSize, offset := CalcPageSizeAndOffset(r.Form.Get("pageSize"), r.Form.Get("offset"))
	rt := make([]string, pageSize)
	db := GetDBHandle()
	const columns = "permissionname"
	rows, err := db.Query("SELECT "+columns+" FROM permissionlist ORDER BY permissionname OFFSET $1 ROWS  LIMIT $2", offset, pageSize)
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for rows.Next() {
		err := rows.Scan(&rt[i])
		if err != nil {
			log.Fatal(err)
		}
	}
	return rt, nil
}

func CalcPageSizeAndOffset(rawPageSize string, rawOffset string) (int, int) {
	offset, err := strconv.Atoi(rawOffset)
	if err != nil {
		offset = 0
	}
	pageSize, err := strconv.Atoi(rawPageSize)
	if err != nil || pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return pageSize, offset
}

func MapAllPostParams(r *http.Request, params []string) map[string]string {
	err := r.ParseForm()
	rt := make(map[string]string)
	if err == nil {
		for _, para := range params {
			rt[para] = r.PostForm[para][0]
		}
		return rt

	}
	return nil
}

func IsEmptyStr(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
