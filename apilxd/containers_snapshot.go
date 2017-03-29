package apilxd

import (
	"net/http"
	"os/exec"
	"encoding/json"
	"bytes"
	"strings"
	"fmt"
	"github.com/lxc/lxd/shared/api"
	"github.com/gorilla/mux"
)

func containerSnapshotsGet(lx *LxdpmApi, r *http.Request) Response {
	name := mux.Vars(r)["name"]

	hostname := getHostnameFromContainername(lx,name)

	resp := doSnapshotsGet(name, hostname[0][0].(string))
	
	return resp
}

func doSnapshotsGet(cname string, hostname string) Response {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket s/1.0/containers/"+cname+"/snapshots"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/containers/"+cname+"/snapshots"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
    out, err := command.Output()
    fmt.Println(string(out))
    fmt.Println(out)
    if err != nil {
        fmt.Println(err)
    }
    resp := parseSyncResponse(out)
    return resp
	
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
	name := mux.Vars(r)["name"]

	req := api.ContainerSnapshotsPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	resp := doContainerSnapshotPost(lx,req,name)
	
	return resp
}

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
}