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


var profileCmd = Command{
	name: "profiles/{name}",
	get: profileGet,
	put: profilePut,
	delete: profileDelete,
	post: profilePost,
//	patch: profilePatch
}

func getHostnameFromProfileName(lx *LxdpmApi, name string) [][]interface{} {
	inargs := []interface{}{}
	outargs := []interface{}{"name"}

	result, err := dbQueryScan(lx.db, `SELECT H.name FROM hosts H, profiles P where P.host_id = H.id and P.name='`+name+`';`,inargs,outargs )
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func profileGet(lx *LxdpmApi, r *http.Request) Response {
	name := mux.Vars(r)["name"]

	hostname := getHostnameFromProfileName(lx,name)

	resp := profileGetMetadata(hostname[0][0].(string),name)

	//meta := resp.Metadata 

	return SyncResponse(true,resp)
}

func profileGetMetadata(hostname string, pname string) LxdResponseRaw {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket s/1.0/profiles/"+pname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/profiles/"+pname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
    out, err := command.Output()
    if err != nil {
        fmt.Println(err)
    }
    meta := parseMetadataFromProfileResponse(out)
    return meta
}

func parseMetadataFromProfileResponse(input []byte) LxdResponseRaw {
	var resp = LxdResponseRaw{}
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    return resp
}

func profilePut(lx *LxdpmApi,  r *http.Request) Response {
	name := mux.Vars(r)["name"]

	req := api.ProfilePut{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)
	res := doProfilePut(lx,req,name)
	return AsyncResponse(true,res)
}


func doProfilePut(lx *LxdpmApi,req api.ProfilePut,pname string) LxdResponseRaw {
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	hostname := getHostnameFromProfileName(lx,pname)
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname[0][0].(string) != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname[0][0].(string)},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X PUT -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/profiles/"+pname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","PUT","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/profiles/"+pname}
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

func profileDelete(lx *LxdpmApi,  r *http.Request) Response {
	name := mux.Vars(r)["name"]
	res := doProfileDelete(lx,name)
	return AsyncResponse(true,res)
}

func doProfileDelete(lx *LxdpmApi,pname string) LxdResponseRaw {
	hostname := getHostnameFromProfileName(lx,pname)
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname[0][0].(string) != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname[0][0].(string)},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X DELETE s/1.0/profiles/"+pname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","DELETE","s/1.0/profiles/"+pname}
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

func profilePost(lx *LxdpmApi,  r *http.Request) Response {
	name := mux.Vars(r)["name"]

	req := api.ProfilePost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)
	res := doProfilePost(lx,req,name)
	return AsyncResponse(true,res)
}


func doProfilePost(lx *LxdpmApi,req api.ProfilePost,pname string) LxdResponseRaw {
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	hostname := getHostnameFromProfileName(lx,pname)
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname[0][0].(string) != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname[0][0].(string)},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/profiles/"+pname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/profiles/"+pname}
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