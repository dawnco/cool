package mysql

// Cfg 数据库配置
type Cfg struct {
	User    string
	Pass    string
	Host    string
	Port    int
	Name    string
	Zone    string // 格式 "+08:00" mysql 连接使用的
	TimeLoc string // 格式 Asia/Shanghai  mysql日期格式转成  time.Time使用的的时区
	Charset string
}
