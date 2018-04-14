package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetResourceType(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	resId := r.Form.Get("resId")
	i, err := strconv.Atoi(resId)
	if err != nil {
		fmt.Fprintf(w, "invalid parameter")
	}
	restype := GetResourceTypeCore(i)
	return restype, err
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

func AddAvailTimePlan(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	fmt.Println("hit AddAvailTimePlan")
	return nil, err
}

type BookingInfo struct {
	Id            int
	ResourceId    int
	Bookedforuser string
	Bookstart     time.Time
	Bookend       time.Time
}
type ResBaseInfo struct {
	Id   int
	Name string
}

type ResourceInfo struct {
	Id          int
	Name        string
	BookingStat []BookingInfo
}

func getAllResourceForUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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
	return rvalues, nil
}

func getAllResourceForUserCore(username string) *[]ResourceInfo {
	userperms := GetAllPermsofUser(username)
	db := GetDBHandle()
	flatStr := strings.Join(userperms[:], "','")
	fmt.Println(flatStr)
	qstr := "SELECT id, displayname FROM resourcelist WHERE viewpermission IN ('" + flatStr + "')"
	rows, err := db.Query(qstr)
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
		result = append(result, ResourceInfo{resid, resname, nil})
	}
	return &result
}

func getAllAvailResource(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	db := GetDBHandle()
	rows, err := db.Query("SELECT id, displayname FROM resourcelist")
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
		result = append(result, ResourceInfo{resid, resname, nil})
	}
	return &result, nil
}

func GetBookableResForUserAtTime(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	startTime, _ := extractParams(r, "startTime")
	endTime, _ := extractParams(r, "endTime")
	//1 get all res this user can see
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}

	allresources := getAllResourceForUserCore(currentuser)
	fmt.Println("allresources = ", allresources)
	timestr := "('" + startTime + "'<=bookstart AND bookstart<=" + "'" + endTime + "'" + ") OR (" + "'" + startTime + "'<=bookend AND bookend<=" + "'" + endTime + "')"

	// 2 get all unavailable resource by lookup booking table
	bqstr := "SELECT resource, bookedforuser, bookstart, bookend FROM resourcebooking WHERE bookedforuser <> '" + currentuser + "' AND " + timestr + ";"
	fmt.Println("bqstr = ", bqstr)
	blocked, _ := queryDBTableAdv(bqstr)
	fmt.Println("blocked = ", blocked)

	// 3 compare and generate new list.

	rt := make([]ResourceInfo, 0)
	for _, element := range *allresources {
		found := false
		for _, bmap := range blocked {
			if bmap["resource"] == strconv.Itoa(element.Id) {
				found = true
			}
		}
		if !found {
			rt = append(rt, element)
		}
	}
	fmt.Println(rt[:])
	return rt[:], nil
}

func ListResourceForAdm(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return queryDBTable("id, displayname", "resourcelist", "", "", "", 0, 0)
}

func FetchResourceDetailAdm(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	rid, _ := extractParams(r, "rid")
	return queryDBTable("*", "resourcelist", "", "", "id="+rid, 0, 0)
}

func AddResource(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	resourceid, _ := extractParams(r, "resourceid")
	displayname, _ := extractParams(r, "displayname")
	restype, _ := extractParams(r, "restype")
	viewpermission, _ := extractParams(r, "viewpermission")
	bookpermission, _ := extractParams(r, "bookpermission")
	capastr, _ := extractParams(r, "capacity")
	capacity, _ := strconv.Atoi(capastr)
	currentuser, err := GetUserNamefromCookie(r)

	if err != nil {
		log.Fatal(err)
	}
	db := GetDBHandle()
	succ := false

	err = db.QueryRow("INSERT INTO resourcelist(resourceid,displayname,type,viewpermission,bookpermission,createdby,created,capacity) VALUES($1,$2,$3,$4,$5,$6,$7,$8);",
		resourceid, displayname, restype, viewpermission, bookpermission, currentuser, time.Now(), capacity).Scan(&succ)

	fmt.Println("hit AddResource")
	return succ, err
}

func EditResource(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	displayname, _ := extractParams(r, "displayname")
	restype, _ := extractParams(r, "restype")
	viewpermission, _ := extractParams(r, "viewpermission")
	bookpermission, _ := extractParams(r, "bookpermission")
	capastr, _ := extractParams(r, "capacity")

	capacity, _ := strconv.Atoi(capastr)
	fmt.Println("capastr", capastr, "capacity", capacity)
	id, _ := extractParams(r, "id")
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}

	db := GetDBHandle()
	succ := false

	err = db.QueryRow("UPDATE resourcelist SET displayname = $1, type = $2, viewpermission = $3, bookpermission = $4, createdby= $5, created = $6, capacity = $7 WHERE id = $8",
		displayname, restype, viewpermission, bookpermission, currentuser, time.Now(), capacity, id).Scan(&succ)

	fmt.Println("hit EditResource")
	return succ, err

}

func BookResource(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

		db.QueryRow("INSERT INTO resourcebooking(resource,bookedforuser,bookstart,bookend,createdby,created) VALUES($1,$2,$3,$4,$5,$6);",
			resourceid, currentuser, starttime, endtime, currentuser, time.Now())

	}

	fmt.Println("hit BookResource")
	return "OK", nil
}

func cancelBooking(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	resourceid := r.Form.Get("resourceid")
	username := r.Form.Get("username")
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	if HasPermission(currentuser, "clientadminperm") || currentuser == username {
		db := GetDBHandle()
		db.QueryRow("DELETE FROM resourcebooking WHERE resource = $1 AND bookedforuser = $2;", resourceid, username)
		fmt.Println("hit deleteBooking")
		return "OK", nil
	} else {
		fmt.Println("hit deleteBooking, but no permission")
		return "OK", errors.New("hit deleteBooking, but no permission")
	}

}

func cancelBookingAdm(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	bookingId := r.Form.Get("bookingId")
	db := GetDBHandle()
	db.QueryRow("DELETE FROM resourcebooking WHERE id = $1;", bookingId)
	fmt.Println("hit deleteBooking")
	return "OK", nil
}

func getBookingListAdm(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fmt.Println("hit getBookingListAdm")
	return queryDBTable("*", "resourcebooking", "", "", "", 100, 0)
}

func getBookingListTodayAdm(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	daybegin := bod(time.Now())
	dayend := daybegin.AddDate(0, 0, 1)
	timestr := "bookstart>='" + bod(time.Now()).Format("2006-01-02 15:04:05") + "' AND bookend<='" + dayend.Format("2006-01-02 15:04:05") + "'"

	fmt.Println("hit getBookingListTodayAdm", "timestr = ", timestr)

	return queryDBTableAdv("SELECT resourcebooking.id, resourcebooking.resource, resourcelist.displayname, resourcebooking.bookedforuser, resourcebooking.bookstart,resourcebooking.bookend FROM resourcebooking INNER JOIN resourcelist ON resourcebooking.resource=resourcelist.id WHERE " + timestr + ";")

	//return queryDBTable("*", "resourcebooking", "", "", timestr, 100, 0)
}

func ArchivePastBooking(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	return "OK", nil
}

func CreateNewresourceType(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	typename := r.Form.Get("typename")
	//displayname := r.Form.Get("displayname")
	displayname := ""
	viewpermission := r.Form.Get("viewpermission")
	bookpermission := r.Form.Get("bookpermission")
	db := GetDBHandle()
	currentuser, err := GetUserNamefromCookie(r)
	/*
		if err != nil {
			log.Fatal(err)
		}
	*/
	db.QueryRow("INSERT INTO resourcetypes(resourcetype, displayname, viewpermission, bookpermission, created,createdby) VALUES($1,$2,$3,$4,$5,$6) returning id;", typename, displayname, viewpermission, bookpermission, time.Now(), currentuser)
	fmt.Println("hit createNewRole")
	return "OK", err
}

// For Admin to create resources
func IsResourceTypeNameAvail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	typename := r.Form.Get("typename")
	exists := false
	db := GetDBHandle()
	err := db.QueryRow("SELECT exists (SELECT resourcetype FROM resourcetypes WHERE resourcetype=$1)", typename).Scan(&exists)

	return exists, err
}

func getResourceOccupiedSlot(resIds []string, dt time.Time) *[]BookingInfo {
	db := GetDBHandle()
	year, month, day := dt.Date()
	dtBegin := time.Date(year, month, day, 0, 0, 0, 0, dt.Location())
	dtEnd := time.Date(year, month, day+1, 0, 0, 0, 0, dt.Location())
	rows, err := db.Query("SELECT id, resource, bookedforuser, bookstart,bookend FROM resourcebooking WHERE resource IN $1 AND bookstart >= $2 AND bookend <= $3;", resIds, dtBegin, dtEnd)
	if err != nil {
		log.Fatal(err)
	}
	var result []BookingInfo
	for rows.Next() {
		var bookid, resid int
		var bookedforuser string
		var bookstart, bookend time.Time
		err := rows.Scan(&bookid, &resid, &bookedforuser, &bookstart, &bookend)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, BookingInfo{bookid, resid, bookedforuser, bookstart, bookend})
	}
	return &result
}

func getResBookingStatus(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	paraName := []string{"rooms", "startTime", "endTime"}
	params, _ := extractMultiParams(r, paraName)
	username, err := GetUserNamefromCookie(r)
	rt := make(map[int]ResourceInfo)
	roomsStr := params["rooms"]
	db := GetDBHandle()

	userperms := GetAllPermsofUser(username)
	flatStr := strings.Join(userperms[:], "','")
	qStr0 := "SELECT id, displayname FROM resourcelist WHERE id IN (" + roomsStr + ") AND viewpermission IN ('" + flatStr + "');"

	rows0, err0 := db.Query(qStr0)
	if err0 != nil {
		log.Fatal(err0)
	}

	for rows0.Next() {
		var id int
		var name string
		err1 := rows0.Scan(&id, &name)
		if err1 != nil {
			log.Fatal(err1)
		}
		resInfo := ResourceInfo{
			Id:          id,
			Name:        name,
			BookingStat: make([]BookingInfo, 0),
		}
		rt[id] = resInfo
	}

	qStr := "SELECT id, resource,bookedforuser,bookstart,bookend FROM resourcebooking WHERE resource IN (" + roomsStr + ")"
	//timestr := "bookstart<'" + endTime + "' OR bookend>'" + startTime + "'"
	rows, err := db.Query(qStr+" AND bookstart >= $1 AND bookend <= $2;", params["startTime"], params["endTime"])
	if err != nil {
		log.Fatal(err)
	}

	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}

	//var result []BookingInfo
	for rows.Next() {
		var bookid, resid int
		var bookedforuser string
		var bookstart, bookend time.Time
		err := rows.Scan(&bookid, &resid, &bookedforuser, &bookstart, &bookend)
		if err != nil {
			log.Fatal(err)
		}
		info := BookingInfo{bookid, resid, bookedforuser, bookstart, bookend}

		if bookedforuser != currentuser {
			info.Bookedforuser = ""
		}

		temp := rt[resid]
		temp.BookingStat = append(temp.BookingStat, info)
		rt[resid] = temp
	}
	return rt, nil

}
