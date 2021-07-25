package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sshr/client"
	"sshr/conf"
	"strings"

	. "sshr/public"

	"golang.org/x/term"
)

const SSHR = "SSHR"

type sshConf struct{ conf.SShAuth }

func (sc *sshConf) parse() *sshConf {
	uhost := os.Args[1]
	if uhost == SSHR {
		uhost = os.Args[2]
	}

	uh := strings.Split(uhost, "@")
	if len(uh) != 2 {
		panic("参数错误")
	}

	addr := strings.Split(uh[1], ":")
	host := fmt.Sprintf("%s:22", addr[0])
	if len(addr) == 2 {
		host = fmt.Sprintf("%s:%s", addr[0], addr[1])
	}

	sc.User = uh[0]
	sc.Host = host
	if sc.Passwd == "" {
		sc.Passwd = conf.GetPasswd(sc.SShAuth)
	}

	return sc
}

func doSSh(c *cmd) (err error) {
	defer Recover(&err)()
	sc := c.args.(*sshConf).parse()
	cli, err := client.NewCli(sc.User, sc.Passwd, sc.Host)

	conf.Save(sc.SShAuth)
	Die("SSH Dial", err)
	Die("SSH Terminal", cli.Terminal())
	return
}

func init() {
	cfg := &sshConf{}
	set := flag.NewFlagSet(SSHR, flag.ExitOnError)
	set.StringVar(&cfg.Passwd, "p", "", "ssh登入的密码")
	set.StringVar(&cfg.Brief, "b", "-", "ssh连接的描述信息")
	set.StringVar(&cfg.Group, "g", "-", "ssh连接的group信息")
	register(&cmd{
		name:    SSHR,
		brief:   "-",
		args:    cfg,
		flagSet: set,
		cb:      doSSh,
	})
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	if os.Getenv("SSHR") == "YES" {
		Start()
		return
	}

	Die("Set env", os.Setenv("SSHR", "YES")) // 设置一个环境变量，防止子进程递归执行

	/* 使用子进程来连接ssh, 防止exit时，终端tty设置没有恢复二乱码 */
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	fd := int(os.Stdin.Fd())
	state, err := term.GetState(fd)
	Die("Term GetState", err)
	defer func() { term.Restore(fd, state) }() // 恢复到默认的tty

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	cmd.Process.Kill()
}
