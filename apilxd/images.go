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
	"github.com/gorilla/mux"
	//"io/ioutil"
	"errors"
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
			out,_ := doImagesGet(DefaultHosts[key].Name)
			metadata_hosts <- parseMetadataFromAllImagesResponse(key,out)
		}(k)
	}
	//fmt.Printf("%+v %s",metadata_hosts,cap(metadata_hosts))
	
	for i :=0 ;i < len(keys); i++ {
		result = append(result,<-metadata_hosts)
	}
	wg.Wait()
	for _,v := range result {
		addImagesToHostDB(lx,v.Images)
		resultLXD = append(resultLXD,(v.Images)...)
	}
	return SyncResponse(true,resultLXD)
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

func createImageDB(lx *LxdpmApi,fingerprint string) error{
	q := `INSERT INTO images (fingerprint) VALUES (?)`
	_,err := dbExec(lx.db,q,fingerprint)
	return err
}

func addImagesToHostDB(lx *LxdpmApi,images []string) {
	var fingerprint []string
	for _,image := range images {

		fingerprint = strings.Split(image,"/")
		id,err := getImageIdDB(lx,fingerprint[len(fingerprint)-1])

		if err != nil {
			fmt.Println(err)
		} 
		if id == "" {
			err := createImageDB(lx,fingerprint[len(fingerprint)-1])
			if err != nil {
				fmt.Println(err)
			}
		} 
	}
}
func addImageToDB(lx *LxdpmApi,fingerprint string) {
	id,err := getImageIdDB(lx,fingerprint)
	if err != nil {
		fmt.Println(err)
	}
	if id == "" {
		err := createImageDB(lx,fingerprint)
		if err != nil {
			fmt.Println(err)
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
		result,err := imagesPostHost(lx,req,r)
		fmt.Println(result)
		fmt.Println(err)
		return nil
	}

	op, err := taskCreate(taskClassOp, nil, nil, run, nil, nil)
	//fmt.Printf("Soy la op: %+v\n",op)
	fmt.Println("Creation error: ",err)
	tk,err := taskGet(op.id)
	//fmt.Printf("Soy la task: %+v\n",tk)
	//tk.Run() //De normal se llama al hacer el render de operation en Response.go
	fmt.Println("Task Get error: ",err)
	return TaskResponse(tk)
}

func imagesPostHost(lx *LxdpmApi,  imageJson api.ImagesPost, realRequest *http.Request) ([]OperationHost, error) {
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
		go func (key string,freq api.ImagesPost,request *http.Request) {
			defer wg.Done()
			operation_info <- doImagePost(key,freq,request)
		}(key,imageJson,realRequest)
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

func doImagePost(hostname string,r api.ImagesPost, originalrequest *http.Request ) OperationImageInfo {
	requestHeaders := getImagePostRequestHeaders(originalrequest)
	opResult := OperationImageInfo{}
	headers := fmt.Sprintf("-H 'X-LXD-filename: %s' -H 'X-LXD-public: %s' -H 'X-LXD-properties: %s'",requestHeaders["X-LXD-filename"],requestHeaders["X-LXD-public"],requestHeaders["X-LXD-properties"])
	body ,err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
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
		headers := []string{"-H",fmt.Sprintf("'X-LXD-filename: %s'",requestHeaders["X-LXD-filename"]),"-H",fmt.Sprintf("'X-LXD-public: %s'",requestHeaders["X-LXD-public"]),"-H",fmt.Sprintf("'X-LXD-properties: %s'",requestHeaders["X-LXD-properties"])}
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

func getImagePostRequestHeaders(originalrequest *http.Request) map[string]string {
	var resultHeaders map[string]string = make(map[string]string,4)

	resultHeaders["X-LXD-filename"] = originalrequest.Header.Get("X-LXD-filename")
	resultHeaders["X-LXD-public"] = originalrequest.Header.Get("X-LXD-public")
	resultHeaders["X-LXD-properties"] = originalrequest.Header.Get("X-LXD-properties")
	resultHeaders["X-LXD-fingerprint"] = originalrequest.Header.Get("X-LXD-fingerprint")
	fmt.Printf("Parsed headers: %+v",resultHeaders)
	return resultHeaders
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
		fmt.Println("There's at least a failure in: ")
		for _,v := range result {
			if v.Response.Status != "Failure" {
				fmt.Println("deleting: "+v.Response.Metadata["fingerprint"].(string)+" en "+v.Host)
				deleteImage(v)
			}else{
				fmt.Println("No deletion because "+v.Response.Metadata["fingerprint"].(string)+" wasn't created.")
			}
		}
		return nil
	}
	fmt.Println("\nEverything went fine.")
	return result
}

func doWatchImageOperation(op OperationImageInfo) OperationHost {

    out,hostname,err := doWatchOperation(op.Host,op.Operation)
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
	
	hostname := oph.Host
	fingerprint := oph.Response.ID
	_,err := doImageDelete(hostname,fingerprint)
	if err != nil {
		fmt.Println(err)
	}
}

var imageCmd = Command{
	name: "images/{fingerprint}",
	get: imageGet,
	put: imagePut,
	delete: imagesDelete,
}

func getHostnameFromFingerprint(lx *LxdpmApi, fingerprint string) [][]interface{} {
	inargs := []interface{}{}
	outargs := []interface{}{"name"}

	result, err := dbQueryScan(lx.db, `SELECT H.name FROM hosts H, images I where I.host_id = H.id and I.fingerprint='`+fingerprint+`';`,inargs,outargs )
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func imageGet(lx *LxdpmApi,  r *http.Request) Response {
	fingerprint := mux.Vars(r)["fingerprint"]
	hostname := "local"

	res,_ := doImageGet(fingerprint,hostname)
	responseType := getResponseType(res)
	if responseType == "sync" {
		endpointResponse,_ := parseResponseRawToSync(res)
		return &endpointResponse
	}else {
		errorResp := parseErrorResponse(res)
		return &errorResp
	}
}

func imagePut(lx *LxdpmApi,  r *http.Request) Response {
	req := api.ImagePut{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequest(err)
	}

	fingerprint := mux.Vars(r)["fingerprint"]
	hostname := getHostnameFromFingerprint(lx,fingerprint)

	res,_ := doImagePut(fingerprint,hostname[0][0].(string),req)
	responseType := getResponseType(res)
	if responseType == "sync" {
		endpointResponse,_ := parseResponseRawToSync(res)
		return &endpointResponse
	}else {
		errorResp := parseErrorResponse(res)
		return &errorResp
	}
}

func imageDelete(lx *LxdpmApi,  r *http.Request) Response {
	var endpointResponse api.Operation
	var endpointUrl string

	fingerprint := mux.Vars(r)["fingerprint"]
	hostname := getHostnameFromFingerprint(lx,fingerprint)

	res, _ := doImageDelete(hostname[0][0].(string),fingerprint)
	responseType := operationOrError(res)
	if responseType == "operation" {
		endpointUrl,endpointResponse,_ = parseResponseRawToOperation(res)
		return OperationResponse(endpointUrl,&endpointResponse)
	}else{
		errorResp := parseErrorResponse(res)
		return &errorResp
		}
}

func imagesDelete(lx *LxdpmApi,  r *http.Request) Response {
	fingerprint := mux.Vars(r)["fingerprint"]

	run := func(op *task) error {
		result,err := imagesDeleteAllHosts(lx,fingerprint)
		fmt.Printf("Images Delete result:%+v",result)
		fmt.Printf("Images Delete error:%+v",err)
		return nil
	}

	op, err := taskCreate(taskClassOp, nil, nil, run, nil, nil)
	//fmt.Printf("Soy la op: %+v\n",op)
	fmt.Println("Creation error: ",err)
	tk,err := taskGet(op.id)
	//fmt.Printf("Soy la task: %+v\n",tk)
	//tk.Run() //De normal se llama al hacer el render de operation en Response.go
	fmt.Println("Task Get error: ",err)
	return TaskResponse(tk)
}

func imagesDeleteAllHosts(lx *LxdpmApi, fingerprint string) ([]ResponseHost, error) {
	var keys []string
	var wg sync.WaitGroup
	//var resultLXD []string
	
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
		go func (key string,fingerprint string) {
			defer wg.Done()
			out,err := doImageDelete(key,fingerprint)
			if err != nil {
				fmt.Println(err)
			}
			op,err := parseOperationResponse(out)
			if err != nil {
				fmt.Println(err)
			}
			operation_info <- OperationImageInfo{Operation: op.ID, Host: key}
		}(key,fingerprint)
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
	opResults := watchDeleteImageOperation(result)

	//res := imagesPost(req)
	//return AsyncResponse(true,"")
	return opResults,nil
}

func watchDeleteImageOperation(opInfo []OperationImageInfo) []ResponseHost {
	var wg sync.WaitGroup
	var result []ResponseHost
	var operationResponse = make(chan ResponseHost,len(opInfo))
	defer close(operationResponse)
	wg.Add(len(opInfo))
	fmt.Println(opInfo)
	for i,_ := range opInfo {
		operation := opInfo[i] 
		go func (op OperationImageInfo) {
			defer wg.Done()
			operationResponse <- doWatchDeleteImageOperation(op)
		}(operation)
	}
	for i :=0 ;i < len(opInfo); i++ {
		result = append(result,<-operationResponse)
		//fmt.Printf("Soy result: %+v",result)
	}
	wg.Wait()
	//After knowing the result,we can decide if the operation is valid, or we have to do a rollback
	failure := false
	for _,v := range result {
		//fmt.Printf("\nSoy v: %+v",v)
		if v.Response.Type == "error" || v.OperationResponse.Status == "Failure"{
			failure = true
			break
		}
		//resultLXD = append(resultLXD,v)
	}
	if failure == true {
		fmt.Println("There's at least a failure in: ")
		for _,v := range result {
			if v.Response.Type == "error" {
				fmt.Println("No deletion because "+v.Response.Error+" in host "+ v.Host)
			}
			if v.OperationResponse.Status == "Failure" {
				fmt.Println("No deletion because "+v.OperationResponse.Err +"in host "+ v.Host)
			}
		}
		return nil
	}
	fmt.Println("Everything went fine.")
	return result
}

type ResponseHost struct {
	Response 	api.Response
	OperationResponse 	api.Operation
	Host 		string
}

func doWatchDeleteImageOperation(op OperationImageInfo) ResponseHost {

    out,hostname,err := doWatchOperation(op.Host,op.Operation)
    if err != nil {
        fmt.Println(err)
    }
    responseType := operationOrError(out)
	if responseType == "operation" {
		if err != nil {
				fmt.Println(err)
			}
		opresp,_ := parseOperationResponse(out)
		hostResp := ResponseHost{OperationResponse: opresp,Host: hostname}
    	return hostResp
		}else{
		errorResp := parseErrorResponseToApiResponse(out)
		return ResponseHost{Response: errorResp,Host: hostname}
		}
}


var imagesExportCmd = Command{name: "images/{fingerprint}/export", get: imageExport}

func imageExport(lx *LxdpmApi,  r *http.Request) Response {
	fingerprint := mux.Vars(r)["fingerprint"]

	hostname := r.Header.Get("X-LXDPM-hostname")
	if hostname == ""{
		return BadRequest(errors.New("Required header X-LXDPM-hostname."))
	}
	resp := doImageExport(fingerprint,hostname)

	return resp
}

func doImageExport(fingerprint string,hostname string) Response {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket s/1.0/images/"+fingerprint+"/export"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/images/"+fingerprint+"/export"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
    out, err := command.Output()
    if err != nil {
        fmt.Println(err)
    }
    filename := doImageGetFileName(fingerprint,hostname)
    saveFile(out,"./images/"+filename)
    return SyncResponse(true,"Image created.")
}
func doImageGetFileName(fingerprint string,hostname string) string {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket s/1.0/images/"+fingerprint}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/images/"+fingerprint}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
    out, err := command.Output()
    if err != nil {
        fmt.Println(err)
    }
    resp := parseMetadataFromContainerResponse(out)
    filename := resp.Metadata.(map[string]interface{})["filename"].(string)

    return filename
}