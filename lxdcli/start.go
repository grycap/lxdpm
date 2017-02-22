package lxdcli

import (
	"fmt"
	"bytes"
	//"log"
	"os/exec"
	"strings"

)

type StartCmd struct {
	cmd 	exec.Cmd
}

func StartCommand(arg ...string) *StartCmd {
	start := &StartCmd{
		cmd: 	*exec.Command("lxc","start",strings.Join(arg," ")),
	}
	return start 
}

func (c *StartCmd) Do() error {
	var out bytes.Buffer
	var errout bytes.Buffer
	c.cmd.Stdout = &out
	c.cmd.Stderr = &errout
	err := c.cmd.Run()
	if err != nil {
		fmt.Printf("%s\n", errout.String())
	}
	fmt.Printf("Test: %s\n", out.String())

	return err
}