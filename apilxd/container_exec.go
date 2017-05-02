package apilxd

import (
	"net/http"
	"os/exec"
	"encoding/json"
	//"bytes"
	"strings"
	"fmt"
	"github.com/lxc/lxd/shared/api"
	"github.com/gorilla/mux"
)

func containerExecPost(lx *LxdpmApi, r *http.Request) Response {
	name := mux.Vars(r)["name"]

	req := api.ContainerExecPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	resp := doContainerExecPost(lx,req,name)
	
	return resp
}

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
}