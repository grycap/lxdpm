package apilxd

import (
	"net/http"
	"encoding/json"
	//"bytes"
	"fmt"
	"github.com/lxc/lxd/shared/api"
	"github.com/gorilla/mux"
)

func snapshotsGet(lx *LxdpmApi, r *http.Request) Response {
	name := mux.Vars(r)["name"]
	snap := mux.Vars(r)["snapshotName"]

	fmt.Println("Llego al endpoint")

	hostname := getHostnameFromContainername(lx,name)

	res,_ := doSnapshotGet(name, hostname[0][0].(string),snap)
	
	endpointResponse,_ := parseResponseRawToSync(res)

	return &endpointResponse 
}

func snapshotsPost(lx *LxdpmApi, r *http.Request) Response {
	var endpointResponse api.Operation
	var endpointUrl string
	cname := mux.Vars(r)["name"]
	snap := mux.Vars(r)["snapshotName"]
	hostname := getHostnameFromContainername(lx,cname)
	req := api.ContainerSnapshotsPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	res,_ := doSnapshotPost(hostname[0][0].(string),req,cname,snap)
	
	responseType := operationOrError(res)
	if responseType == "operation" {
		endpointUrl,endpointResponse,_ = parseResponseRawToOperation(res)
		return OperationResponse(endpointUrl,&endpointResponse)
		}else{
		errorResp := parseErrorResponse(res)
		return &errorResp
		}
}

func snapshotsDelete(lx *LxdpmApi,  r *http.Request) Response {
	var endpointResponse api.Operation
	var endpointUrl string
	cname := mux.Vars(r)["name"]
	snap := mux.Vars(r)["snapshotName"]
	hostname := getHostnameFromContainername(lx,cname)

	res, _ := doSnapshotDelete(hostname[0][0].(string),cname,snap)
	responseType := operationOrError(res)
	if responseType == "operation" {
		endpointUrl,endpointResponse,_ = parseResponseRawToOperation(res)
		return OperationResponse(endpointUrl,&endpointResponse)
		}else{
		errorResp := parseErrorResponse(res)
		return &errorResp
		}
}