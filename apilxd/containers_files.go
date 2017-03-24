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

	resp := containerFileGetMetadata(hostname[0][0].(string),name,path,r)
	//resp := containerFilesGet(hostname[0][0].(string),name)

	//meta := resp.Metadata 

	return resp
}
/*func containerFileGet(hostname string, cname string,path string r *http.Request) Response {

}*/
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

	//defer os.Remove(temp.Name())

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