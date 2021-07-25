package main

import (
	"flag"
	"fmt"
)

type scpConf struct{}

func doScp(c *cmd) error {
	return fmt.Errorf("Not supported yet!")
}

func init() {
	name := "scp"
	cfg := &scpConf{}
	set := flag.NewFlagSet(name, flag.ExitOnError)

	register(&cmd{
		name:    name,
		brief:   "scp拷贝文件",
		args:    cfg,
		flagSet: set,
		cb:      doScp,
	})
}
