package wredis

type Cfg struct {
	Name     string `json:",default=default,optional"`
	Addr     string `json:",default=127.0.0.1,optional"`
	Port     int    `json:",default=6379,optional"`
	Password string `json:",optional"`
	Db       int    `json:",default=1,optional"`
}
