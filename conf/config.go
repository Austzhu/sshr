package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type SShAuth struct {
	User   string // 用户名
	Passwd string // 密码
	Host   string // ssh地址, eg:127.0.0.1:22
	Brief  string // 简单描述信息
	Group  string // 配置的group
}

type SShConf struct {
	sync.RWMutex
	Auth     []SShAuth
	fileName string // 配置文件名称
}

var sconf = &SShConf{}

func (c *SShConf) GetAuth(au SShAuth) (*SShAuth, bool) {
	c.RLock()
	defer c.RUnlock()

	for i, v := range c.Auth {
		if v.User == au.User && v.Host == au.Host {
			if au.Passwd == "" || v.Passwd == au.Passwd {
				needSave := false

				if au.Group != "-" {
					needSave = true
					c.Auth[i].Group = au.Group
				}

				if au.Brief != "-" {
					needSave = true
					c.Auth[i].Brief = au.Brief
				}

				return &v, needSave
			}
		}
	}

	return nil, true
}

func dedup(list []string) []string {
	var dump = make(map[string]bool)
	var resp []string
	for _, v := range list {
		if ok := dump[v]; ok {
			continue
		}

		resp = append(resp, v)
		dump[v] = true
	}

	return resp
}

func GetPasswd(au SShAuth) string {
	a, _ := sconf.GetAuth(au)
	if a == nil {
		return ""
	}

	return a.Passwd
}

func (c *SShConf) genList(g string) []string {
	var slist []string
	for _, v := range c.Auth {
		uh := fmt.Sprintf("%s@%s", v.User, v.Host)
		s := fmt.Sprintf("%-32s %-8s %-8s %s", uh, v.Passwd, v.Group, v.Brief)
		slist = append(slist, s)
	}

	return dedup(slist)
}

func List(g string) []string {
	sconf.RLock()
	defer sconf.RUnlock()
	return sconf.genList(g)
}

func Save(au SShAuth) {
	c := sconf
	auth, needSave := c.GetAuth(au)

	if auth == nil {
		c.Lock()
		defer c.Unlock()
		c.Auth = append(c.Auth, au)
	}

	if needSave {
		str := strings.Join(c.genList(""), "\n") + "\n"
		ioutil.WriteFile(c.fileName, []byte(str), 0644)
	}
}

func oneParse(text string) *SShAuth {
	l := strings.Fields(text)
	if len(l) < 4 {
		return nil
	}

	u := strings.Split(l[0], "@")
	if len(u) != 2 {
		return nil
	}

	return &SShAuth{User: u[0], Passwd: l[1], Host: u[1], Brief: l[3], Group: l[2]}
}

func confParse(text string) (auth []SShAuth) {
	text = strings.Trim(text, " ")
	if text == "" {
		return
	}

	list := strings.Split(text, "\n")
	for _, v := range list {
		if tmp := oneParse(v); tmp != nil {
			auth = append(auth, *tmp)
		}
	}

	return
}

func init() {
	c := sconf
	home := os.Getenv("HOME")
	if runtime.GOOS != "linux" {
		home = os.Getenv("USERPROFILE")
	}

	c.fileName = filepath.Join(home, ".sshr/sshr.conf")
	os.MkdirAll(filepath.Dir(c.fileName), 0755)

	buf, _ := ioutil.ReadFile(c.fileName)
	if len(buf) > 0 {
		c.Auth = confParse(string(buf))
	}
}
