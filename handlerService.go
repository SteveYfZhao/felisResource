package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

type EndPoint struct {
	Function   func(http.ResponseWriter, *http.Request) (interface{}, error)
	Permission string
}

var endPointList = []EndPoint{
	EndPoint{UserBasicInfo, "Public"},
	EndPoint{createUserbySignup, "Public"},

	// user section
	//EndPoint{createUser, "clientadminperm"},
	EndPoint{enableUser, "clientadminperm"},
	EndPoint{disableUser, "clientadminperm"},
	EndPoint{ListUsers, "Public"},
	EndPoint{FindUser, "Public"},

	//Role and perm section
	EndPoint{assignRoletoUser, "Public"},
	EndPoint{removeRolefromUser, "Public"},
	EndPoint{createNewPerm, "Public"},
	EndPoint{assignRoletoPerm, "Public"},
	EndPoint{removeRolefromPerm, "Public"},
	EndPoint{deletePerm, "Public"},

	//restrictions
	EndPoint{CreateRestrcitionType, "clientadminperm"},
	EndPoint{CreateRestrcition, "clientadminperm"},
	EndPoint{UpdateRestrcition, "clientadminperm"},
	EndPoint{RemoveRestrcition, "clientadminperm"},
	EndPoint{GetRestrictionValuesCL, "clientadminperm"},

	//resource and bookings
	EndPoint{ArchivePastBooking, "clientadminperm"},
	EndPoint{AddAvailTimePlan, "clientadminperm"},

	// superadmin
	EndPoint{removeUser, "superadminperm"},
	EndPoint{createNewRole, "superadminperm"},
	EndPoint{CreateNewresourceType, "superadminperm"},
	EndPoint{IsResourceTypeNameAvail, "superadminperm"},
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func AddUserPermHandler() {

	for _, ep := range endPointList {
		tokens := strings.Split(strings.ToLower(GetFunctionName(ep.Function)), ".")
		endPoint := "/" + tokens[len(tokens)-1]
		fmt.Println("endPoint:", endPoint)
		http.HandleFunc(endPoint, makeRestrictiedHandlerbyPerm(ep.Permission, ep.Function))
	}
}

type HandleResponse struct {
	Data  interface{}
	Error error
}

func processFuncResp(w http.ResponseWriter, r *http.Request, rt interface{}, err error) {
	resp := HandleResponse{nil, nil}

	log.Print("preprocess Data", rt, "preprocess error", err)

	if err == nil && rt != nil {
		resp.Data = rt
	}

	if err != nil && IsDEV == true {
		resp.Data = rt
		resp.Error = err
		log.Fatal(err)
	}

	log.Print("postprocess Data", rt, "postprocess error", err)

	if resp.Data != nil || resp.Error != nil {
		rvalues, err := json.Marshal(resp)
		log.Print("Marshal Data", rvalues)
		fmt.Fprintf(w, string(rvalues))
		if err != nil {
			log.Fatal(err)
		}
	}
}

/*
func makePublicHandler(funcName func(http.ResponseWriter, *http.Request) (interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// preprocess w and r
		w, r = preprocessRequestAndReponse(w, r)
		rt, err := funcName(w, r)
		processFuncResp(w, r, rt, err)
	}
}
*/
func makeRestrictiedHandlerbyPerm(requirePerm string, funcName func(http.ResponseWriter, *http.Request) (interface{}, error)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := GetUserNamefromCookie(r)
		restrictionNotEmpty := (len(strings.TrimSpace(requirePerm)) != 0)
		restrictionPublic := (requirePerm == "Public")
		//restrictionAnon := (requirePerm == "Anonymous")
		if restrictionNotEmpty && (restrictionPublic || HasPermission(userID, requirePerm)) {
			w, r = preprocessRequestAndReponse(w, r)
			rt, err := funcName(w, r)
			log.Print("raw reply from kernel", rt)

			processFuncResp(w, r, rt, err)
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

/*
func makeRestrictiedHandlerbyRole(requireRole string, funcName func(http.ResponseWriter, *http.Request) (interface{}, error)) func(http.ResponseWriter, *http.Request) {
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
*/
