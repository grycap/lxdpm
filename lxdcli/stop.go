package lxdcli

import (
	"fmt"
	"bytes"
	//"log"
	"os/exec"
	"strings"
)

type StopCmd struct {
	cmd 	exec.Cmd
}

func StopCommand(arg ...string) *StopCmd {
	start := &StopCmd{
		cmd: 	*exec.Command("lxc","stop",strings.Join(arg," ")),
	}
	return start 
}

func (c *StopCmd) Do() error {
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