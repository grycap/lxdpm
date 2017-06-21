package apilxd

import (
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
	"github.com/lxc/lxd/shared/api"
	"github.com/gorilla/mux"
)

func containerSnapshotsGet(lx *LxdpmApi, r *http.Request) Response {
	name := mux.Vars(r)["name"]

	hostname := getHostnameFromContainername(lx,name)

	res,_ := doContainerSnapshotsGet(hostname[0][0].(string),name)
	
	endpointResponse,_ := parseResponseRawToSync(res)

	return &endpointResponse 
}

func parseSyncResponse(input []byte) Response {
	req := api.ResponseRaw{}
	fmt.Println(input)
	fmt.Println(string(input))
	if err := json.NewDecoder(bytes.NewReader(input)).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)


	return SyncResponse(true,req.Metadata)

}

func containerSnapshotsPost(lx *LxdpmApi, r *http.Request) Response {
	var endpointResponse api.Operation
	var endpointUrl string
	cname := mux.Vars(r)["name"]
	hostname := getHostnameFromContainername(lx,cname)

	req := api.ContainerSnapshotsPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	res,_ := doContainerSnapshotPost(hostname[0][0].(string),req,cname)
	
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
func doContainerSnapshotPost(lx *LxdpmApi,req api.ContainerSnapshotsPost,cname string) Response {
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
		argstr = []string{strings.Join([]string{"troig","@",hostname[0][0].(string)},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname+"/snapshots"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname+"/snapshots"}
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