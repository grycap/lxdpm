package apilxd

import (
	"net/http"
	"os/exec"
	"encoding/json"
	//"bytes"
	"strings"
	"fmt"
	"github.com/lxc/lxd/shared"
	"github.com/gorilla/mux"
)

func snapshotsGet(lx *LxdpmApi, r *http.Request) Response {
	name := mux.Vars(r)["name"]
	snap := mux.Vars(r)["snapshotName"]

	fmt.Println("Llego al endpoint")

	hostname := getHostnameFromContainername(lx,name)

	resp := doSnapshotGet(name, hostname[0][0].(string),snap)
	
	return resp
}

func doSnapshotGet(cname string, hostname string, snapname string) Response {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket s/1.0/containers/"+cname+"/snapshots/"+snapname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/containers/"+cname+"/snapshots/"+snapname}
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

func snapshotsPost(lx *LxdpmApi, r *http.Request) Response {
	name := mux.Vars(r)["name"]
	snap := mux.Vars(r)["snapshotName"]

	req := shared.Jmap{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	resp := doSnapshotPost(lx,req,name,snap)
	
	return resp
}

func doSnapshotPost(lx *LxdpmApi,req shared.Jmap,cname string,snap string) Response {
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
		argstr = []string{strings.Join([]string{"troig","@",hostname[0][0].(string)},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname+"/snapshots/"+snap}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname+"/snapshots/"+snap}
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

func snapshotsDelete(lx *LxdpmApi,  r *http.Request) Response {
	name := mux.Vars(r)["name"]
	snap := mux.Vars(r)["snapshotName"]
	res := doSnapshotDelete(lx,name,snap)
	return AsyncResponse(true,res)
}

func doSnapshotDelete(lx *LxdpmApi,cname string,snapname string) LxdResponseRaw {
	hostname := getHostnameFromContainername(lx,cname)
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname[0][0].(string) != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname[0][0].(string)},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X DELETE s/1.0/containers/"+cname+"/snapshots/"+snapname }
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","DELETE","s/1.0/containers/"+cname+"/snapshots/"+snapname }
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
    out, err := command.Output()
    fmt.Println("\nOut: ",string(out))
    if err != nil {
        fmt.Println(err)
    }
    meta := parseMetadataFromOperationResponse(out)
    return meta
}