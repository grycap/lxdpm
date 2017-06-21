package apilxd

import (
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
	"github.com/lxc/lxd/shared/api"
	"github.com/gorilla/mux"
)


var profileCmd = Command{
	name: "profiles/{name}",
	get: profileGet,
	put: profilePut,
	delete: profileDelete,
	post: profilePost,
//	patch: profilePatch
}

func getHostnameFromProfileName(lx *LxdpmApi, name string) [][]interface{} {
	inargs := []interface{}{}
	outargs := []interface{}{"name"}

	result, err := dbQueryScan(lx.db, `SELECT H.name FROM hosts H, profiles P where P.host_id = H.id and P.name='`+name+`';`,inargs,outargs )
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func profileGet(lx *LxdpmApi, r *http.Request) Response {
	name := mux.Vars(r)["name"]

	hostname := getHostnameFromProfileName(lx,name)

	resp,_ := doProfileGet(hostname[0][0].(string),name)

	endpointResponse,_ := parseResponseRawToSync(resp)

	return &endpointResponse
}

func parseMetadataFromProfileResponse(input []byte) LxdResponseRaw {
	var resp = LxdResponseRaw{}
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    return resp
}

func profilePut(lx *LxdpmApi,  r *http.Request) Response {
	pname := mux.Vars(r)["name"]
	hostname := getHostnameFromProfileName(lx,pname)
	req := api.ProfilePut{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)
	res,_ := doProfilePut(hostname[0][0].(string),req,pname)
	responseType := getResponseType(res)
	if responseType == "sync" {
		endpointResponse,_ := parseResponseRawToSync(res)
		return &endpointResponse
	}else {
		errorResp := parseErrorResponse(res)
		return &errorResp
	}

}

func profileDelete(lx *LxdpmApi,  r *http.Request) Response {
	pname := mux.Vars(r)["name"]
	hostname := getHostnameFromProfileName(lx,pname)
	res,_ := doProfileDelete(hostname[0][0].(string),pname)
	responseType := getResponseType(res)
	if responseType == "sync" {
		endpointResponse,_ := parseResponseRawToSync(res)
		return &endpointResponse
	}else {
		errorResp := parseErrorResponse(res)
		return &errorResp
	}
}

func profilePost(lx *LxdpmApi,  r *http.Request) Response {
	pname := mux.Vars(r)["name"]
	hostname := getHostnameFromProfileName(lx,pname)

	req := api.ProfilePost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)
	res,_ := doProfilePost(hostname[0][0].(string),req,pname)
	responseType := getResponseType(res)
	if responseType == "sync" {
		endpointResponse,_ := parseResponseRawToSync(res)
		return &endpointResponse
	}else {
		errorResp := parseErrorResponse(res)
		return &errorResp
	}
}