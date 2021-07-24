package main

import (
	"log"
	"os"

	"sshr/conf"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func die(msg string, err error) {
	if err != nil {
		log.Fatalf("%s fail: %v\n", msg, err)
	}
}

func warn(msg string, err error) {
	if err != nil {
		log.Printf("%s fail: %v\n", msg, err)
	}
}

func main() {
	var (
		c    = conf.GetSshConf()
		auth = []ssh.AuthMethod{ssh.Password(c.Passwd)}
		cb   = ssh.InsecureIgnoreHostKey()
		cc   = &ssh.ClientConfig{User: c.User, Auth: auth, HostKeyCallback: cb}
	)

	cli, err := ssh.Dial("tcp", c.Host, cc)
	die("SSH dial", err)

	session, err := cli.NewSession() // 建立会话
	die("New session", err)
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	die("Terminal MakeRaw", err)
	defer term.Restore(fd, oldState)

	session.Stdout = os.Stdout
	session.Stderr = os.Stdin
	session.Stdin = os.Stdin

	w, h, err := term.GetSize(fd)
	warn("Terminal GetSize", err)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}

	die("Request pty", session.RequestPty("xterm-256color", h, w, modes))
	die("Start shell", session.Shell())
	die("Ssh Wait", session.Wait())
}
