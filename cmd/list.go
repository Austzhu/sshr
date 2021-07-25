package main

import (
	"flag"
	"fmt"
	"sshr/conf"
	"strings"
)

type listConf struct{ regx string }

func doList(c *cmd) error {
	l := conf.List("")
	fmt.Println(strings.Join(l, "\n"))
	return nil
}

func init() {
	name := "list"
	cfg := &listConf{}
	set := flag.NewFlagSet(name, flag.ExitOnError)
	set.StringVar(&cfg.regx, "r", "", "匹配ssh连接的正则")
	register(&cmd{
		name:    name,
		brief:   "查看ssh的配置信息",
		args:    cfg,
		flagSet: set,
		cb:      doList,
	})
}
