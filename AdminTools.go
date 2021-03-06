package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	Enabled     bool
	Createdby   string
	Roles       map[string]bool
	Permissions map[string]bool
}

type UserInfoNullable struct {
	UserID      sql.NullString
	Email       sql.NullString
	Created     sql.NullString
	Lastlogin   sql.NullString
	Enabled     sql.NullBool
	Createdby   sql.NullString
	Roles       []string
	Permissions []string
}

const MaxPageSize = 200

func ListUsers(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	pageSize, offset := RegulatePageSizeAndOffset(r.Form.Get("pageSize"), r.Form.Get("offset"))

	rt := make([]UserInfo, pageSize)
	db := GetDBHandle()
	//const columns = "username, email, created, lastlogin, enabled, createdby"
	const columns = "username, email, enabled"
	rows, err := db.Query("SELECT "+columns+" FROM useraccount ORDER BY username OFFSET $1 ROWS LIMIT $2", offset, pageSize)
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	rtprt := &rt
	for rows.Next() {
		err := rows.Scan(&(rt[i].UserID), &(rt[i].Email), &(rt[i].Enabled))
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
	log.Print("enter finduser ep")

	if r.Method == "POST" {
		//params := []string{"userid", "email", "pageSize", "offset"}
		paraMap, err := MapURLEncodedPostParams(r)
		log.Println("paraMap", paraMap, err)

		pageSize, offset := RegulatePageSizeAndOffset(paraMap["pageSize"], paraMap["offset"])
		log.Println("offset, pageSize", offset, pageSize)

		if IsEmptyStr(paraMap["userid"]) && IsEmptyStr(paraMap["email"]) {
			return nil, errors.New("No search parameters")
		}
		rt := make([]UserInfo, pageSize)
		db := GetDBHandle()
		const columns string = "username, email, enabled"
		rows, err := db.Query("SELECT "+columns+" FROM useraccount WHERE username = $1 OR email=$2 OFFSET $3 ROWS LIMIT $4", paraMap["userid"], paraMap["email"], offset*pageSize, pageSize)
		if err != nil {
			log.Fatal(err)
		}

		i := 0
		for rows.Next() {
			err := rows.Scan(&(rt[i].UserID), &(rt[i].Email), &(rt[i].Enabled))
			if err != nil {
				log.Fatal(err)
			}
			log.Print("i", i, rt)
			i++
		}
		return rt[:i], nil

	}
	return nil, nil
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	uID := r.URL.Query().Get("uid")
	uData := UserInfoNullable{}
	db := GetDBHandle()
	log.Print("uid", uID, len(uID))
	err := db.QueryRow("SELECT username, email, enabled, created, lastlogin, Createdby FROM useraccount WHERE username = $1", uID).Scan(&(uData.UserID), &(uData.Email), &(uData.Enabled), &(uData.Created), &(uData.Lastlogin), &(uData.Createdby))
	if err != nil {
		log.Print("Error getting basic info", err)
		return nil, err
	}
	uDataValue := UserInfo{}
	// init maps to avoid errors
	uDataValue.Roles = make(map[string]bool)
	uDataValue.Permissions = make(map[string]bool)
	if uData.UserID.Valid {
		uDataValue.UserID = uData.UserID.String
	}

	if uData.Email.Valid {
		uDataValue.Email = uData.Email.String
	}
	if uData.Enabled.Valid {
		uDataValue.Enabled = uData.Enabled.Bool
	}
	if uData.Created.Valid {
		layout := "2006-01-02T11:04:05-04"
		time, err := time.Parse(uData.Created.String, layout)
		if err == nil {
			uDataValue.Created = time
		}
	}
	if uData.Lastlogin.Valid {
		layout := "2006-01-02T11:04:05-04"
		time, err := time.Parse(uData.Lastlogin.String, layout)
		if err == nil {
			uDataValue.Lastlogin = time
		}
	}
	if uData.Createdby.Valid {
		uDataValue.Createdby = uData.Createdby.String
	}
	allRoles := AllRolesSlim()
	allPerms := AllPermsSlim()
	uRoles := GetAllRolesOfUser(uID)
	uPerms := GetAllPermsofUser(uID)

	for _, roleName := range allRoles {
		uDataValue.Roles[roleName] = false
	}
	for _, uR := range uRoles {
		uDataValue.Roles[uR] = true
	}

	for _, permName := range allPerms {
		uDataValue.Permissions[permName] = false
	}
	for _, uP := range uPerms {
		uDataValue.Permissions[uP] = true
	}
	return uDataValue, err
}

func ListAllRolls(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	pageSize, offset := RegulatePageSizeAndOffset(r.Form.Get("pageSize"), r.Form.Get("offset"))
	rt := make([]string, pageSize)
	db := GetDBHandle()
	const columns = "rolename"
	rows, err := db.Query("SELECT "+columns+" FROM rolelist ORDER BY rolename OFFSET $1 ROWS  LIMIT $2;", offset, pageSize)
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

func AllRolesSlim() []string {
	db := GetDBHandle()
	rows, err := db.Query("SELECT rolename FROM rolelist ORDER BY rolename;")
	if err != nil {
		log.Fatal(err)
	}
	return RowsToStringSlice(rows)
}

func AllPermsSlim() []string {
	db := GetDBHandle()
	rows, err := db.Query("SELECT permissionname FROM permissionlist ORDER BY permissionname;")
	if err != nil {
		log.Fatal(err)
	}
	return RowsToStringSlice(rows)
}

func ListAllPerms(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	pageSize, offset := RegulatePageSizeAndOffset(r.Form.Get("pageSize"), r.Form.Get("offset"))
	rt := make([]string, pageSize)
	db := GetDBHandle()
	const columns = "permissionname"
	rows, err := db.Query("SELECT "+columns+" FROM permissionlist ORDER BY permissionname OFFSET $1 ROWS  LIMIT $2;", offset, pageSize)
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

func RegulatePageSizeAndOffset(rawPageSize string, rawOffset string) (int, int) {
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

func GetContentType(r *http.Request) string {
	ct := r.Header.Get("Content-Type")
	log.Print("Content-Type", ct)
	return ct
}

func MapURLEncodedPostParams(r *http.Request) (map[string]string, error) {

	const ContTypeJSON = "application/json"
	const ContTypeFormEncode = "application/x-www-form-urlencoded"
	cType := GetContentType(r)

	if strings.Contains(cType, ContTypeJSON) {
		rt := make(map[string]string)
		rawData := make(map[string]json.RawMessage)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error reading body:", err)
		}
		bodystr := string(body)
		log.Println(bodystr)

		err = json.Unmarshal(body, &rawData)
		if err != nil {
			fmt.Println("error Unmarshal:", err)
			return nil, err
		}
		// And now set a new body, which will simulate the same data we read:
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		for k, v := range rawData {
			if v != nil {
				rt[k] = string(v)
			}
		}

		quoteByte := []byte("\"")[0]
		for k, v := range rawData {

			if v != nil {
				vVal := string(v)
				if len(vVal) > 1 && vVal[0] == quoteByte && vVal[len(vVal)-1] == quoteByte {
					if len(vVal) > 2 {
						vVal = vVal[1 : len(vVal)-1]
					} else {
						vVal = ""
					}
				}
				rt[k] = vVal
			}
		}

		return rt, err

	} else if strings.Contains(cType, ContTypeFormEncode) {
		err := r.ParseForm()
		if err != nil {
			return nil, err
		}

		rt := make(map[string]string)
		rawData := r.Form
		for k, v := range rawData {
			if v != nil && len(v) > 0 {
				rt[k] = v[0]
			}
		}

		return rt, err

		/*
			if err == nil && len(params) > 0 {

				for _, para := range params {
					log.Print("Para: ", para)
					log.Print("r. PostFormValue(para): ", r.PostFormValue(para))
					rt[para] = r.PostFormValue(para)
				}
				return rt
			}
		*/
	} else {
		log.Print("Unknown Content-Type", cType)
	}
	return nil, errors.New("Unknown Content-Type")
}

func IsEmptyStr(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
