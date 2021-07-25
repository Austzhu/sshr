package conf

import (
	"flag"
)

type SShConf struct {
	User   string // 用户名
	Passwd string // 密码
	Host   string // ssh地址, eg:127.0.0.1:22
	Brief  string // 简单描述信息
	Group  string // 配置的group
}

var sc = &SShConf{}

func Init() {
	flag.StringVar(&sc.User, "u", "", "用户名")
	flag.StringVar(&sc.Passwd, "p", "", "密码")
	flag.StringVar(&sc.Brief, "b", "-", "描述信息")
	flag.StringVar(&sc.Host, "h", "127.0.0.1:22", "ssh连接的host信息")
	flag.Parse()

	if sc.User == "" {
		panic("User is empty!")
	}
}

func GetSshConf() *SShConf { return sc }
