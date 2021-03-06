package apilxd

import (
	"os/exec"
	"strings"
	"fmt"
	"github.com/lxc/lxd/shared/api"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strconv"
	"regexp"
)

//strings.Join([]string{systemUser,"@",hostname},""),
const systemUser string = "troig" 
const remoteArgs string = "curl -s --unix-socket /var/lib/lxd/unix.socket"
var localArgs = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket"}

func doContainersGet(hostname string) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr := []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" s/1.0/containers"}
		fmt.Println("\nArgs: ",argstr) 
		command = exec.Command("ssh", argstr...)
	}else {
		argstr = append(localArgs,[]string{"s/1.0/containers"}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	out, err := command.Output()
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return out,nil
}
/*
func doContainersPost(req ContainersHostPost) ([]byte,error) {
	buf ,err := json.Marshal(req.ContainersPost)
	if err != nil {
		fmt.Println(err)
	}
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if req.Hostname != "" {
		argstr = []string{strings.Join([]string{systemUser,"@",req.Hostname},""),remoteArgs+" -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers"}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
    out, err := command.Output()
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return out,nil
}*/

func doContainersPlannerPost(req api.ContainersPost,hostname string) ([]byte,error) {
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers"}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
    out, err := command.Output()
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return out,nil
}

func doContainerGet(hostname string, cname string) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" s/1.0/containers/"+cname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,"s/1.0/containers/"+cname)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	out, err := command.Output()
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return out,nil
}

func doContainerPut(hostname string,req api.ContainerPut,cname string) ([]byte,error) {
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" -X PUT -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","PUT","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
	out, err := command.Output()
	if err != nil {
		fmt.Println("fallo al ejecutar el comando: ",err,"out:",out)
		return nil,err
	}

	return out,nil
}

func doContainerDelete(hostname string,cname string) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" -X DELETE s/1.0/containers/"+cname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","DELETE","s/1.0/containers/"+cname}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
}


func doContainerPost(hostname string,req api.ContainerPost,cname string) ([]byte,string,error) {
	argstr := []string{}
	newname := req.Name
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs + " -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,"",err
	}
	return out,newname,nil
}

func doContainerExecPost(hostname string,req api.ContainerExecPost,cname string) ([]byte,error) {
	argstr := []string{}
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	command := exec.Command("curl",argstr...)
	
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),"curl -k --unix-socket /var/lib/lxd/unix.socket -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname+"/exec"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname+"/exec"}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
}

func doContainerStateGet(hostname string, cname string) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" s/1.0/containers/"+cname+"/state"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"s/1.0/containers/"+cname+"/state"}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
}

func doContainerStatePut(hostname string, req api.ContainerStatePut, cname string) ([]byte,error){
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" -X PUT -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname+"/state"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","PUT","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname+"/state"}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
}

func doSnapshotGet(cname string, hostname string, snapname string) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" s/1.0/containers/"+cname+"/snapshots/"+snapname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"s/1.0/containers/"+cname+"/snapshots/"+snapname}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}

	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
	
}

func doSnapshotPost(hostname string, req api.ContainerSnapshotPost, cname string, snap string) ([]byte,error) {
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname+"/snapshots/"+snap}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname+"/snapshots/"+snap}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
}

func doSnapshotDelete(hostname string,cname string,snapname string) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" -X DELETE s/1.0/containers/"+cname+"/snapshots/"+snapname }
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","DELETE","s/1.0/containers/"+cname+"/snapshots/"+snapname }...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
}

func doContainerSnapshotsGet(hostname string, cname string) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" s/1.0/containers/"+cname+"/snapshots"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"s/1.0/containers/"+cname+"/snapshots"}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
	
}

func doContainerSnapshotPost(hostname string,req api.ContainerSnapshotsPost,cname string) ([]byte,error) {
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname+"/snapshots"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname+"/snapshots"}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
}

func containerFileGet(hostname string, cname string, filepath string) ([]byte,[]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	headers := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" s/1.0/containers/"+cname+"/files?path="+filepath}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" -D - s/1.0/containers/"+cname+"/files?path="+filepath+" -o /dev/null"}
		headers = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"s/1.0/containers/"+cname+"/files?path="+filepath}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
		argstr = append(localArgs,[]string{"-D","-","s/1.0/containers/"+cname+"/files?path="+filepath,"-o","/dev/null"}...)
		headers = exec.Command("curl", argstr...)
	}
	out, err := command.Output()
	if err != nil {
		fmt.Println(err)
		return nil,nil,err
	}
	outhead, err := headers.Output()
	if err != nil {
		fmt.Println(err)
		return out,nil,err
	}

	fmt.Println("Out files: \n"+string(out))
	fmt.Println("Out files: \n"+string(outhead))
	return out,outhead,nil

}

func containerFilePost(hostname string, cname string,filepath string, r *http.Request) ([]byte,error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	strbody := string(body)
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" -X POST -H \"Content-Type: application/octet-stream\" -d '"+strbody+"' s/1.0/containers/"+cname+"/files?path="+filepath}
		fmt.Println("\nArgs: ",strings.Join(argstr," "))
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","POST","-H", "Content-Type: application/octet-stream","-d","'"+strbody+"'","s/1.0/containers/"+cname+"/files?path="+filepath}...)
		fmt.Println("\nArgs: ",strings.Join(argstr," "))
		command = exec.Command("curl", argstr...)
	}
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
}

func containerFileDelete(hostname string, cname string,filepath string, r *http.Request) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" -X DELETE s/1.0/containers/"+cname+"/files?path="+filepath}
		fmt.Println("\nArgs: ",strings.Join(argstr," "))
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","DELETE","s/1.0/containers/"+cname+"/files?path="+filepath}...)
		fmt.Println("\nArgs: ",strings.Join(argstr," "))
		command = exec.Command("curl", argstr...)
	}
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
}

func doProfilesGet(hostname string) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr := []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" s/1.0/profiles"}
		fmt.Println("\nArgs: ",argstr) 
		command = exec.Command("ssh", argstr...)
	}else {
		argstr = append(localArgs,[]string{"s/1.0/profiles"}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	out, err := command.Output()
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return out,nil
}

func doProfilesPost(req ProfilesHostPost) ([]byte,error) {
	buf ,err := json.Marshal(req.ProfilesPost)
	if err != nil {
		fmt.Println(err)
	}
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if req.Hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",req.Hostname},""),remoteArgs+" -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/profiles"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/profiles"}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
    out, err := command.Output()
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return out,nil
}

func doProfileGet(hostname string, pname string) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" s/1.0/profiles/"+pname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"s/1.0/profiles/"+pname}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
}

func doProfilePut(hostname string,req api.ProfilePut,pname string) ([]byte,error) {
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" -X PUT -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/profiles/"+pname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","PUT","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/profiles/"+pname}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
}

func doProfileDelete(hostname string,pname string) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" -X DELETE s/1.0/profiles/"+pname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","DELETE","s/1.0/profiles/"+pname}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
}

func doProfilePost(hostname string,req api.ProfilePost,pname string) ([]byte,error) {
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/profiles/"+pname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(localArgs,[]string{"-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/profiles/"+pname}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return nil,err
	}
	return out,nil
}

func doImagesGet(hostname string) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr := []string{strings.Join([]string{systemUser,"@",hostname},""),remoteArgs+" s/1.0/images"}
		fmt.Println("\nArgs: ",argstr) 
		command = exec.Command("ssh", argstr...)
	}else {
		argstr = append(localArgs,[]string{"s/1.0/images"}...)
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	out, err := command.Output()
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return out,nil
}

func doImageGet(fingerprint string,hostname string) ([]byte,error) {
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
		return nil,err
	}
	return out,nil
}

func doImagePut(fingerprint string,hostname string,req api.ImagePut) ([]byte,error) {
	body ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	strbody := string(body)
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket -X PUT -d '"+strbody+"' s/1.0/images/"+fingerprint}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","PUT","-d",strbody,"s/1.0/images/"+fingerprint}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
    out, err := command.Output()
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return out,nil
}

func doImageDelete(hostname string,fingerprint string) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -s --unix-socket /var/lib/lxd/unix.socket -X DELETE s/1.0/images/"+fingerprint}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","-X","DELETE","s/1.0/images/"+fingerprint}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
    out, err := command.Output()
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return out,nil
}

func doGetAvailableMemory(hostname string) (int,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),"bash -c free | awk '/Mem/ { print $4 }'"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-c","free | awk '/Mem/ { print $4 }'"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("bash", argstr...)
	}
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return 0,err
	}
	re := regexp.MustCompile("[0-9]+")
	freeMemory,err := strconv.ParseInt(re.FindAllString(string(out),1)[0],10,0)
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return 0,err
	}
	return int(freeMemory),nil
}

func doGetCores(hostname string) (int,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),"bash -c nproc"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-c","nproc"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("bash", argstr...)
	}
	out, err := command.Output()
	fmt.Println("Este es el out: ",string(out))
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return 0,err
	}
	re := regexp.MustCompile("[0-9]+")
	freeMemory,err := strconv.ParseInt(re.FindAllString(string(out),1)[0],10,0)
	if err != nil {
		fmt.Println("Este es el error: ",err)
		return 0,err
	}
	return int(freeMemory),nil
}

func doWatchOperation(hostname string, ophash string) ([]byte,string,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	//fmt.Println("\n"+string(buf))
	//fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{"troig","@",hostname},""),"curl -k --unix-socket /var/lib/lxd/unix.socket s/1.0/operations/"+ophash+"/wait"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket","s/1.0/operations/"+ophash+"/wait"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("curl", argstr...)
	}
	out, err := command.Output()
	if err != nil {
		fmt.Println(err)
		return nil,"",err
	}
	return out,hostname,nil
}