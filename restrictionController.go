package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func CreateRestrcitionType(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	typename := r.Form.Get("typename")
	currentuser, err := GetUserNamefromCookie(r)
	/*
		if err != nil {
			log.Fatal(err)
		}
	*/
	db := GetDBHandle()
	db.QueryRow("INSERT INTO restrictiontypes(type,createdby,created) VALUES($1,$2,$3);", typename, currentuser, time.Now())

	fmt.Println("hit createGlobalRestrcition")
	return "OK", err
}

func CreateRestrcition(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	restrictionid := r.Form.Get("restrictionid")
	restrictiontype := r.Form.Get("restrictiontype")
	restrictionvalue := r.Form.Get("restrictionvalue")
	resourcetype := r.Form.Get("resourcetype")
	resource := r.Form.Get("resource")
	restag := r.Form.Get("restag")
	restagvalue := r.Form.Get("restagvalue")
	userperm := r.Form.Get("userperm")
	currentuser, err := GetUserNamefromCookie(r)
	/*
		if err != nil {
			log.Fatal(err)
		}
	*/
	db := GetDBHandle()
	db.QueryRow("INSERT INTO restrictions(restrictionid,restrictiontype,restrictionvalue,resourcetype,resource,restag,restagvalue,userperm,createdby,created) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);",
		restrictionid, restrictiontype, restrictionvalue, resourcetype, resource, restag, restagvalue, userperm, currentuser, time.Now())

	fmt.Println("hit CreateRestrcition")
	return "OK", err
}

func UpdateRestrcition(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	restrictionid := r.Form.Get("restrictionid")
	restrictiontype := r.Form.Get("restrictiontype")
	restrictionvalue := r.Form.Get("restrictionvalue")
	resourcetype := r.Form.Get("resourcetype")
	resource := r.Form.Get("resource")
	restag := r.Form.Get("restag")
	restagvalue := r.Form.Get("restagvalue")
	userperm := r.Form.Get("userperm")
	id := r.Form.Get("id")
	db := GetDBHandle()
	db.QueryRow("UPDATE restrictions SET restrictionid=$1,restrictiontype=$2,restrictionvalue=$3,resourcetype=$4,resource=$5,restag=$6,restagvalue=$7,userperm=$8) WHERE id = &9;",
		restrictionid, restrictiontype, restrictionvalue, resourcetype, resource, restag, restagvalue, userperm, id)

	fmt.Println("hit CreateRestrcition")
	return "OK", nil
}

func RemoveRestrcition(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id := r.Form.Get("id")
	db := GetDBHandle()
	db.QueryRow("DELETE FROM restrictions WHERE id = $1);", id)
	fmt.Println("hit RemoveRestriction")
	return "OK", nil
}

func GetRestrictionValuesCL(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	resourcetype := r.Form.Get("resourcetype")
	resource := r.Form.Get("resource")
	//restag := r.Form.Get("restag")
	//restagvalue := r.Form.Get("restagvalue")
	currentuser, err := GetUserNamefromCookie(r)
	/*
		if err != nil {
			log.Fatal(err)
		}
	*/

	userperm := GetAllPermsofUser(currentuser)
	rvalues, err := json.Marshal(GetRestrictionValuesCoreLite(resourcetype, resource, "", userperm))

	fmt.Fprintf(w, string(rvalues))
	fmt.Println("hit GetRestrictionValuesCL")
	return "OK", err
}

func GetRestrictionValuesCoreLite(resourcetype, resource, restrictiontype string, userperm []string) []string {
	db := GetDBHandle()
	var rvalue []string
	rows, err := db.Query(`SELECT type, restrictionvalue FROM restrictions WHERE 
		(resourcetype=$1 OR resourcetype = '') AND 
		(resource = $2 OR resource = '') AND 
		(userperm = "" OR userperm IN $3)
		`, resourcetype, resource, userperm)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var tempResType, tempResValue string
		err := rows.Scan(&tempResType, &tempResValue)
		if err != nil {
			log.Fatal(err)
		}
		if restrictiontype == "" || restrictiontype == tempResValue {
			rvalue = append(rvalue, tempResValue)
		}
	}
	return rvalue
}

func GetRestrictionValuesCore(resourcetype, restrictiontype string, resource int, userperm []string, restagvalue map[string]string) []string {
	db := GetDBHandle()
	var rvalue []string
	var buffer bytes.Buffer
	for k, v := range restagvalue {
		//tempStr := "(restagvalue=" + k + " AND (restagvalue = '' OR restagvalue = " +v+ ")) OR"
		buffer.WriteString("(restagvalue=")
		buffer.WriteString(k)
		buffer.WriteString(" AND (restagvalue = '' OR restagvalue = ")
		buffer.WriteString(v)
		buffer.WriteString(")) OR")
	}
	genQStr := buffer.String()
	genQStr = genQStr[:len(genQStr)-2]
	fmt.Println(genQStr)
	rows, err := db.Query(`SELECT type, restrictionvalue FROM restrictions WHERE 
		(resourcetype=$1 OR resourcetype = '') AND 
		(resource = $2 OR resource = '') AND
		( $3 ) AND  
		(userperm = "" OR userperm IN $4)
		`, resourcetype, resource, genQStr, userperm)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var tempResType, tempResValue string
		err := rows.Scan(&tempResType, &tempResValue)
		if err != nil {
			log.Fatal(err)
		}
		if restrictiontype == "" || restrictiontype == tempResValue {
			rvalue = append(rvalue, tempResValue)
		}
	}
	return rvalue
}

func IsRequestLegal(username string, resourceId string, startTime, endTime time.Time) bool {
	// resourcetype := GetResourceTypeCore(resourceId)
	// restrctions := GetRestrictionValuesCore(resourcetype, "", resourceId, )
	// Restrictions should be separated by type and parsed by their corresponding functions.
	return true
}
