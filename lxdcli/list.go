package lxdcli

import (
	"fmt"
	"bytes"
	"log"
	"os/exec"

)

type ListCmd struct {
	cmd 	exec.Cmd
}

func ListCommand() *ListCmd {
	info := &ListCmd{
		cmd: 	*exec.Command("lxc","list"),
	}
	return info 
}

func (c *ListCmd) Do() error {
	var out bytes.Buffer
	c.cmd.Stdout = &out
	err := c.cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%s\n", out.String())

	return err
}