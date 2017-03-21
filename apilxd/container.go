package apilxd

import (
	"net/http"
	/*"os/exec"
	"encoding/json"
	"bytes"
	"strings"*/
	"fmt"
	//"github.com/gorilla/mux"
)
const DB_FILL string = `
	INSERT INTO hosts (id, name, ip) VALUES (1,'localhost','');
	`
func containerGet(lx *LxdpmApi, r *http.Request) Response {
	/*_, err := lx.db.Exec(DB_FILL)
	if err != nil {
		return BadRequest(err)
	}*/
	result := getDBhosts(lx)
	return SyncResponse(true,result)
}

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
/*func containerGet(lx *LxdpmApi, r *http.Request) Response {
	name := mux.Vars(r)["name"]

	resp := containerGetMetadata("lxdpm02",name)

	meta := resp.Metadata 

	return SyncResponse(true,meta)
}

func containerGetMetadata(hostname string, contname string) LxdResponseRaw {

	argstr := []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket s/1.0/containers/"+contname}  
    out, err := exec.Command("ssh", argstr...).Output()
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
}*/