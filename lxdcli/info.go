package lxdcli

import (
	"fmt"
	"bytes"
	"log"
	"os/exec"

)

type InfoCmd struct {
	cmd 	exec.Cmd
}

func InfoCommand() *InfoCmd {
	info := &InfoCmd{
		cmd: 	*exec.Command("lxc","info"),
	}
	return info 
}

func (c *InfoCmd) Do() error {
	var out bytes.Buffer
	c.cmd.Stdout = &out
	err := c.cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%s\n", out.String())

	return err
}