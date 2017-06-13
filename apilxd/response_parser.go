package apilxd

import (
	"github.com/lxc/lxd/shared/api"
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

func parseResponseRawToSync(input []byte) (syncResponse,error) {
	var rawResp = api.ResponseRaw{}
	var syncResp = syncResponse{}

    json.NewDecoder(bytes.NewReader(input)).Decode(&rawResp)

    syncResp.success = true

    if rawResp.Response.Status != "Success" {
    	syncResp.success = false
	}
	if rawResp.Metadata != nil {
		syncResp.metadata = rawResp.Metadata
	}

    return syncResp,nil
}

func parseResponseRawToAsync(input []byte) (asyncResponse,error) {
	var rawResp = api.ResponseRaw{}
	var asyncResp = asyncResponse{}

    json.NewDecoder(bytes.NewReader(input)).Decode(&rawResp)
    fmt.Printf("%+v",rawResp)
    asyncResp.success = true

    if rawResp.Response.Status != "Success" {
    	asyncResp.success = false
	}
	if rawResp.Metadata != nil {
		asyncResp.metadata = rawResp.Metadata
	}

    return asyncResp,nil
}

func parseResponseRawToOperation(input []byte) (string,api.Operation,error) {
	var rawResp = api.ResponseRaw{}
	var opResp = api.Operation{}
	var url = ""

    json.NewDecoder(bytes.NewReader(input)).Decode(&rawResp)
    metadata := rawResp.Metadata.(map[string]interface{})
    
    createdAt,err := parseOperationCreatedAt(metadata)
    if err != nil {
    	fmt.Println(err)
    	return url,opResp,err
    }
    updatedAt,err2 := parseOperationUpdatedAt(metadata)
    if err2 != nil {
    	fmt.Println(err2)
    }
    statusCode := parseOperationStatusCode(metadata)
    parsedResources := parseOperationResources(metadata)
    //This way of parsing the value avoids panic, and leaves the uninitialized value on the variable in case it is nil
	metadataValue, _ := metadata["metadata"].(map[string]interface{})
    
    opResp = api.Operation{
		ID:         metadata["id"].(string),
		Class:      metadata["class"].(string),
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		Status:     metadata["status"].(string),
		StatusCode: statusCode,
		Resources:  parsedResources,
		Metadata:   metadataValue,
		MayCancel:  metadata["may_cancel"].(bool),
		Err:        metadata["err"].(string),
	}
    url = rawResp.Response.Operation

    return url,opResp,nil
}

func parseOperationCreatedAt(metadata map[string]interface{}) (time.Time,error) {
	created := metadata["created_at"].(string)
    createdAt,err := time.Parse(time.RFC3339,created)
    if err != nil {
    	fmt.Println(err)
    	return time.Time{},err
    }
    return createdAt,nil
}

func parseOperationUpdatedAt(metadata map[string]interface{}) (time.Time,error) {
	updated := metadata["updated_at"].(string)
    updatedAt,err := time.Parse(time.RFC3339,updated)
    if err != nil {
    	fmt.Println(err)
    	return time.Time{},err
    }
    return updatedAt,nil
}

func parseOperationStatusCode(metadata map[string]interface{}) api.StatusCode {
	statusCodeFloat := metadata["status_code"].(float64)
    statusCodeInt := int(statusCodeFloat)
    statusCode := api.StatusCode(statusCodeInt)
    return statusCode
}

func parseOperationResources(metadata map[string]interface{}) map[string][]string {
	resources := metadata["resources"].(map[string]interface{})
    var parsedResources map[string][]string = make(map[string][]string,len(resources))
    fmt.Printf("%+v\n\n",resources)
    result := []string{}
    for key, value := range resources {
    	switch vv := value.(type) {
    	case []interface{}:
        	fmt.Println(key, "is an array:")
        	for _, u := range vv {
        		result = append(result,u.(string))
        	}
        }
    		parsedResources[key] = result
	}
	return parsedResources
}
func parseErrorResponse(input []byte) errorResponse {
	var resp = api.Response{}
	var errorInfoResponse = errorResponse{}
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    errorInfoResponse.code = resp.Code
    errorInfoResponse.msg = resp.Error 
    
    return errorInfoResponse
}

func operationOrError(input []byte) string{
	var resp = api.Response{}
	json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
	if resp.Type == "error" {
		return "error"
	}
	return "operation"
}

func parseOperationResponse(input []byte) (api.Operation,error) {
	var resp = api.Response{}
	var operationInfo = api.Operation{} 
    json.NewDecoder(bytes.NewReader(input)).Decode(&resp)
    err := resp.MetadataAsStruct(&operationInfo)
    if err != nil {
			fmt.Println(err)
			return operationInfo,err
		}
    return operationInfo,nil
}
