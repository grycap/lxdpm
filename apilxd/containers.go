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
	get:  containersGetAll,
	//post: containersPost,
}

func containersGetAll(lx *LxdpmApi,  r *http.Request) Response {
	var result []HostContainerMetadata
	var wg sync.WaitGroup
	var metadata_hosts = make(chan HostContainerMetadata)
	defer close(metadata_hosts)
	var keys []string
	for k := range DefaultHosts {
		keys = append(keys,k)
	}
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
	
	fmt.Println(metadata_hosts)
	go func() {
        for response := range metadata_hosts {
        	//fmt.Println(response)
            result = append(result,response)
            //fmt.Println(result)
        }
    }()
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