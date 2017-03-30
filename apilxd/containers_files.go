package apilxd

import (
	"net/http"
	"os/exec"
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"log"
	//"encoding/json"
	//"bytes"
	"strings"
	"fmt"
	//"github.com/lxc/lxd/shared/api"
	"github.com/gorilla/mux"
)

func containerFileHandler(lx *LxdpmApi, r *http.Request) Response {
	name := mux.Vars(r)["name"]
	path := r.FormValue("path")

	hostname := getHostnameFromContainername(lx,name)
	switch r.Method {
	case "GET":
			resp := containerFileGetMetadata(hostname[0][0].(string),name,path,r)
			return resp
	case "POST":
			resp := containerFilePost(hostname[0][0].(string),name,path,r)
			return resp
	case "DELETE":
			resp := containerFileDelete(hostname[0][0].(string),name,path,r)
			return resp
	default:
		return NotFound
	}
	
	//resp := containerFilesGet(hostname[0][0].(string),name)

	//meta := resp.Metadata
	return NotFound

}

func containerFileGetMetadata(hostname string, cname string,filepath string,r *http.Request) Response{
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	headers := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket s/1.0/containers/"+cname+"/files?path="+filepath}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket -D - s/1.0/containers/"+cname+"/files?path="+filepath+" -o /dev/null"}
		headers = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/containers/"+cname+"/files?path="+filepath}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-D","-","s/1.0/containers/"+cname+"/files?path="+filepath,"-o","/dev/null"}
		headers = exec.Command("curl", argstr...)
	}
    out, err := command.Output()
    if err != nil {
        fmt.Println(err)
    }
    outhead, err := headers.Output()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Out files: \n"+string(out))
    fmt.Println("Out files: \n"+string(outhead))
    return createFileResponse(string(out),string(outhead),filepath,r)

    //return string(out)+"\n"+string(outhead)
}


func createFileResponse(body string,headers string,path string,r *http.Request) Response {
	var lxdHeaders map[string]string = map[string]string{}

	scanner := bufio.NewScanner(strings.NewReader(headers))
	scanner.Scan()
	for scanner.Scan() {
		splitted := strings.Split(scanner.Text(),": ")
		if strings.HasPrefix(splitted[0],"X-Lxd"){
			lxdHeaders[splitted[0]] = splitted[1]
			fmt.Println(lxdHeaders)
		}
	}
	if err := scanner.Err();err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	temp, err := ioutil.TempFile("", "lxd_forkgetfile_")
	if err != nil {
		return InternalError(err)
	}

	if _, err := temp.Write([]byte(body)); err != nil {
		log.Fatal(err)
	}
	defer temp.Close()

	files := make([]fileResponseEntry, 1)
	files[0].identifier = filepath.Base(path)
	files[0].path = temp.Name()
	files[0].filename = filepath.Base(path)

	return FileResponse(r,files,lxdHeaders,true)

}

func containerFilePost(hostname string, cname string,filepath string, r *http.Request) Response {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	strbody := string(body)
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X POST -H \"Content-Type: application/octet-stream\" -d '"+strbody+"' s/1.0/containers/"+cname+"/files?path="+filepath}
		fmt.Println("\nArgs: ",strings.Join(argstr," "))
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","POST","-H", "Content-Type: application/octet-stream","-d","'"+strbody+"'","s/1.0/containers/"+cname+"/files?path="+filepath}
		fmt.Println("\nArgs: ",strings.Join(argstr," "))
		command = exec.Command("curl", argstr...)
	}
    out, err := command.Output()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Out files: \n"+string(out))
    //Handle this better in case of error
    return SyncResponse(true,"")
}

func containerFileDelete(hostname string, cname string,filepath string, r *http.Request) Response {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X DELETE s/1.0/containers/"+cname+"/files?path="+filepath}
		fmt.Println("\nArgs: ",strings.Join(argstr," "))
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","DELETE","s/1.0/containers/"+cname+"/files?path="+filepath}
		fmt.Println("\nArgs: ",strings.Join(argstr," "))
		command = exec.Command("curl", argstr...)
	}
    out, err := command.Output()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Out files: \n"+string(out))
    //Handle this better in case of error + check if server has required API extension
    return SyncResponse(true,"")
}