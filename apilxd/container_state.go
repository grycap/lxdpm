package apilxd

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/lxc/lxd/shared/api"
	"github.com/gorilla/mux"
)

func containerState(lx *LxdpmApi, r *http.Request) Response {
	cname := mux.Vars(r)["name"]

	hostname := getHostnameFromContainername(lx,cname)

	res,_ := doContainerStateGet(hostname[0][0].(string),cname)
	responseType := operationOrError(res)
	if responseType == "operation" {
		endpointResponse,_ := parseResponseRawToSync(res)
		return &endpointResponse
		}else{
		errorResp := parseErrorResponse(res)
		return &errorResp
	}
}

func containerStatePut(lx *LxdpmApi,  r *http.Request) Response {
	var endpointResponse api.Operation
	var endpointUrl string
	cname := mux.Vars(r)["name"]
	hostname := getHostnameFromContainername(lx,cname)
	req := api.ContainerStatePut{}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)

	res,_ := doContainerStatePut(hostname[0][0].(string),req,cname)

	responseType := operationOrError(res)
	if responseType == "operation" {
		endpointUrl,endpointResponse,_ = parseResponseRawToOperation(res)
		return OperationResponse(endpointUrl,&endpointResponse)
		}else{
		errorResp := parseErrorResponse(res)
		return &errorResp
		}
}