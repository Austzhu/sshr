package client

import (
	"fmt"
	"os"
	"runtime"

	. "sshr/public"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type Client struct {
	User    string
	Passwd  string
	Host    string
	scli    *ssh.Client
	session *ssh.Session
	fd      int
	state   *term.State
}

func recv(err *error) func() {
	return func() {
		if rec := recover(); rec != nil {
			*err = fmt.Errorf("%v", rec)
		}
	}
}

func close(c *Client) {
	if c.session != nil {
		Warn("Session close", c.session.Close())
	}

	if c.scli != nil {
		Warn("ssh cli close", c.scli.Close())
	}

	if c.state != nil {
		Warn("Term restore", term.Restore(c.fd, c.state))
	}

}

func NewCli(user, pwd, host string) (cli *Client, err error) {
	defer recv(&err)() // 错误恢复

	c := &Client{User: user, Passwd: pwd, Host: host}
	c.scli, err = ssh.Dial("tcp", host, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	Die("SSH dial", err)

	c.session, err = c.scli.NewSession() // 建立会话
	Die("New session", err)

	runtime.SetFinalizer(c, close) // 设置垃圾自动回收时关闭资源

	cli = c
	return
}

func (c *Client) Terminal() (err error) {
	defer recv(&err)()

	session := c.session
	c.fd = int(os.Stdin.Fd())
	c.state, err = term.MakeRaw(c.fd)
	Die("Terminal MakeRaw", err)

	session.Stdout = os.Stdout
	session.Stderr = os.Stdin
	session.Stdin = os.Stdin

	w, h, err := term.GetSize(c.fd)
	Warn("Terminal GetSize", err)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}

	Die("Request pty", session.RequestPty("xterm-256color", h, w, modes))
	Die("Start shell", session.Shell())
	Die("Ssh Wait", session.Wait())
	return
}
