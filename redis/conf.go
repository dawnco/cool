package redis

type Cfg struct {
	Name string `json:",optional"`
	Host string
	Pass string `json:",optional"`
}
