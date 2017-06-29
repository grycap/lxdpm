package apilxd

import (
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
	"github.com/lxc/lxd/shared/api"
	"github.com/gorilla/mux"
)

func getDBhosts(lx *LxdpmApi) [][]interface{} {
	inargs := []interface{}{}
	outargs := []interface{}{"id","localhost","ip"}
	//cash, err := lx.db.Query(`SELECT * FROM hosts`)
	result, err := dbQueryScan(lx.db, `SELECT * FROM hosts`,inargs,outargs )
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func getHostnameFromContainername(lx *LxdpmApi, name string) [][]interface{} {
	inargs := []interface{}{}
	outargs := []interface{}{"name"}

	result, err := dbQueryScan(lx.db, `SELECT H.name FROM hosts H, containers C where C.host_id = H.id and C.name='`+name+`';`,inargs,outargs )
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func containerGet(lx *LxdpmApi, r *http.Request) Response {
	cname := mux.Vars(r)["name"]

	hostname := getHostnameFromContainername(lx,cname)

	resp,_ := doContainerGet(hostname[0][0].(string),cname)

	endpointResponse,_ := parseResponseRawToSyncContainerGet(resp,hostname[0][0].(string))

	return &endpointResponse 
}

func parseMetadataFromContainerResponse(input []byte) LxdResponseRaw {
	var resp = LxdResponseRaw{}
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    return resp
}


func containerPut(lx *LxdpmApi,  r *http.Request) Response {
	cname := mux.Vars(r)["name"]
	hostname := getHostnameFromContainername(lx,cname)
	req := api.ContainerPut{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}

	res,_ := doContainerPut(hostname[0][0].(string),req,cname)

	endpointUrl,endpointResponse,_ := parseResponseRawToOperation(res) 

	return OperationResponse(endpointUrl,&endpointResponse)
}

func containerDelete(lx *LxdpmApi,  r *http.Request) Response {
	var endpointResponse api.Operation
	var endpointUrl string

	cname := mux.Vars(r)["name"]

	hostname := getHostnameFromContainername(lx,cname)

	res,_ := doContainerDelete(hostname[0][0].(string),cname)

	responseType := operationOrError(res)
	if responseType == "operation" {
		err := deleteContainerDB(lx,cname)
		if err != nil {
				fmt.Println(err)
			}
		endpointUrl,endpointResponse,_ = parseResponseRawToOperation(res)
		return OperationResponse(endpointUrl,&endpointResponse)
		}else{
		errorResp := parseErrorResponse(res)
		return &errorResp
		}
}
func containerPost(lx *LxdpmApi,  r *http.Request) Response {
	var endpointResponse api.Operation
	var endpointUrl string
	cname := mux.Vars(r)["name"]
	hostname := getHostnameFromContainername(lx,cname)
	req := api.ContainerPost{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)
	res,newname,_ := doContainerPost(hostname[0][0].(string),req,cname)

	responseType := operationOrError(res)
	if responseType == "operation" {
		err := updateContainerDB(lx,cname,newname)
		if err != nil {
				fmt.Println(err)
			}
		endpointUrl,endpointResponse,_ = parseResponseRawToOperation(res)
		return OperationResponse(endpointUrl,&endpointResponse)
		}else{
		errorResp := parseErrorResponse(res)
		return &errorResp
		}
}
