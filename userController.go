package main

import (
	"fmt"
	"log"
	"net/http"
)

func addFavResforCurrentUser(w http.ResponseWriter, r *http.Request) {
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	resourceid := r.Form.Get("resourceid")
	db := GetDBHandle()

	db.QueryRow("INSERT INTO userfavresource(resource,username) VALUES($1,$2);",
		resourceid, currentuser)

	fmt.Println("hit addFavResforCurrentUser")
}

func removeFavResforCurrentUser(w http.ResponseWriter, r *http.Request) {
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	resourceid := r.Form.Get("resourceid")
	db := GetDBHandle()

	db.QueryRow("DELETE FROM userfavresource WHERE resource = $1 AND username=$2;", resourceid, currentuser)

	fmt.Println("hit removeFavResforCurrentUser")
}
