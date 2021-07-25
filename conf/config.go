package conf

type SShConf struct {
	User   string // 用户名
	Passwd string // 密码
	Host   string // ssh地址, eg:127.0.0.1:22
	Brief  string // 简单描述信息
	Group  string // 配置的group
}
