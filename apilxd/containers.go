package apilxd

import (
	"fmt"
	"net/http"
	"sort"
	"os/exec"
	"encoding/json"
	"bytes"
	"github.com/lxc/lxd/shared/api"
	"strings"
	"sync"
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

type HostContainerMetadata struct {
	Name 		string 	`json:"name"`
	Containers 	[]string `json:"containers"`
}

var containersCmd = Command{
	name: "containers",
	get:  containersGetAllLXD,
	post: containerPostHost,
}

func containersGetAllLXD(lx *LxdpmApi,  r *http.Request) Response {
	var keys []string
	for k := range DefaultHosts {
		keys = append(keys,k)
	}
	var result []HostContainerMetadata
	var wg sync.WaitGroup
	var resultLXD []string
	var metadata_hosts = make(chan HostContainerMetadata,len(keys))
	defer close(metadata_hosts)
	
	wg.Add(len(keys))
	sort.Strings(keys)
	fmt.Println(keys)
	for _,k := range keys {

		go func (key string) {
			defer wg.Done()
			if DefaultHosts[key].Name == "local" {
					metadata_hosts <- containersGetMetadataLocal()
			} else {
					metadata_hosts <- containersGetMetadata(DefaultHosts[key].Name)
			}
			
		}(k)
	}
	//fmt.Printf("%+v %s",metadata_hosts,cap(metadata_hosts))
	
	for i :=0 ;i < len(keys); i++ {
		result = append(result,<-metadata_hosts)
	}
	wg.Wait()
	for _,v := range result {
		resultLXD = append(resultLXD,(v.Containers)...)
	}
	return SyncResponse(true,resultLXD)
}

func containersGetAll(lx *LxdpmApi,  r *http.Request) Response {
	var keys []string
	for k := range DefaultHosts {
		keys = append(keys,k)
	}
	var result []HostContainerMetadata
	var wg sync.WaitGroup
	var metadata_hosts = make(chan HostContainerMetadata,len(keys))
	defer close(metadata_hosts)
	
	wg.Add(len(keys))
	sort.Strings(keys)
	fmt.Println(keys)
	for _,k := range keys {

		go func (key string) {
			defer wg.Done()
			if DefaultHosts[key].Name == "local" {
					metadata_hosts <- containersGetMetadataLocal()
			} else {
					metadata_hosts <- containersGetMetadata(DefaultHosts[key].Name)
			}
			
		}(k)
	}
	//fmt.Printf("%+v %s",metadata_hosts,cap(metadata_hosts))
	
	for i :=0 ;i < len(keys); i++ {
		result = append(result,<-metadata_hosts)
	}
	wg.Wait()
	return SyncResponse(true,result)
}
func containersGetMetadata(hostname string) HostContainerMetadata {

	argstr := []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket s/1.0/containers"}  
    out, err := exec.Command("ssh", argstr...).Output()
    if err != nil {
        fmt.Println(err)
    }
    meta := parseMetadataFromResponse(hostname,out)
    //fmt.Println(meta)
    //fmt.Println(result)
    return meta
}

func containersGetMetadataLocal() HostContainerMetadata {

	argstr := []string{"-s","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/containers"}
    out, err := exec.Command("curl", argstr...).Output()
    if err != nil {
        fmt.Println(err)
    }
    meta := parseMetadataFromResponse("local",out)
    //fmt.Println(meta)
    //fmt.Println(result)
    return meta
}
func parseMetadataFromResponse(hostname string, input []byte) (res HostContainerMetadata) {
	var resp = api.Response{}
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    res.Name = hostname
    json.NewDecoder(bytes.NewReader(resp.Metadata)).Decode(&res.Containers)
    return res
}

type ContainersHostPost struct {

	Hostname   string          `json:"hostname" yaml:"hostname"`
	ContainersPost api.ContainersPost `json:"containersPost" yaml:"containerPost"`
}

func parseMetadataFromOperationResponse(input []byte) LxdResponseRaw {
	var resp = LxdResponseRaw{}
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    //fmt.Println(resp.Metadata)
    //fmt.Println(string(resp.Metadata))
    return resp
}
func containerPostHost(lx *LxdpmApi,  r *http.Request) Response {
	req := ContainersHostPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)
	res := containerPost(req)
	return AsyncResponse(true,res)
}


func containerPost(req ContainersHostPost) LxdResponseRaw {
	buf ,err := json.Marshal(req.ContainersPost)
	if err != nil {
		fmt.Println(err)
	}
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if req.Hostname != "" {
		argstr = []string{strings.Join([]string{"troig","@",req.Hostname},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers"}
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

var containerCmd = Command{
	name:   "containers/{name}",
	get:    containerGet,
	/*put:    containerPut,
	delete: containerDelete,
	post:   containerPost,
	patch:  containerPatch,
	*/
}
/*
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