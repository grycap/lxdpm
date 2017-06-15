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
/*
func doContainerExecPost(lx *LxdpmApi,req api.ContainerExecPost,cname string) Response {
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	hostname := getHostnameFromContainername(lx,cname)
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname[0][0].(string) != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname[0][0].(string)},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname+"/exec"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname+"/exec"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
    out, err := command.Output()
    fmt.Println("\nOut: ",string(out))
    if err != nil {
        fmt.Println(err)
    }
    meta := parseMetadataFromOperationResponse(out)
    return AsyncResponse(true,meta)
}*/