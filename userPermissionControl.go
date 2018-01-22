package main

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

//TODO: Add remove and edit logic for everything

//Role and Permission design
/*
The system only have four level of permissions:
1 basic user
2 pro user
3 admin
4 super admin

each perm contains all power of the previous one and some more.

by default there should be 4 roles match the permissions. admin can create more roles.


*/

/*
func createUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	username := r.Form.Get("username")
	db := GetDBHandle()
	db.QueryRow("INSERT INTO useraccount(username,userid,created) VALUES($1,$2,$3) returning id;", username, "", time.Now())
	fmt.Println("hit createUser")
	return "OK", nil
}
*/

func createUserbySignup(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fmt.Println("hit createUsersignup")
	err := r.ParseForm()
	if err == nil {
		username := r.PostForm["username"][0]
		email := r.PostForm["email"][0]
		password := r.PostForm["password"][0]

		fmt.Println("username:", username)
		fmt.Println("email:", email)
		fmt.Println("password:", password)
		// test if userid or email exists
		exists := false
		//username := r.Form.Get("username")
		//rolename := r.Form.Get("rolename")
		db := GetDBHandle()
		err := db.QueryRow("SELECT exists (SELECT username FROM useraccount WHERE username=$1 OR email=$2)", username, email).Scan(&exists)
		if err != nil {
			log.Fatal(err)
		}
		if !exists {
			salt, _ := GenerateRandomString(128)
			hasher := sha1.New()
			hasher.Write([]byte(email + password + salt))
			passwordHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
			salt1, _ := GenerateRandomString(128)
			hasher1 := sha1.New()
			hasher1.Write([]byte(email + salt1))
			verificationhash := base64.URLEncoding.EncodeToString(hasher1.Sum(nil))

			db.QueryRow("INSERT INTO useraccount(username,verificationhash,created,createdby,email,salt,passwordhash,disabled) VALUES($1,$2,$3,$4,$5,$6,$7,$8) returning id;", username, verificationhash, time.Now(), "self-registered", email, salt, passwordHash, true)
			fmt.Println("complete createUser")
			scheme, hostname := GetRootURL(r)
			/*

				fmt.Println("r.URL")
				fmt.Println(r.URL.String())
				fmt.Println("r.host")
				fmt.Println(r.Host)
				fmt.Println("r.URL.Hostname")
				fmt.Println(r.URL.Hostname())
			*/
			redirectURL := scheme + "://" + hostname + ":" + FrontEndPort + "/login"
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return "OK", nil
		}
		fmt.Println("user/email exists. Cannot createUser")
		return "OK", errors.New("user/email exists. Cannot createUser")

	}
	return nil, err
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

func disableUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	username := r.Form.Get("username")
	db := GetDBHandle()
	db.QueryRow("UPDATE useraccount SET disabled = 1 WHERE username=$1;", username)
	fmt.Println("hit disableUser")
	return "OK", nil
}

func enableUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	username := r.Form.Get("username")
	db := GetDBHandle()
	db.QueryRow("UPDATE useraccount SET disabled = 0 WHERE username=$1;", username)
	fmt.Println("hit enableUser")
	return "OK", nil
}

func removeUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	//username := r.Form.Get("username")
	fmt.Println("hit removeUser")
	return "OK", nil
}

func createNewRole(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	rolename := r.Form.Get("rolename")
	db := GetDBHandle()
	currentuser, err := GetUserNamefromCookie(r)
	/*
		if err != nil {
			log.Fatal(err)
		}
	*/
	db.QueryRow("INSERT INTO rolelist(rolename,created,createdby) VALUES($1,$2,$3) returning id;", rolename, time.Now(), currentuser)
	fmt.Println("hit createNewRole")
	return "OK", err
}

func assignRoletoUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	db := GetDBHandle()
	username := r.Form.Get("username")
	rolename := r.Form.Get("rolename")
	currentuser, err := GetUserNamefromCookie(r)
	/*
		if err != nil {
			log.Fatal(err)
		}
	*/
	db.QueryRow("INSERT INTO roleassignment(username,rolename,created,createdby) VALUES($1,$2,$3,$4) returning id;", username, rolename, time.Now(), currentuser)
	fmt.Println("hit assignRoletoUser")
	return "OK", err
}

func createNewPerm(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	permissionname := r.Form.Get("permissionname")
	db := GetDBHandle()
	currentuser, err := GetUserNamefromCookie(r)
	/*
		if err != nil {
			log.Fatal(err)
		}
	*/
	db.QueryRow("INSERT INTO permissionlist(permissionname,created,createdby) VALUES($1,$2,$3) returning id;", permissionname, time.Now(), currentuser)
	fmt.Println("hit createNewRole")
	return "OK", err
}

func deletePerm(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	permissionname := r.Form.Get("permissionname")
	if len(strings.TrimSpace(permissionname)) != 0 {
		db := GetDBHandle()
		db.QueryRow("DELETE FROM permissionlist WHERE permissionname = $1;", permissionname)
		clearAssignmentForPermCore(permissionname)
		fmt.Println("hit deletePerm")
		return "OK", nil
	}
	return nil, errors.New("permission name is empty")
}

func clearAssignmentForPermCore(permissionname string) {

	db := GetDBHandle()
	db.QueryRow("DELETE FROM permissionassignment WHERE permissionname = $1;", permissionname)
}

func assignRoletoPerm(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	db := GetDBHandle()
	permissionname := r.Form.Get("permissionname")
	rolename := r.Form.Get("rolename")
	currentuser, err := GetUserNamefromCookie(r)
	/*
		if err != nil {
			log.Fatal(err)
		}
	*/
	db.QueryRow("INSERT INTO permissionassignment(permissionname,rolename,created,createdby) VALUES($1,$2,$3,$4) returning id;", permissionname, rolename, time.Now(), currentuser)
	fmt.Println("hit assignRoletoUser")
	return "OK", err
}

func removeRolefromUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	db := GetDBHandle()
	username := r.Form.Get("username")
	rolename := r.Form.Get("rolename")
	db.QueryRow("DELETE FROM roleassignment WHERE username = $1 AND rolename =$2);", username, rolename)
	fmt.Println("hit removeRolefromUser")
	return "OK", nil
}

func removeRolefromPerm(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	db := GetDBHandle()
	permissionname := r.Form.Get("permissionname")
	rolename := r.Form.Get("rolename")
	db.QueryRow("DELETE FROM permissionassignment WHERE permissionname = $1 AND rolename =$2);", permissionname, rolename)
	fmt.Println("hit removeRolefromPerm")
	return "OK", nil
}

type UserBaseInfo struct {
	Username          string
	CommonPermissions []string
}

func UserBasicInfo(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	cookieUsername, _ := GetUserNamefromCookie(r)
	userBasicPerms := []string{"canAccess", "basicClient", "basicAdmin"}

	rt := UserBaseInfo{cookieUsername, userBasicPerms}

	/*
		rvalues, err := json.Marshal(rt)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		fmt.Println("rt = ", rt)
		fmt.Println("rvalues = ", rvalues)
		fmt.Fprintf(w, string(rvalues))
	*/
	fmt.Println("hit UserBasicInfo")
	return rt, nil
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

	roles := RowsToStringSlice(rows)
	return roles
}

func GetAllPermsOfRole(rolename string) []string {
	db := GetDBHandle()
	rows, err := db.Query("SELECT permissionname FROM permissionassignment WHERE rolename=$1", rolename)
	if err != nil {
		log.Fatal(err)
	}
	perms := RowsToStringSlice(rows)
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

func InitDBTablewithValue() {
	//db:= GetDBHandle()
	//db.QueryRow("INSERT INTO useraccount(username,userid,created) VALUES($1,$2,$3) returning id;", "sysadmin0", "", time.Now())
	//db.QueryRow("INSERT INTO rolelist(rolename,created) VALUES($1,$2);", "sysadminrole", time.Now())
	//db.QueryRow("INSERT INTO roleassignment(username, rolename, created) VALUES($1,$2,$3);", "sysadmin0", "sysadminrole", time.Now())
	//db.QueryRow("INSERT INTO permissionassignment(permissionname, rolename, created) VALUES($1,$2,$3);", "sysadminperm", "sysadminrole", time.Now())

	//db.QueryRow("INSERT INTO permissionlist(permissionname, created) VALUES($1,$2);", "clientadminperm", time.Now())
	//db.QueryRow("INSERT INTO permissionassignment(permissionname, rolename, created) VALUES($1,$2,$3);", "clientadminperm", "sysadminrole", time.Now())
}
