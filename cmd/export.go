package main

import (
	"flag"
	"fmt"
)

type exportConf struct{}

func doExport(c *cmd) error {
	return fmt.Errorf("Not supported yet!")
}

func init() {
	name := "export"
	cfg := &exportConf{}
	set := flag.NewFlagSet(name, flag.ExitOnError)

	register(&cmd{
		name:    name,
		brief:   "导出ssh连接的配置信息",
		args:    cfg,
		flagSet: set,
		cb:      doExport,
	})
}
