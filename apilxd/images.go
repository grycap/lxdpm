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
	//"io/ioutil"
)

var imagesCmd = Command{
	name: "images",
	get: imagesGetAllLXD,
	post: imagesPost,
}

type HostImageMetadata struct {
	Name 		string 	`json:"name"`
	Images 	[]string `json:"containers"`
}

func imagesGetAllLXD(lx *LxdpmApi,  r *http.Request) Response {
	var keys []string
	for k := range DefaultHosts {
		keys = append(keys,k)
	}
	var result []HostImageMetadata
	var wg sync.WaitGroup
	var resultLXD []string
	var metadata_hosts = make(chan HostImageMetadata,len(keys))
	defer close(metadata_hosts)
	
	wg.Add(len(keys))
	sort.Strings(keys)
	fmt.Println(keys)
	for _,k := range keys {

		go func (key string) {
			defer wg.Done()
			if DefaultHosts[key].Name == "local" {
					metadata_hosts <- imagesGetMetadataLocal()
			} else {
					metadata_hosts <- imagesGetMetadata(DefaultHosts[key].Name)
			}
			
		}(k)
	}
	//fmt.Printf("%+v %s",metadata_hosts,cap(metadata_hosts))
	
	for i :=0 ;i < len(keys); i++ {
		result = append(result,<-metadata_hosts)
	}
	wg.Wait()
	for _,v := range result {
		addImagesToHostDB(lx, v.Name ,v.Images)
		resultLXD = append(resultLXD,(v.Images)...)
	}
	return SyncResponse(true,resultLXD)
}

func imagesGetMetadata(hostname string) HostImageMetadata {

	argstr := []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket s/1.0/images"}  
    out, err := exec.Command("ssh", argstr...).Output()
    if err != nil {
        fmt.Println(err)
    }
    meta := parseImagesMetadataFromResponse(hostname,out)
    //fmt.Println(meta)
    //fmt.Println(result)
    return meta
}

func imagesGetMetadataLocal() HostImageMetadata {

	argstr := []string{"-s","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/images"}
    out, err := exec.Command("curl", argstr...).Output()
    if err != nil {
        fmt.Println(err)
    }
    meta := parseImagesMetadataFromResponse("local",out)
    //fmt.Println(meta)
    //fmt.Println(result)
    return meta
}

func parseImagesMetadataFromResponse(hostname string, input []byte) (res HostImageMetadata) {
	var resp = api.Response{}
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    res.Name = hostname
    json.NewDecoder(bytes.NewReader(resp.Metadata)).Decode(&res.Images)
    return res
}

func getImageIdDB(lx *LxdpmApi,name string) (string,error) {
	inargs := []interface{}{}
	outargs := []interface{}{"id"}
	//cash, err := lx.db.Query(`SELECT * FROM hosts`)
	result, err := dbQueryScan(lx.db, `SELECT id FROM images where fingerprint='`+name+`'`,inargs,outargs )
	if err != nil {
		fmt.Println(err)
		return "",err
	}
	if len(result) == 0 {
		return "" , nil
	}

	return result[0][0].(string) ,nil
}

func createImageDB(lx *LxdpmApi,hostid string,fingerprint string) error{
	q := `INSERT INTO images (fingerprint,host_id) VALUES (?,?)`
	_,err := dbExec(lx.db,q,fingerprint,hostid)
	return err
}

func addImagesToHostDB(lx *LxdpmApi,hostname string,images []string) {
	hostid := getHostId(lx,hostname)
	var fingerprint []string
	for _,image := range images {

		fingerprint = strings.Split(image,"/")
		id,err := getImageIdDB(lx,fingerprint[len(fingerprint)-1])

		if err != nil {
			fmt.Println(err)
		} 
		if id == "" {
			err := createImageDB(lx,hostid,fingerprint[len(fingerprint)-1])
			if err != nil {
				fmt.Println(err)
			}
		} 
	}
}

func getHostsDB(lx *LxdpmApi) ([][]interface {},error) {
	inargs := []interface{}{}
	outargs := []interface{}{"id","name","ip"}
	result, err := dbQueryScan(lx.db, `SELECT * FROM hosts;`,inargs,outargs )
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	if len(result) == 0 {
		return nil , nil
	}
	fmt.Println(result)
	return result ,nil
}

/*func imagesPostHost(lx *LxdpmApi,  r *http.Request) Response {
	var keys []string
	var wg sync.WaitGroup
	//var resultLXD []string
	
	req := api.ImagesPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	fmt.Printf("\nReq: %+v",req)
	hosts, err:= getHostsDB(lx)
	if err != nil {
		fmt.Println(err)
	}
	for _, host := range hosts {
		keys = append(keys,host[1].(string))
	}
	var result []LxdResponseRaw
	var metadata_hosts = make(chan LxdResponseRaw,len(keys))
	defer close(metadata_hosts)
	wg.Add(len(keys))
	sort.Strings(keys)
	fmt.Println(keys)
	for i,_ := range keys {
		key := keys[i] 
		go func (key string,freq api.ImagesPost) {
			defer wg.Done()
			metadata_hosts <- doImagePost(key,freq)
		}(key,req)
	}
	for i :=0 ;i < len(keys); i++ {
		result = append(result,<-metadata_hosts)
		//fmt.Printf("Soy result: %+v",result)
	}
	wg.Wait()
	for _,v := range result {
		fmt.Printf("\nSoy v: %+v",v)
		//resultLXD = append(resultLXD,v)
	}
	//res := imagesPost(req)
	return AsyncResponse(true,"")
}*/

type OperationImageInfo struct {
	Operation 	string
	Host 		string
}

func imagesPost(lx *LxdpmApi,  r *http.Request) Response {
	req := api.ImagesPost{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}
	run := func(op *task) error {
		result,err := imagesPostHost(lx,req)
		fmt.Println(result)
		fmt.Println(err)
		return nil
	}

	op, err := taskCreate(taskClassOp, nil, nil, run, nil, nil)
	//fmt.Printf("Soy la op: %+v\n",op)
	fmt.Println("Soy el error de creación: ",err)
	tk,err := taskGet(op.id)
	//fmt.Printf("Soy la task: %+v\n",tk)
	//tk.Run() //De normal se llama al hacer el render de operation en Response.go
	fmt.Println("Soy el error de obtención: ",err)
	return TaskResponse(tk)
}

func imagesPostHost(lx *LxdpmApi,  imageJson api.ImagesPost ) ([]OperationHost, error) {
	var keys []string
	var wg sync.WaitGroup
	//var resultLXD []string
	
	fmt.Printf("\nReq: %+v",imageJson)
	hosts, err:= getHostsDB(lx)
	if err != nil {
		fmt.Println(err)
	}
	for _, host := range hosts {
		keys = append(keys,host[1].(string))
	}
	var result []OperationImageInfo
	var operation_info = make(chan OperationImageInfo,len(keys))
	defer close(operation_info)
	wg.Add(len(keys))
	sort.Strings(keys)
	fmt.Println(keys)
	for i,_ := range keys {
		key := keys[i] 
		go func (key string,freq api.ImagesPost) {
			defer wg.Done()
			operation_info <- doImagePost(key,freq)
		}(key,imageJson)
	}
	for i :=0 ;i < len(keys); i++ {
		result = append(result,<-operation_info)
		//fmt.Printf("Soy result: %+v",result)
	}
	wg.Wait()
	/*for _,v := range result {
		fmt.Printf("\nSoy v: %+v",v)
		//resultLXD = append(resultLXD,v)
	}*/
	opResults := watchImageOperation(result)

	//res := imagesPost(req)
	//return AsyncResponse(true,"")
	return opResults,nil
}

func doImagePost(hostname string,r api.ImagesPost ) OperationImageInfo {
	opResult := OperationImageInfo{}
	headers := "-H 'X-LXD-filename: imageTest' -H 'X-LXD-public: true' -H 'X-LXD-properties: os=Ubuntu&alias=pruebaComun'"
	body ,err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
	/*body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}*/
	strbody := string(body)
	//fmt.Println(strbody)
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	//fmt.Println("\n"+string(buf))
	//fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -k --unix-socket /var/lib/lxd/unix.socket "+headers+" -X POST -d '"+strbody+"' s/1.0/images"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		headers := []string{"-H","'X-LXD-filename: imageTest'","-H","'X-LXD-public: true'","-H","'X-LXD-properties: os=Ubuntu&alias=pruebaComun'"}
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket"}
		argstr = append(argstr,headers...)
		argstr = append(argstr,[]string{"-X","POST","-d",fmt.Sprintf(""+strbody+""),"s/1.0/images"}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
    out, err := command.Output()
    fmt.Println("\nOut: ",string(out))
    if err != nil {
        fmt.Println(err)
    }
    resp := parseMetadataFromOperationResponseClean(out)
    //fmt.Println("\n",resp.Operation)
    opResult.Operation = resp.Operation
    opResult.Host = hostname
    return opResult
}

type OperationHost struct {
	Response 	api.Operation
	Host 		string
}

func watchImageOperation(opInfo []OperationImageInfo) []OperationHost {
	var wg sync.WaitGroup
	var result []OperationHost
	var operationResponse = make(chan OperationHost,len(opInfo))
	defer close(operationResponse)
	wg.Add(len(opInfo))
	fmt.Println(opInfo)
	for i,_ := range opInfo {
		operation := opInfo[i] 
		go func (op OperationImageInfo) {
			defer wg.Done()
			operationResponse <- doWatchImageOperation(op)
		}(operation)
	}
	for i :=0 ;i < len(opInfo); i++ {
		result = append(result,<-operationResponse)
		//fmt.Printf("Soy result: %+v",result)
	}
	wg.Wait()
	//After knowing the result,we can decide if the operation is valid, or we have to do a rollback
	deleteImages := false
	for _,v := range result {
		//fmt.Printf("\nSoy v: %+v",v)
		if v.Response.Status == "Failure" {
			deleteImages = true
			break
		}
		//resultLXD = append(resultLXD,v)
	}
	if deleteImages == true {
		fmt.Println("hay al menos un failure en:")
		for _,v := range result {
			if v.Response.Status != "Failure" {
				fmt.Println("estoy borrando: "+v.Response.Metadata["fingerprint"].(string)+" en "+v.Host)
				deleteImage(v)
			}
			fmt.Println("No borro porque "+v.Response.Metadata["fingerprint"].(string)+" no se ha llegado a crear")
		}
		return nil
	}
	fmt.Println("Al parecer todo se ha creado bien.")
	return result
}

func doWatchImageOperation(op OperationImageInfo) OperationHost {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	hostname := op.Host
	//fmt.Println("\n"+string(buf))
	//fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -k --unix-socket /var/lib/lxd/unix.socket s"+op.Operation+"/wait"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","s"+op.Operation+"/wait"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
    out, err := command.Output()
    fmt.Println("\nOut: ",string(out))
    if err != nil {
        fmt.Println(err)
    }
    resp := parseMetadataFromOperationResponseClean(out)
    fmt.Printf("After wait: %+v",resp)
    opresp := parseOperationFromMetadata(resp.Metadata)
    fmt.Printf("After get operation: %+v",resp)
    hostResp := OperationHost{Response: opresp,Host: hostname}
    return hostResp
}

func parseOperationFromMetadata(input []byte) api.Operation {
	var resp = api.Operation{}
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    //fmt.Println(resp.Metadata)
    //fmt.Println(string(resp.Metadata))
    return resp
}

func deleteImage(oph OperationHost) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	hostname := oph.Host
	//fmt.Println("\n"+string(buf))
	//fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -k --unix-socket /var/lib/lxd/unix.socket s/1.0/images/"+oph.Response.Metadata["fingerprint"].(string)}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/images/"+oph.Response.ID}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
    out, err := command.Output()
    fmt.Println("\nOut: ",string(out))
    if err != nil {
        fmt.Println(err)
    }
}
