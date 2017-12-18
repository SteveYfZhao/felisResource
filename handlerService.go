package main

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func AddUserPermHandler() {

	publicEndPoints := []func(http.ResponseWriter, *http.Request){
		UserBasicInfo,
		createUserbySignup,
	}
	reqClientAdmin := []func(http.ResponseWriter, *http.Request){
		createUser,
		enableUser,
		disableUser,
		assignRoletoUser,
		removeRolefromUser,
		createNewPerm,
		CreateRestrcitionType,
		CreateRestrcition,
		UpdateRestrcition,
		RemoveRestrcition,
		GetRestrictionValuesCL,
	}
	reqSuperAdmin := []func(http.ResponseWriter, *http.Request){
		removeUser,
		createNewRole,
	}

	for _, funcName := range publicEndPoints {
		tokens := strings.Split(strings.ToLower(GetFunctionName(funcName)), ".")
		endPoint := "/" + tokens[len(tokens)-1]
		fmt.Println("endPoint:", endPoint)
		http.HandleFunc(endPoint, makePublicHandler(funcName))
	}

	for _, funcName := range reqClientAdmin {
		tokens := strings.Split(strings.ToLower(GetFunctionName(funcName)), ".")
		endPoint := "/" + tokens[len(tokens)-1]
		fmt.Println("endPoint:", endPoint)
		http.HandleFunc(endPoint, makeRestrictiedHandlerbyPerm("clientadminperm", funcName))
	}

	for _, funcName := range reqSuperAdmin {
		tokens := strings.Split(strings.ToLower(GetFunctionName(funcName)), ".")
		endPoint := "/" + tokens[len(tokens)-1]
		fmt.Println("endPoint:", endPoint)
		http.HandleFunc(endPoint, makeRestrictiedHandlerbyPerm("superadminperm", funcName))

	}
}

func makePublicHandler(funcName func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// preprocess w and r
		w, r = preprocessRequestAndReponse(w, r)
		funcName(w, r)

	}
}

func makeRestrictiedHandlerbyRole(requireRole string, funcName func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := GetUserNamefromCookie(r)
		if HasRole(userID, requireRole) {
			w, r = preprocessRequestAndReponse(w, r)
			funcName(w, r)
		} else {
			http.NotFound(w, r)
			return
		}
	}
}

func makeRestrictiedHandlerbyPerm(requirePerm string, funcName func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := GetUserNamefromCookie(r)
		if HasPermission(userID, requirePerm) {
			w, r = preprocessRequestAndReponse(w, r)
			funcName(w, r)
		} else {
			http.NotFound(w, r)
			return
		}
	}
}

func preprocessRequestAndReponse(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	return w, r
}
