package lxdcli


import (
	"fmt"
	"bytes"
	"os/exec"
	"strings"

)

type DeleteCmd struct {
	cmd 	exec.Cmd
}

func DeleteCommand(arg ...string) *DeleteCmd {
	del := &DeleteCmd{
		cmd: 	*exec.Command("lxc","delete",strings.Join(arg," ")),
	}
	return del
}


func (c *DeleteCmd) Do() error {
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