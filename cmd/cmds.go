package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	. "github.com/zhuzongzhen/sshr/public"
)

type cmd struct {
	name    string
	brief   string
	args    interface{}
	flagSet *flag.FlagSet
	cb      func(*cmd) error
}

var sCmd = &struct {
	sync.RWMutex
	c []*cmd
}{}

func usage(_ *cmd) error {
	title := filepath.Base(os.Args[0]) + "{CMD | USER@IPADDR[:PORT]} [OPTIONS]"
	cmds := "CMD:\n"
	opts := "OPTIONS:"

	for _, v := range sCmd.c {
		if v.name != SSHR {
			cmds += fmt.Sprintf("    %-12s %s\n", v.name, v.brief)
		} else {
			out := bytes.NewBuffer(nil)
			v.flagSet.SetOutput(out)
			v.flagSet.Usage()
			r := regexp.MustCompile("Usage of .*")
			opts += "    " + r.ReplaceAllString(out.String(), "") + "\n"
		}
	}

	fmt.Println(title)
	fmt.Println(cmds)
	fmt.Println(opts)
	return nil
}

func Usage() { usage(nil) }

func register(c *cmd) {
	sCmd.Lock()
	defer sCmd.Unlock()
	sCmd.c = append(sCmd.c, c)
}

func getCmd(name string) *cmd {
	sCmd.RLock()
	defer sCmd.RUnlock()

	for _, v := range sCmd.c {
		if v.name == name {
			return v
		}
	}
	return nil
}

func Start() {
	c := getCmd(os.Args[1])
	if c == nil {
		c = getCmd(SSHR)
	}

	if c.flagSet != nil {
		Die("flagSet Parse", c.flagSet.Parse(os.Args[2:]))
	}

	if err := c.cb(c); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
