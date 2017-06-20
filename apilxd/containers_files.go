package apilxd

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)

func containerFileHandler(lx *LxdpmApi, r *http.Request) Response {
	cname := mux.Vars(r)["name"]
	path := r.FormValue("path")
	fmt.Printf("%+v",r)

	hostname := getHostnameFromContainername(lx,cname)
	switch r.Method {
	case "GET":
			res,headers,_ := containerFileGet(hostname[0][0].(string),cname,path)
			responseType := getResponseType(res)
			if responseType == "sync" {
				endpointResponse,_ := parseResponseRawToSync(res)
				return &endpointResponse
			}else if responseType == "error"{
				errorResp := parseErrorResponse(res)
				return &errorResp
			}else {
				return createFileResponse(string(res),string(headers),path,r)
			}
	case "POST":
			res,_ := containerFilePost(hostname[0][0].(string),cname,path,r)
			responseType := getResponseType(res)
			if responseType == "sync" {
				endpointResponse,_ := parseResponseRawToSync(res)
				return &endpointResponse
			}else if responseType == "error"{
				errorResp := parseErrorResponse(res)
				return &errorResp
			}
	case "DELETE":
			res,_ := containerFileDelete(hostname[0][0].(string),cname,path,r)
			responseType := getResponseType(res)
			if responseType == "sync" {
				endpointResponse,_ := parseResponseRawToSync(res)
				return &endpointResponse
			}else if responseType == "error"{
				errorResp := parseErrorResponse(res)
				return &errorResp
			}
	default:
		return NotFound
	}
	
	return NotFound

}