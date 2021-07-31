package main

import (
	"flag"
	"fmt"
	"sshr/conf"
	"strings"
)

type listConf struct{ regx string }

func doList(c *cmd) error {
	var l []string
	for _, v := range conf.List("") {
		uh := fmt.Sprintf("%s@%s", v.User, v.Host)
		s := fmt.Sprintf("%-32s %-8s %s", uh, v.Group, v.Brief)
		l = append(l, s)
	}

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
