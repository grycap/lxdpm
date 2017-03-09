package apilxd

import (
	"fmt"
	"net/http"
	"os/exec"
	"encoding/json"
	"bytes"
	"github.com/lxc/lxd/shared/api"
	"strings"
	//"os"
	//"log"
	/*"io"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/lxc/go-lxc.v2"

	"github.com/lxc/lxd/lxd/types"
	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/osarch"*/
)


var containersCmd = Command{
	name: "containers",
	get:  containersGet,
	//post: containersPost,
}

/*Returns metadata array
func containersGetAll(lx *LxdpmApi,  r *http.Request) Response {
	var result []string = []string{}  
	argstr := []string{"troig@lxdpm02", "curl -s --unix-socket /var/lib/lxd/unix.socket s/1.0/containers"}
    out, err := exec.Command("ssh", argstr...).Output()
    if err != nil {
        fmt.Println(err)
    }
    metadata := parseMetadataFromResponse(out)
    result = append(result,metadata)
    argstr = []string{"-s","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/containers"}
    out, err = exec.Command("curl", argstr...).Output()
    if err != nil {
        fmt.Println(err)
    }
    result = append(result,parseMetadataFromResponse(out))
    //fmt.Println(string(result))
	return SyncResponse(true,interface{}(result))
}*/
/*
func containersGetAllAppendString(lx *LxdpmApi,  r *http.Request) Response {
	var result string = "" 
	argstr := []string{"-s","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/containers"}
    out, err := exec.Command("curl", argstr...).Output()
    if err != nil {
        fmt.Println(err)
    }
    result = result + string(out)

	return SyncResponse(true,interface{}(string(result)))
}*/

func containersGet(lx *LxdpmApi,  r *http.Request) Response {
	var result []string = []string{}  
	argstr := []string{"-s","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/containers"}
    out, err := exec.Command("curl", argstr...).Output()
    if err != nil {
        fmt.Println(err)
    }
    metadata := parseMetadataFromResponse(out)
    fmt.Println(strings.Split(metadata,""))
    splitted := strings.Split(metadata,"")
    fmt.Println(len(splitted))
    splitted = splitted[2:len(splitted)-2]
    result = append(result,strings.Join(splitted,""))
    fmt.Println(result)
	return SyncResponse(true,interface{}(result))
}

func parseMetadataFromResponse(input []byte) (res string) {
	var resp = api.Response{}
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    res = string(resp.Metadata)
    return res
} 


/*
	var result1 = api.Response{}
    json.NewDecoder(bytes.NewReader(out)).Decode(&result1)
    fmt.Printf("%+v",result1)
    meta := string(result1.Metadata)
    fmt.Printf("%v",meta)

func containersGet(lx *LxdpmApi,  r *http.Request) Response {
	argstr := []string{"troig@lxdpm02", "curl -s --unix-socket /var/lib/lxd/unix.socket s/1.0/containers"}
    out, err := exec.Command("ssh", argstr...).Output()
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    fmt.Println(string(out))
	fmt.Println("Tutto beneeeeeeeeeeeeee")
	return SyncResponse(true,interface{}(string(out)))
}*/
/*func containersGet(lx *LxdpmApi,  r *http.Request) Response {
	var containers,err = lx.Cli.ListContainers()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(containers)
	}
	fmt.Println("Tutto beneeeeeeeeeeeeee")
	return SyncResponse(true,interface{}(containers))
}
*/

/*
var containerCmd = Command{
	name:   "containers/{name}",
	get:    containerGet,
	put:    containerPut,
	delete: containerDelete,
	post:   containerPost,
	patch:  containerPatch,
}

var containerStateCmd = Command{
	name: "containers/{name}/state",
	get:  containerState,
	put:  containerStatePut,
}

var containerFileCmd = Command{
	name:   "containers/{name}/files",
	get:    containerFileHandler,
	post:   containerFileHandler,
	delete: containerFileHandler,
}

var containerSnapshotsCmd = Command{
	name: "containers/{name}/snapshots",
	get:  containerSnapshotsGet,
	post: containerSnapshotsPost,
}

var containerSnapshotCmd = Command{
	name:   "containers/{name}/snapshots/{snapshotName}",
	get:    snapshotHandler,
	post:   snapshotHandler,
	delete: snapshotHandler,
}

var containerExecCmd = Command{
	name: "containers/{name}/exec",
	post: containerExecPost,
}*/