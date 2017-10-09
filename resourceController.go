package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetResourceTypeCL(w http.ResponseWriter, r *http.Request) {
	resId := r.Form.Get("resId")
	i, err := strconv.Atoi(resId)
	if err != nil {
		fmt.Fprintf(w, "invalid parameter")
	}
	restype := GetResourceTypeCore(i)
	fmt.Fprintf(w, restype)
}

func GetResourceTypeCore(resId int) string {
	db := GetDBHandle()
	var restype string
	err := db.QueryRow("SELECT type FROM resourcelist WHERE id = resId").Scan(&restype)
	if err != nil {
		log.Fatal(err)
	}
	return restype
}

func GetResourceTags(resId int) *map[string]string {
	db := GetDBHandle()
	var result map[string]string
	rows, err := db.Query("SELECT resourcetaglist.tagid, resourcetagvalues.value FROM resourcetagvalues INNER JOIN resourcetaglist ON resourcetagvalues.tagid = resourcetaglist.id WHERE resourcetagvalues.resource = $1", resId)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var k, v string
		err := rows.Scan(&k, &v)
		if err != nil {
			log.Fatal(err)
		}
		result[k] = v
	}
	return &result
}

func AddAvailTime(w http.ResponseWriter, r *http.Request) {
	resID := r.Form.Get("resid")
	resType := r.Form.Get("restype")
	rulestartdate := r.Form.Get("rulestartdate")
	ruleenddate := r.Form.Get("ruleenddate")
	availstarttime := r.Form.Get("availstarttime")
	availendtime := r.Form.Get("availendtime")
	endonnextday := r.Form.Get("endonnextday")
	freq := r.Form.Get("freq")
	bywkday := r.Form.Get("bywkday")
	bydate := r.Form.Get("bydate")
	userperm := r.Form.Get("userperm")
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	db := GetDBHandle()

	db.QueryRow("INSERT INTO resourceavailabletime(resource,resourcetype,rulestartdate,ruleenddate,availstarttime,availendtime,endonnextday,freq,bywkday,bydate,userperm,createdby,created) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);",
		resID, resType, rulestartdate, ruleenddate, availstarttime, availendtime, endonnextday, freq, bywkday, bydate, userperm, currentuser, time.Now())

	fmt.Println("hit AddAvailTime")
}

type ResourceInfo struct {
	id   int
	name string
}

func getAllResourceForUser(w http.ResponseWriter, r *http.Request) {
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	resources := getAllResourceForUserCore(currentuser)
	rvalues, err := json.Marshal(resources)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(rvalues))
	fmt.Println("hit getAllResourceForUser")
}
func getAllResourceForUserCore(username string) *[]ResourceInfo {
	userperms := GetAllPermsofUser(username)
	db := GetDBHandle()
	rows, err := db.Query("SELECT id, displayname FROM resourcelist WHERE viewpermission IN $1", userperms)
	if err != nil {
		log.Fatal(err)
	}
	var result []ResourceInfo
	for rows.Next() {
		var resid int
		var resname string
		err := rows.Scan(&resid, &resname)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, ResourceInfo{resid, resname})
	}
	return &result
}

func getAllAvailResource(w http.ResponseWriter, r *http.Request) {

	db := GetDBHandle()
	db.Query("SELECT ")
}

func AddResource(w http.ResponseWriter, r *http.Request) {
	resourceid := r.Form.Get("resourceid")

	displayname := r.Form.Get("displayname")
	restype := r.Form.Get("type")
	viewpermission := r.Form.Get("viewpermission")
	bookpermission := r.Form.Get("bookpermission")
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	db := GetDBHandle()

	db.QueryRow("INSERT INTO resourcelist(resourceid,displayname,type,viewpermission,bookpermission,createdby,created) VALUES($1,$2,$3,$4,$5,$6,$7);",
		resourceid, displayname, restype, viewpermission, bookpermission, currentuser, time.Now())

	fmt.Println("hit AddResource")
}

func BookResource(w http.ResponseWriter, r *http.Request) {
	resourceid := r.Form.Get("resourceid")
	starttime := r.Form.Get("starttime")
	endtime := r.Form.Get("endtime")
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}

	layout := "2006-01-02T11:04:05-0400"

	st, err := time.Parse(layout, starttime)

	if err != nil {
		fmt.Println(err)
	}
	et, err := time.Parse(layout, endtime)

	if err != nil {
		fmt.Println(err)
	}

	islegal := IsRequestLegal(currentuser, resourceid, st, et)
	if islegal {
		db := GetDBHandle()

		db.QueryRow("INSERT INTO resourcelist(resource,bookedforuser,bookstart,bookend,createdby,created) VALUES($1,$2,$3,$4,$5,$6);",
			resourceid, currentuser, starttime, endtime, currentuser, time.Now())

	}

	fmt.Println("hit AddResource")
}

func ArchivePastBooking() {

}
