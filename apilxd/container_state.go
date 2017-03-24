package apilxd

import (
	"net/http"
	"os/exec"
	"strings"
	"fmt"
	"encoding/json"
	"github.com/lxc/lxd/shared/api"
	"github.com/gorilla/mux"
)

func containerState(lx *LxdpmApi, r *http.Request) Response {
	name := mux.Vars(r)["name"]

	hostname := getHostnameFromContainername(lx,name)

	resp := containerStateGetMetadata(hostname[0][0].(string),name)

	//meta := resp.Metadata 

	return SyncResponse(true,resp)
}

func containerStateGetMetadata(hostname string, cname string) LxdResponseRaw {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket s/1.0/containers/"+cname+"/state"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/containers/"+cname+"/state"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
    out, err := command.Output()
    if err != nil {
        fmt.Println(err)
    }
    meta := parseMetadataFromContainerResponse(out)
    return meta
}

func containerStatePut(lx *LxdpmApi,  r *http.Request) Response {
	name := mux.Vars(r)["name"]

	req := api.ContainerStatePut{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)
	res := doContainerStatePut(lx,req,name)
	return AsyncResponse(true,res)
}

func doContainerStatePut(lx *LxdpmApi,req api.ContainerStatePut,cname string) LxdResponseRaw {
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
		argstr = []string{strings.Join([]string{"troig","@",hostname[0][0].(string)},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X PUT -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname+"/state"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","PUT","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname+"/state"}
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