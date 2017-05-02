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

/*func containerGet(lx *LxdpmApi, r *http.Request) Response {
	/*_, err := lx.db.Exec(DB_FILL)
	if err != nil {
		return BadRequest(err)
	}
	result := getHostnameFromContainername(lx,"otrohost2")
	return SyncResponse(true,result[0][0])
}*/

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
	name := mux.Vars(r)["name"]

	hostname := getHostnameFromContainername(lx,name)

	resp := containerGetMetadata(hostname[0][0].(string),name)

	//meta := resp.Metadata 

	return SyncResponse(true,resp)
}

func containerGetMetadata(hostname string, cname string) LxdResponseRaw {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket s/1.0/containers/"+cname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/containers/"+cname}
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

func parseMetadataFromContainerResponse(input []byte) LxdResponseRaw {
	var resp = LxdResponseRaw{}
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    return resp
}


func containerPut(lx *LxdpmApi,  r *http.Request) Response {
	name := mux.Vars(r)["name"]

	req := api.ContainerPut{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)
	res := doContainerPut(lx,req,name)
	return AsyncResponse(true,res)
}


func doContainerPut(lx *LxdpmApi,req api.ContainerPut,cname string) LxdResponseRaw {
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
		argstr = []string{strings.Join([]string{"troig","@",hostname[0][0].(string)},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X PUT -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","PUT","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname}
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

func containerDelete(lx *LxdpmApi,  r *http.Request) Response {
	name := mux.Vars(r)["name"]
	res := doContainerDelete(lx,name)
	return AsyncResponse(true,res)
}

func doContainerDelete(lx *LxdpmApi,cname string) LxdResponseRaw {
	hostname := getHostnameFromContainername(lx,cname)
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname[0][0].(string) != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname[0][0].(string)},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X DELETE s/1.0/containers/"+cname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","DELETE","s/1.0/containers/"+cname}
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

func containerPost(lx *LxdpmApi,  r *http.Request) Response {
	name := mux.Vars(r)["name"]

	req := api.ContainerPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)
	res := doContainerPost(lx,req,name)
	return AsyncResponse(true,res)
}


func doContainerPost(lx *LxdpmApi,req api.ContainerPost,cname string) LxdResponseRaw {
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
		argstr = []string{strings.Join([]string{"troig","@",hostname[0][0].(string)},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname}
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
