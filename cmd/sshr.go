package main

import (
	"os"
	"os/exec"
	"runtime"
	"sshr/client"
	"sshr/conf"

	. "sshr/public"

	"golang.org/x/term"
)

func start() {
	conf.Init()
	c := conf.GetSshConf()
	cli, err := client.NewCli(c.User, c.Passwd, c.Host)
	Die("SSH Dial", err)
	Die("SSH Terminal", cli.Terminal())
	cli = nil
	runtime.GC()
}

func main() {
	if os.Getenv("SSHR") == "YES" {
		start()
		return
	}

	Die("Set env", os.Setenv("SSHR", "YES")) // 设置一个环境变量，防止子进程递归执行

	/* 使用子进程来连接ssh, 防止exit时，终端tty设置没有恢复二乱码 */
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	Die("Term MakeRaw", err)
	defer func() { term.Restore(fd, state) }() // 恢复到默认的tty

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	Die("CMD RUN", cmd.Run())
	cmd.Process.Kill()
}
