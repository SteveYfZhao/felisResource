package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

//TODO: Add remove and edit logic for everything

func createUser(w http.ResponseWriter, r *http.Request) {
	username := r.Form.Get("username")
	db := GetDBHandle()
	db.QueryRow("INSERT INTO useraccount(username,userid,created) VALUES($1,$2,$3) returning id;", username, "", time.Now())
	fmt.Println("hit createUser")
}

func createUserbySignup(w http.ResponseWriter, r *http.Request) {
	username := r.Form.Get("username")
	email := r.Form.Get("email")
	salt, _ := GenerateRandomString(128)
	hasher := sha1.New()
	hasher.Write([]byte(email + r.Form.Get("password") + salt))
	passwordHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	db := GetDBHandle()
	db.QueryRow("INSERT INTO useraccount(username,userid,created,email,salt,passwordhash) VALUES($1,$2,$3,$4,$5,$6) returning id;", username, "", time.Now(), email, salt, passwordHash)
	fmt.Println("hit createUser")
}

// test if user can login with password
func LoginPW(user string, pass string) bool {
	var pwHash, email, salt string
	db := GetDBHandle()
	row := db.QueryRow("SELECT email,salt,passwordhash FROM useraccount WHERE username=$1;", user)
	err := row.Scan(&email, &salt, &pwHash)

	if err == nil {
		hasher := sha1.New()
		hasher.Write([]byte(email + pass + salt))
		passwordHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		if passwordHash == pwHash {
			return true
		}
	}
	return false
}

func disableUser(w http.ResponseWriter, r *http.Request) {
	username := r.Form.Get("username")
	db := GetDBHandle()
	db.QueryRow("UPDATE useraccount SET disabled = 1 WHERE username=$1;", username)
	fmt.Println("hit disableUser")
}

func enableUser(w http.ResponseWriter, r *http.Request) {
	username := r.Form.Get("username")
	db := GetDBHandle()
	db.QueryRow("UPDATE useraccount SET disabled = 0 WHERE username=$1;", username)
	fmt.Println("hit enableUser")
}

func removeUser(w http.ResponseWriter, r *http.Request) {
	//username := r.Form.Get("username")
	fmt.Println("hit removeUser")
}

func createNewRole(w http.ResponseWriter, r *http.Request) {
	rolename := r.Form.Get("rolename")
	db := GetDBHandle()
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	db.QueryRow("INSERT INTO rolelist(rolename,created,createdby) VALUES($1,$2,$3) returning id;", rolename, time.Now(), currentuser)
	fmt.Println("hit createNewRole")
}

func assignRoletoUser(w http.ResponseWriter, r *http.Request) {
	db := GetDBHandle()
	username := r.Form.Get("username")
	rolename := r.Form.Get("rolename")
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	db.QueryRow("INSERT INTO roleassignment(username,rolename,created,createdby) VALUES($1,$2,$3,$4) returning id;", username, rolename, time.Now(), currentuser)
	fmt.Println("hit assignRoletoUser")
}

func createNewPerm(w http.ResponseWriter, r *http.Request) {
	permissionname := r.Form.Get("permissionname")
	db := GetDBHandle()
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	db.QueryRow("INSERT INTO permissionlist(permissionname,created,createdby) VALUES($1,$2,$3) returning id;", permissionname, time.Now(), currentuser)
	fmt.Println("hit createNewRole")
}

func assignRoletoPerm(w http.ResponseWriter, r *http.Request) {
	db := GetDBHandle()
	permissionname := r.Form.Get("permissionname")
	rolename := r.Form.Get("rolename")
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	db.QueryRow("INSERT INTO permissionassignment(permissionname,rolename,created,createdby) VALUES($1,$2,$3,$4) returning id;", permissionname, rolename, time.Now(), currentuser)
	fmt.Println("hit assignRoletoUser")
}

func removeRolefromUser(w http.ResponseWriter, r *http.Request) {
	db := GetDBHandle()
	username := r.Form.Get("username")
	rolename := r.Form.Get("rolename")
	db.QueryRow("DELETE FROM roleassignment WHERE username = $1 AND rolename =$2);", username, rolename)
	fmt.Println("hit removeRolefromUser")
}

func removeRolefromPerm(w http.ResponseWriter, r *http.Request) {
	db := GetDBHandle()
	permissionname := r.Form.Get("permissionname")
	rolename := r.Form.Get("rolename")
	db.QueryRow("DELETE FROM permissionassignment WHERE permissionname = $1 AND rolename =$2);", permissionname, rolename)
	fmt.Println("hit removeRolefromPerm")
}

func HasRole(username string, rolename string) bool {
	fmt.Println("hit HasRole", username, rolename)
	db := GetDBHandle()
	exists := false
	//username := r.Form.Get("username")
	//rolename := r.Form.Get("rolename")
	err := db.QueryRow("SELECT exists (SELECT rolename FROM roleassignment WHERE username=$1 AND rolename=$2)", username, rolename).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("HasRole", exists)
	return exists
}

func HasPermission(username string, permission string) bool {
	fmt.Println("hit HasPermission", username, permission)
	db := GetDBHandle()
	roles := GetAllRolesOfUser(username)
	/*
		rows, err := db.Query("SELECT rolename FROM roleassignment WHERE username=$1", username)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
	*/
	for _, role := range roles {
		//var role string
		fmt.Println("role", role, permission)
		hasPerm := false
		err := db.QueryRow("SELECT exists (SELECT permissionname FROM permissionassignment WHERE permissionname=$1 AND rolename=$2)", permission, role).Scan(&hasPerm)
		if err != nil {
			log.Fatal(err)
		}
		if hasPerm {
			return hasPerm
		}
	}
	return false
}

func GetAllRolesOfUser(username string) []string {
	db := GetDBHandle()

	rows, err := db.Query("SELECT rolename FROM roleassignment WHERE username=$1", username)
	if err != nil {
		log.Fatal(err)
	}

	roles := rowsToStringSlice(rows)
	return roles
}

func GetAllPermsOfRole(rolename string) []string {
	db := GetDBHandle()
	rows, err := db.Query("SELECT permissionname FROM permissionassignment WHERE rolename=$1", rolename)
	if err != nil {
		log.Fatal(err)
	}
	perms := rowsToStringSlice(rows)
	return perms
}

func GetAllPermsofUser(username string) []string {

	var perms []string
	roles := GetAllRolesOfUser(username)
	for _, role := range roles {
		roleperms := GetAllPermsOfRole(role)
		perms = append(perms, roleperms...)
	}
	return perms
}

func createNewresourceType(w http.ResponseWriter, r *http.Request) {
	typename := r.Form.Get("typename")
	displayname := r.Form.Get("displayname")
	viewpermission := r.Form.Get("viewpermission")
	bookpermission := r.Form.Get("bookpermission")
	db := GetDBHandle()
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	db.QueryRow("INSERT INTO resourcetypes(resourcetype, displayname, viewpermission, bookpermission, created,createdby) VALUES($1,$2,$3,$4,$5,$6) returning id;", typename, displayname, viewpermission, bookpermission, time.Now(), currentuser)
	fmt.Println("hit createNewRole")
}

func InitDBTablewithValue() {
	//db:= GetDBHandle()
	//db.QueryRow("INSERT INTO useraccount(username,userid,created) VALUES($1,$2,$3) returning id;", "sysadmin0", "", time.Now())
	//db.QueryRow("INSERT INTO rolelist(rolename,created) VALUES($1,$2);", "sysadminrole", time.Now())
	//db.QueryRow("INSERT INTO roleassignment(username, rolename, created) VALUES($1,$2,$3);", "sysadmin0", "sysadminrole", time.Now())
	//db.QueryRow("INSERT INTO permissionassignment(permissionname, rolename, created) VALUES($1,$2,$3);", "sysadminperm", "sysadminrole", time.Now())

	//db.QueryRow("INSERT INTO permissionlist(permissionname, created) VALUES($1,$2);", "clientadminperm", time.Now())
	//db.QueryRow("INSERT INTO permissionassignment(permissionname, rolename, created) VALUES($1,$2,$3);", "clientadminperm", "sysadminrole", time.Now())
}
