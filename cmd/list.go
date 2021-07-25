package main

import "flag"

type listConf struct{ regx string }

func doList(c *cmd) error {
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
