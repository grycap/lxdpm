package apilxd

import (
	"fmt"
	"net/http"
	"sort"
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
	post: containerPostPlanner,
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
			out,_ := doContainersGet(DefaultHosts[key].Name)
			metadata_hosts <- parseMetadataFromMultipleContainersResponse(key,out)
		}(k)
	}
	//fmt.Printf("%+v %s",metadata_hosts,cap(metadata_hosts))
	
	for i :=0 ;i < len(keys); i++ {
		result = append(result,<-metadata_hosts)
	}
	wg.Wait()
	for _,v := range result {
		addContainersToHostDB(lx, v.Name ,v.Containers)
		resultLXD = append(resultLXD,(v.Containers)...)
	}

	return SyncResponse(true,resultLXD)
}

type ContainersHostPost struct {

	Hostname   string          `json:"hostname" yaml:"hostname"`
	ContainersPost api.ContainersPost `json:"containersPost" yaml:"containersPost"`
}

func parseMetadataFromOperationResponse(input []byte) LxdResponseRaw {
	var resp = LxdResponseRaw{}
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    //fmt.Println(resp.Metadata)
    //fmt.Println(string(resp.Metadata))
    return resp
}

func parseMetadataFromOperationResponseClean(input []byte) api.Response {
	var resp = api.Response{}
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    //fmt.Println(resp.Metadata)
    //fmt.Println(string(resp.Metadata))
    return resp
}
/*
func containerPostHost(lx *LxdpmApi,  r *http.Request) Response {
	var endpointResponse api.Operation
	var endpointUrl string
	req := ContainersHostPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)
	res,_ := doContainersPost(req)
	responseType := operationOrError(res)
	if responseType == "operation" {
		endpointUrl,endpointResponse,_ = parseResponseRawToOperation(res)
		return OperationResponse(endpointUrl,&endpointResponse)
		}else{
		errorResp := parseErrorResponse(res)
		return &errorResp
		}
}
*/
func containerPostPlanner(lx *LxdpmApi,  r *http.Request) Response {
	var endpointResponse api.Operation
	var endpointUrl string
	req := api.ContainersPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)
	hostname,err := lx.planner.HostToDeploy()
	if err != nil {
		return BadRequest(err)
	}
	res,_ := doContainersPlannerPost(req,hostname)
	responseType := operationOrError(res)
	if responseType == "operation" {
		endpointUrl,endpointResponse,_ = parseResponseRawToOperation(res)
		return OperationResponse(endpointUrl,&endpointResponse)
		}else{
		errorResp := parseErrorResponse(res)
		return &errorResp
		}
}

func getHostId(lx *LxdpmApi,name string) string {
	inargs := []interface{}{}
	outargs := []interface{}{"id"}

	result, err := dbQueryScan(lx.db, `SELECT id FROM hosts where name='`+name+`'`,inargs,outargs )
	if err != nil {
		fmt.Println(err)
	}
	if len(result) == 0 {
		return ""
	}

	return result[0][0].(string)
}

func getContainerIdDB(lx *LxdpmApi,name string) (string,error) {
	inargs := []interface{}{}
	outargs := []interface{}{"id"}
	//cash, err := lx.db.Query(`SELECT * FROM hosts`)
	result, err := dbQueryScan(lx.db, `SELECT id FROM containers where name='`+name+`'`,inargs,outargs )
	if err != nil {
		fmt.Println(err)
		return "",err
	}
	if len(result) == 0 {
		return "" , nil
	}

	return result[0][0].(string) ,nil
}

func createContainerDB(lx *LxdpmApi,hostid string,cname string) error{
	q := `INSERT INTO containers (name,host_id) VALUES (?,?)`
	_,err := dbExec(lx.db,q,cname,hostid)
	return err
}

func deleteContainerDB(lx *LxdpmApi,cname string) error{
	q := `DELETE FROM containers WHERE name=?`
	_,err := dbExec(lx.db,q,cname)
	return err
}

func updateContainerDB(lx *LxdpmApi,cname string,newname string) error{
	q := `UPDATE containers SET name=? WHERE name=?;`
	_,err := dbExec(lx.db,q,newname,cname)
	return err
}

func addContainersToHostDB(lx *LxdpmApi,hostname string,containers []string) {
	hostid := getHostId(lx,hostname)
	var cname []string
	for _,container := range containers {

		cname = strings.Split(container,"/")
		id,err := getContainerIdDB(lx,cname[len(cname)-1])

		if err != nil {
			fmt.Println(err)
		} 
		if id == "" {
			err := createContainerDB(lx,hostid,cname[len(cname)-1])
			if err != nil {
				fmt.Println(err)
			}
		} 
	}
}

var containerCmd = Command{
	name:   "containers/{name}",
	get:    containerGet,
	put:    containerPut,
	delete: containerDelete,
	post:   containerPost,
	/*patch:  containerPatch,
	*/
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
	get:    snapshotsGet,
	post:   snapshotsPost,
	delete: snapshotsDelete,
}


var containerExecCmd = Command{
	name: "containers/{name}/exec",
	post: containerExecPost,
}