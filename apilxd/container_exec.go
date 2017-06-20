package apilxd

import (
	"net/http"
	"encoding/json"
	"github.com/lxc/lxd/shared/api"
	"github.com/gorilla/mux"
	"fmt"
	"errors"
)

func containerExecPost(lx *LxdpmApi, r *http.Request) Response {
	var endpointResponse api.Operation
	var endpointUrl string
	cname := mux.Vars(r)["name"]
	hostname := getHostnameFromContainername(lx,cname)

	if len(hostname) == 0 {
		return BadRequest(errors.New("Container"+cname+" not in database. Do get all containers and try again."))
    }

	req := api.ContainerExecPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	res, err := doContainerExecPost(hostname[0][0].(string),req,cname)
	if err != nil {
        fmt.Println("Este es el error: ",err)
    }

	responseType := operationOrError(res)
	if responseType == "operation" {
		endpointUrl,endpointResponse,_ = parseResponseRawToOperation(res)
		return OperationResponse(endpointUrl,&endpointResponse)
		}else{
		errorResp := parseErrorResponse(res)
		return &errorResp
		}
}
