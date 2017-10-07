package main

import (
	"net/http"
	"time"
	"log"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func createUser(r *http.Request) {
	username := r.Form.Get("username")
	db:= GetDBHandle()
	db.QueryRow("INSERT INTO useraccount(username,userid,created) VALUES($1,$2,$3) returning id;", username, "", time.Now())
	fmt.Println("hit createUser")
}

func disableUser(r *http.Request) {
	username := r.Form.Get("username")
	db:= GetDBHandle()
	db.QueryRow("UPDATE useraccount SET disabled = 1 WHERE username=$1;", username)
	fmt.Println("hit disableUser")
}

func enableUser(r *http.Request) {
	username := r.Form.Get("username")
	db:= GetDBHandle()
	db.QueryRow("UPDATE useraccount SET disabled = 0 WHERE username=$1;", username)
	fmt.Println("hit enableUser")
}

func removeUser(r *http.Request) {
	//username := r.Form.Get("username")
	fmt.Println("hit removeUser")
}

func createNewRole(r *http.Request) {
	rolename := r.Form.Get("rolename")
	db:= GetDBHandle()
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	db.QueryRow("INSERT INTO rolelist(rolename,created,createdby) VALUES($1,$2,$3) returning id;", rolename, time.Now(), currentuser)
	fmt.Println("hit createNewRole")
}

func assignRoletoUser(r *http.Request) {
	db:= GetDBHandle()
	username := r.Form.Get("username")
	rolename := r.Form.Get("rolename")
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	} 
	db.QueryRow("INSERT INTO roleassignment(username,rolename,created,createdby) VALUES($1,$2,$3,$4) returning id;", username, rolename, time.Now(), currentuser)
	fmt.Println("hit assignRoletoUser")		
}

func createNewPerm(r *http.Request) {
	permissionname := r.Form.Get("permissionname")
	db:= GetDBHandle()
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	db.QueryRow("INSERT INTO permissionlist(permissionname,created,createdby) VALUES($1,$2,$3) returning id;", permissionname, time.Now(), currentuser)
	fmt.Println("hit createNewRole")
}

func assignRoletoPerm(r *http.Request) {
	db:= GetDBHandle()
	permissionname := r.Form.Get("permissionname")
	rolename := r.Form.Get("rolename")
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	} 
	db.QueryRow("INSERT INTO permissionassignment(permissionname,rolename,created,createdby) VALUES($1,$2,$3,$4) returning id;", permissionname, rolename, time.Now(), currentuser)
	fmt.Println("hit assignRoletoUser")		
}

func removeRolefromUser(r *http.Request) {
	//username := r.Form.Get("username")
	//rolename := r.Form.Get("rolename")
	fmt.Println("hit removeRolefromUser")
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
	db:= GetDBHandle()
	
	rows, err := db.Query("SELECT rolename FROM roleassignment WHERE username=$1", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
			var role string
			if err := rows.Scan(&role); err != nil {
					log.Fatal(err)
			}
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
	if err := rows.Err(); err != nil {
			log.Fatal(err)
	}
	return false
}

func createNewresourceType(r *http.Request) {
	typename := r.Form.Get("typename")
	displayname := r.Form.Get("displayname")
	viewpermission := r.Form.Get("viewpermission")
	bookpermission := r.Form.Get("bookpermission")
	db:= GetDBHandle()
	currentuser, err := GetUserNamefromCookie(r)
	if err != nil {
		log.Fatal(err)
	}
	db.QueryRow("INSERT INTO resourcetypes(resourcetype, displayname, viewpermission, bookpermission, created,createdby) VALUES($1,$2,$3,$4,$5,$6) returning id;", typename, displayname, viewpermission,bookpermission, time.Now(), currentuser)
	fmt.Println("hit createNewRole")
}
	
func makeRestrictiedHandlerbyRole(requireRole string, funcName func(*http.Request)) func (http.ResponseWriter, *http.Request) {

	return func (w http.ResponseWriter, r *http.Request) {
		userId, _ := GetUserNamefromCookie(r)
		if HasRole(userId, requireRole){
			funcName(r)
		} else {
			http.NotFound(w, r)
            return
		}
	}
}

func makeRestrictiedHandlerbyPerm(requirePerm string, funcName func(*http.Request)) func (http.ResponseWriter, *http.Request) {
	
		return func (w http.ResponseWriter, r *http.Request) {
			userId, _ := GetUserNamefromCookie(r)
			if HasPermission(userId, requirePerm){
				funcName(r)
			} else {
				http.NotFound(w, r)
				return
			}
		}
	}

func InitDBTablewithValue(){
	//db:= GetDBHandle()
	//db.QueryRow("INSERT INTO useraccount(username,userid,created) VALUES($1,$2,$3) returning id;", "sysadmin0", "", time.Now())
	//db.QueryRow("INSERT INTO rolelist(rolename,created) VALUES($1,$2);", "sysadminrole", time.Now())
	//db.QueryRow("INSERT INTO roleassignment(username, rolename, created) VALUES($1,$2,$3);", "sysadmin0", "sysadminrole", time.Now())
	//db.QueryRow("INSERT INTO permissionassignment(permissionname, rolename, created) VALUES($1,$2,$3);", "sysadminperm", "sysadminrole", time.Now())

	//db.QueryRow("INSERT INTO permissionlist(permissionname, created) VALUES($1,$2);", "clientadminperm", time.Now())
	//db.QueryRow("INSERT INTO permissionassignment(permissionname, rolename, created) VALUES($1,$2,$3);", "clientadminperm", "sysadminrole", time.Now())
}

func GetFunctionName(i interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
func AddUserPermHandler(){
	reqClientAdmin := []func(*http.Request){createUser, enableUser, disableUser, assignRoletoUser, removeRolefromUser, createNewPerm}
	reqSuperAdmin := []func(*http.Request){removeUser, createNewRole}
	for _, funcName := range reqClientAdmin{
		tokens := strings.Split(strings.ToLower(GetFunctionName(funcName)), ".")
		endPoint := "/"+ tokens[len(tokens)-1]
		fmt.Println("endPoint:", endPoint)
		http.HandleFunc(endPoint,  makeRestrictiedHandlerbyPerm("clientadminperm", funcName))
	}

	for _, funcName := range reqSuperAdmin{
		tokens := strings.Split(strings.ToLower(GetFunctionName(funcName)), ".")
		endPoint := "/"+ tokens[len(tokens)-1]
		fmt.Println("endPoint:", endPoint)
		http.HandleFunc(endPoint,  makeRestrictiedHandlerbyPerm("superadminperm", funcName))
		
	}
}