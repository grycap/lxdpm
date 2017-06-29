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
)

type HostProfileMetadata struct {
	Name 		string 	`json:"name"`
	Profiles 	[]string `json:"profiles"`
}

var profilesCmd = Command{
	name: "profiles",
	get:  profilesGet,
	post: profilesPostHost,
}


func profilesGet(lx *LxdpmApi,  r *http.Request) Response {
	var keys []string
	for k := range DefaultHosts {
		keys = append(keys,k)
	}
	var result []HostProfileMetadata
	var wg sync.WaitGroup
	var resultLXD []string
	var metadata_hosts = make(chan HostProfileMetadata,len(keys))
	defer close(metadata_hosts)
	
	wg.Add(len(keys))
	sort.Strings(keys)
	fmt.Println(keys)
	for _,k := range keys {

		go func (key string) {
			defer wg.Done()
			out,_ := doProfilesGet(DefaultHosts[key].Name)
			metadata_hosts <- parseProfileMetadataFromResponse(key,out)
		}(k)
	}
	//fmt.Printf("%+v %s",metadata_hosts,cap(metadata_hosts))
	
	for i :=0 ;i < len(keys); i++ {
		result = append(result,<-metadata_hosts)
	}
	wg.Wait()
	for _,v := range result {
		addProfilesToHostDB(lx, v.Name ,v.Profiles)
		resultLXD = append(resultLXD,(v.Profiles)...)
	}
	return SyncResponse(true,resultLXD)
}

func parseProfileMetadataFromResponse(hostname string, input []byte) (res HostProfileMetadata) {
	var resp = api.Response{}
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    res.Name = hostname
    json.NewDecoder(bytes.NewReader(resp.Metadata)).Decode(&res.Profiles)
    return res
}

func getProfileIdDB(lx *LxdpmApi,name string) (string,error) {
	inargs := []interface{}{}
	outargs := []interface{}{"id"}
	//cash, err := lx.db.Query(`SELECT * FROM hosts`)
	result, err := dbQueryScan(lx.db, `SELECT id FROM profiles where name='`+name+`'`,inargs,outargs )
	if err != nil {
		fmt.Println(err)
		return "",err
	}
	if len(result) == 0 {
		return "" , nil
	}

	return result[0][0].(string) ,nil
}

func createProfileDB(lx *LxdpmApi,hostid string,pname string) error{
	q := `INSERT INTO profiles (name,host_id) VALUES (?,?)`
	_,err := dbExec(lx.db,q,pname,hostid)
	return err
}

func addProfilesToHostDB(lx *LxdpmApi,hostname string,profiles []string) {
	hostid := getHostId(lx,hostname)
	var pname []string
	for _,profile := range profiles {

		pname = strings.Split(profile,"/")
		id,err := getProfileIdDB(lx,pname[len(pname)-1])

		if err != nil {
			fmt.Println(err)
		} 
		if id == "" {
			err := createProfileDB(lx,hostid,pname[len(pname)-1])
			if err != nil {
				fmt.Println(err)
			}
		} 
	}
}

type ProfilesHostPost struct {

	Hostname   string          `json:"hostname" yaml:"hostname"`
	ProfilesPost api.ProfilesPost `json:"profilesPost" yaml:"profilesPost"`
}

func profilesPostHost(lx *LxdpmApi,  r *http.Request) Response {
	req := ProfilesHostPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)
	res,_ := doProfilesPost(req)
	responseType := getResponseType(res)
	if responseType == "sync" {
		endpointResponse,_ := parseResponseRawToSync(res)
		return &endpointResponse
	}else {
		errorResp := parseErrorResponse(res)
		return &errorResp
	}
}
