package lxdcli


import (
	"fmt"
	"bytes"
	"os/exec"
	"os"

)

type LaunchPostJSONModel struct {
	Image 	string
	Name 	string
	Args 	string
}

type LaunchCmd struct {
	cmd 	exec.Cmd
}

func LaunchCommand(arg ...string) *LaunchCmd {
	def_args := []string{"launch"}
	args := append(def_args,arg...)
	launch := &LaunchCmd{
		cmd: 	*exec.Command("lxc",args...),
	}
	launch.cmd.Env = os.Environ()
	fmt.Println("Comando:",args)
	fmt.Println("Env:",launch.cmd.Env)
	return launch
}


func (c *LaunchCmd) Do() error {
	var out bytes.Buffer
	var errout bytes.Buffer
	c.cmd.Stdout = &out
	c.cmd.Stderr = &errout

	err := c.cmd.Run()
	if err != nil {
		fmt.Printf("\n%s\n", errout.String())
	}
	fmt.Printf("\n%s\n", out.String())

	return err
}