package main

import (
	"flag"
	"fmt"
)

type importConf struct{}

func doImport(c *cmd) error {
	return fmt.Errorf("Not supported yet!")
}

func init() {
	name := "import"
	cfg := &importConf{}
	set := flag.NewFlagSet(name, flag.ExitOnError)
	register(&cmd{
		name:    name,
		brief:   "导入ssh连接的配置信息",
		args:    cfg,
		flagSet: set,
		cb:      doImport,
	})
}
