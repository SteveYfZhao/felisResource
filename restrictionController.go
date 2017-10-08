package main

import (
	"net/http"
	"time"
	"log"
	"fmt"

)

func CreateRestrcitionType(r *http.Request){	
	typename := r.Form.Get("typename")
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	db:= GetDBHandle()
	db.QueryRow("INSERT INTO restrictiontypes(type,createdby,created) VALUES($1,$2,$3);", typename, currentuser, time.Now())

	fmt.Println("hit createGlobalRestrcition")

}

func CreateRestrcition(r *http.Request){
	restrictionid := r.Form.Get("restrictionid")
	restrictiontype := r.Form.Get("restrictiontype")
	restrictionvalue := r.Form.Get("restrictionvalue")
	resourcetype := r.Form.Get("resourcetype")
	resource := r.Form.Get("resource")
	restag := r.Form.Get("restag")
	restagvalue := r.Form.Get("restagvalue")
	userperm := r.Form.Get("userperm")
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	db:= GetDBHandle()
	db.QueryRow("INSERT INTO restrictions(restrictionid,restrictiontype,restrictionvalue,resourcetype,resource,restag,restagvalue,userperm,createdby,created) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);", 
	restrictionid, restrictiontype, restrictionvalue, resourcetype, resource, restag, restagvalue, userperm, currentuser, time.Now())
	

	fmt.Println("hit CreateRestrcition")

}

func UpdateRestrcition(r *http.Request){
	restrictionid := r.Form.Get("restrictionid")
	restrictiontype := r.Form.Get("restrictiontype")
	restrictionvalue := r.Form.Get("restrictionvalue")
	resourcetype := r.Form.Get("resourcetype")
	resource := r.Form.Get("resource")
	restag := r.Form.Get("restag")
	restagvalue := r.Form.Get("restagvalue")
	userperm := r.Form.Get("userperm")
	id := r.Form.Get("id")
	db:= GetDBHandle()
	db.QueryRow("UPDATE restrictions SET restrictionid=$1,restrictiontype=$2,restrictionvalue=$3,resourcetype=$4,resource=$5,restag=$6,restagvalue=$7,userperm=$8) WHERE id = &9;", 
	restrictionid, restrictiontype, restrictionvalue, resourcetype, resource, restag, restagvalue, userperm, id)
	
	fmt.Println("hit CreateRestrcition")
}

func RemoveRestriction (r *http.Request){
	id := r.Form.Get("id")
	db:= GetDBHandle()
	db.QueryRow("DELETE FROM restrictions WHERE id = $1);", id)	
	fmt.Println("hit RemoveRestriction")
}
