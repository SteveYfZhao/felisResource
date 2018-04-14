package main

import (
	"encoding/json"
	"errors"
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
	EndPoint{createUserbySignup, "Anonymous"},

	// user section
	//EndPoint{createUser, "clientadminperm"},
	EndPoint{enableUser, "clientadminperm"},
	EndPoint{disableUser, "clientadminperm"},
	EndPoint{ListUsers, "clientadminperm"},
	EndPoint{FindUser, "clientadminperm"},
	EndPoint{GetUserDetails, "clientadminperm"},

	//Role and perm section
	EndPoint{assignRoletoUser, "clientadminperm"},
	EndPoint{removeRolefromUser, "clientadminperm"},
	EndPoint{createNewPerm, "clientadminperm"},
	EndPoint{assignRoletoPerm, "clientadminperm"},
	EndPoint{removeRolefromPerm, "clientadminperm"},
	EndPoint{deletePerm, "clientadminperm"},

	//restrictions
	EndPoint{CreateRestrcitionType, "clientadminperm"},
	EndPoint{CreateRestrcition, "clientadminperm"},
	EndPoint{UpdateRestrcition, "clientadminperm"},
	EndPoint{RemoveRestrcition, "clientadminperm"},
	EndPoint{GetRestrictionValuesCL, "clientadminperm"},

	//resource and bookings
	EndPoint{ArchivePastBooking, "clientadminperm"},
	EndPoint{AddAvailTimePlan, "clientadminperm"},
	EndPoint{AddResource, "clientadminperm"},
	EndPoint{BookResource, "Public"},
	EndPoint{getAllResourceForUser, "Public"},
	EndPoint{getResBookingStatus, "Public"},
	EndPoint{GetBookableResForUserAtTime, "Public"},
	EndPoint{getAllAvailResource, "clientadminperm"},
	EndPoint{ListResourceForAdm, "clientadminperm"},
	EndPoint{FetchResourceDetailAdm, "clientadminperm"},
	EndPoint{EditResource, "clientadminperm"},
	EndPoint{getBookingListAdm, "clientadminperm"},
	EndPoint{getBookingListTodayAdm, "clientadminperm"},
	EndPoint{cancelBookingAdm, "clientadminperm"},

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

	//log.Print("preprocess Data", rt, "preprocess error", err)

	if err == nil && rt != nil {
		resp.Data = rt
	}

	if err != nil && IsDEV == true {
		resp.Data = rt
		resp.Error = err
		//log.Fatal(err)
		log.Print(err)
	}

	//log.Print("postprocess Data", rt, "postprocess error", err)

	if resp.Data != nil || resp.Error != nil {
		rvalues, err := json.Marshal(resp)
		//log.Print("Marshal Data", rvalues)
		fmt.Fprintf(w, string(rvalues))
		if err != nil {
			//log.Fatal(err)
			log.Print(err)
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

	// Need to ensure request have application/x-www-form-urlencoded header on the client side. Otherwise the server will not get the posted data.
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			log.Print("Received preflight")
			if r.Header.Get("Origin") == ClientURL {
				log.Print("preflight origin match")
				allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
				w.Header().Set("Access-Control-Allow-Origin", ClientURL)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
				w.Header().Set("Access-Control-Expose-Headers", "Authorization")
				fmt.Fprintf(w, "")
			}
		} else {
			log.Print("Enter handler, Request: ", r)
			restrictionNotEmpty := (len(strings.TrimSpace(requirePerm)) != 0)
			restrictionAnon := (requirePerm == "Anonymous")
			restrictionPublic := (requirePerm == "Public")
			userID := ""
			var err error

			if !restrictionAnon {
				userID, err = GetUserNamefromCookie(r)
				if err == nil && !IsEmptyStr(userID) {
					log.Print("Got user name, ", userID)
				} else {
					log.Print("Failed to get userid.")
				}
			}

			if restrictionNotEmpty && (restrictionAnon || (err == nil && !IsEmptyStr(userID) && (restrictionPublic || HasPermission(userID, requirePerm)))) {
				w, r = preprocessRequestAndReponse(w, r)
				perr := r.ParseForm()
				log.Print("try parse request parameters", perr)
				rt, err := funcName(w, r)
				log.Print("raw reply from kernel", rt, "error", err)
				processFuncResp(w, r, rt, err)
			} else {
				log.Print("Not enough permission")
				http.NotFound(w, r)
			}
		}
	}
}

func preprocessRequestAndReponse(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	w.Header().Set("Access-Control-Allow-Origin", ClientURL)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")

	return w, r
}

func extractParams(r *http.Request, param string) (string, error) {
	if r.Method == "POST" {
		paraMap, err := MapURLEncodedPostParams(r)

		return paraMap[param], err
	}

	if r.Method == "GET" {
		return r.Form.Get(param), nil
	}

	return "", errors.New("Unsupported request type")
}

func extractMultiParams(r *http.Request, params []string) (map[string]string, error) {
	if r.Method == "POST" {
		paraMap, err := MapURLEncodedPostParams(r)

		return paraMap, err
	}

	if r.Method == "GET" {
		var rt = make(map[string]string)
		for _, val := range params {
			rt[val] = r.Form.Get(val)
		}
		return rt, nil
	}

	return make(map[string]string), errors.New("Unsupported request type")
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
