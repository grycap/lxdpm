package apilxd

import (
	"os/exec"
	"strings"
	"fmt"
	"github.com/lxc/lxd/shared/api"
	"encoding/json"
)

//strings.Join([]string{systemUser,"@",hostname},""),
const systemUser string = "troig" 
const localArgs string = "curl -s --unix-socket /var/lib/lxd/unix.socket"
var remoteArgs = []string{"-k","--unix-socket","/var/lib/lxd/unix.socket"}

type CommandExecutor struct {

}

func doContainerGet(hostname string, cname string) ([]byte,error) {
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),localArgs+"s/1.0/containers/"+cname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(remoteArgs,"s/1.0/containers/"+cname)
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
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),localArgs+" -X PUT -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(remoteArgs,[]string{"-X","PUT","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname}...)
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
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),localArgs+" -X DELETE s/1.0/containers/"+cname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(remoteArgs,[]string{"-X","DELETE","s/1.0/containers/"+cname}...)
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


func doContainerPost(hostname string,req api.ContainerPost,cname string) ([]byte,error) {
	argstr := []string{}
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),localArgs + " -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(remoteArgs,[]string{"-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname}...)
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
		argstr = append(remoteArgs,[]string{"-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname+"/exec"}...)
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
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),localArgs+" s/1.0/containers/"+cname+"/state"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(remoteArgs,[]string{"s/1.0/containers/"+cname+"/state"}...)
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
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),localArgs+" -X PUT -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname+"/state"}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(remoteArgs,[]string{"-X","PUT","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname+"/state"}...)
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
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),localArgs+" s/1.0/containers/"+cname+"/snapshots/"+snapname}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(remoteArgs,[]string{"s/1.0/containers/"+cname+"/snapshots/"+snapname}...)
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

func doSnapshotPost(hostname string, req api.ContainerSnapshotsPost, cname string, snap string) ([]byte,error) {
	buf ,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	argstr := []string{}
	command := exec.Command("curl",argstr...)
	fmt.Println("\n"+string(buf))
	fmt.Println("\n"+fmt.Sprintf("'"+string(buf)+"'"))
	if hostname != "local" {
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),localArgs+" -X POST -d "+fmt.Sprintf("'"+string(buf)+"'")+" s/1.0/containers/"+cname+"/snapshots/"+snap}
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(remoteArgs,[]string{"-X","POST","-d",fmt.Sprintf(""+string(buf)+""),"s/1.0/containers/"+cname+"/snapshots/"+snap}...)
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
		argstr = []string{strings.Join([]string{systemUser,"@",hostname},""),localArgs+" -X DELETE s/1.0/containers/"+cname+"/snapshots/"+snapname }
		fmt.Println("\nArgs: ",argstr)
		command = exec.Command("ssh", argstr...)
	} else {
		argstr = append(remoteArgs,[]string{"-X","DELETE","s/1.0/containers/"+cname+"/snapshots/"+snapname }...)
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
