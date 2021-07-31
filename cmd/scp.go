package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sshr/client"
	"sshr/conf"
	. "sshr/public"
	"strings"

	"github.com/pkg/sftp"
)

type scpConf struct{}
type scPath struct {
	conf.SShAuth
	path string
}

func scpUsage() {
	fmt.Println("Usage of scp:\n    scp [USER@IP:PORT]<SRC> [USER@IP:PORT]<DST>")
	fmt.Println("    eg: scp root@127.0.0.1:22:/root/test  /root/abc")
}

func getScPath(str string) *scPath {
	s := strings.Split(str, ":")
	r := &scPath{}
	p := "22"

	switch len(s) {
	case 1:
		return &scPath{path: s[0]}

	case 2:
		r.path = s[1]

	case 3:
		p = s[1]
		r.path = s[2]

	default:
		return nil
	}

	t := strings.Split(s[0], "@")
	if len(t) != 2 {
		return nil
	}

	r.User = t[0]
	r.Host = t[1] + ":" + p
	return r
}

func isDir(cli *sftp.Client, name string) bool {
	if name == "~" {
		return true
	}

	if name[len(name)-1:] == "/" {
		return true
	}

	info, err := cli.Stat(name)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// 从远程拷贝到本地
func scpFromRemote(cli *sftp.Client, s, d string) (err error) {
	defer Recover(&err)()

	dpath := d
	info, err := os.Stat(d)
	if err == nil && info.IsDir() {
		dpath = filepath.Join(d, filepath.Base(s))
	}

	Die("MkdirAll", os.MkdirAll(filepath.Dir(dpath), 0755))

	if s[:2] == "~/" {
		s = s[2:]
	}

	sf, err := cli.OpenFile(s, os.O_RDONLY)
	Die("sftp Create", err)
	defer sf.Close()

	info, _ = sf.Stat()
	of, err := os.OpenFile(d, os.O_RDWR|os.O_CREATE, info.Mode())
	Die("os OpenFile", err)
	defer of.Close()

	_, err = io.Copy(of, sf)
	Die("io.Copy", err)
	return
}

// 从本地拷贝到远程
func scpToRemote(cli *sftp.Client, s, d string) (err error) {
	defer Recover(&err)()

	rpath := d
	if isDir(cli, rpath) {
		rpath = filepath.Join(d, filepath.Base(s))
	}

	if rpath[:2] == "~/" {
		rpath = rpath[2:]
	}

	cli.MkdirAll(filepath.Dir(rpath))

	of, err := os.Open(s)
	Die("os open", err)
	defer of.Close()

	info, _ := of.Stat()

	sf, err := cli.Create(rpath)
	Die("sftp Create", err)
	defer sf.Close()

	sf.Chmod(info.Mode())
	_, err = io.Copy(sf, of)
	Die("io.Copy", err)
	return
}

func scp(s, d *scPath) (err error) {
	defer Recover(&err)()

	var fn func(c *sftp.Client, s, d string) error
	var au conf.SShAuth

	if d.User != "" {
		d.Passwd = conf.GetPasswd(d.SShAuth)
		au = d.SShAuth
		fn = scpToRemote
	} else {
		s.Passwd = conf.GetPasswd(s.SShAuth)
		au = s.SShAuth
		fn = scpFromRemote
	}

	cli, err := client.NewCli(au.User, au.Passwd, au.Host)
	Die("NewCli", err)

	scli, err := sftp.NewClient(cli.SShCli())
	Die("sftp NewClient", err)
	defer scli.Close()

	return fn(scli, s.path, d.path)
}

func doScp(c *cmd) error {
	if len(os.Args) < 4 {
		return fmt.Errorf("参数错误")
	}

	src := getScPath(os.Args[2])
	dst := getScPath(os.Args[3])
	if src == nil || dst == nil || (src.User != "" && dst.User != "") {
		return fmt.Errorf("参数错误")
	}

	return scp(src, dst)
}

func init() {
	name := "scp"
	cfg := &scpConf{}
	set := flag.NewFlagSet(name, flag.ExitOnError)
	set.Usage = scpUsage

	register(&cmd{
		name:    name,
		brief:   "scp拷贝文件",
		args:    cfg,
		flagSet: set,
		cb:      doScp,
	})
}
