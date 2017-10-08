package main

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"fmt"
)


func GetFunctionName(i interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func AddUserPermHandler(){
	reqClientAdmin := []func(http.ResponseWriter, *http.Request) {
		createUser, 
		enableUser, 
		disableUser, 
		assignRoletoUser, 
		removeRolefromUser, 
		createNewPerm,
		CreateRestrcitionType,
		CreateRestrcition,
		UpdateRestrcition,
		UpdateRestrcition,
		GetRestrictionValuesCL,
	}
	reqSuperAdmin := []func(http.ResponseWriter, *http.Request){
		removeUser, 
		createNewRole,
	}
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
	
	func makeRestrictiedHandlerbyPerm(requirePerm string, funcName func(http.ResponseWriter,*http.Request)) func (http.ResponseWriter, *http.Request) {
		
			return func (w http.ResponseWriter, r *http.Request) {
				userId, _ := GetUserNamefromCookie(r)
				if HasPermission(userId, requirePerm){
					funcName(w,r)
				} else {
					http.NotFound(w, r)
					return
				}
			}
	}