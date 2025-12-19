package mysql

// Cfg 数据库配置
type Cfg struct {
	User    string `json:",default=root,optional"`
	Pass    string `json:",default=root,optional"`
	Host    string `json:",default=localhost,optional"`
	Port    int    `json:",default=3306,optional"`
	Name    string `json:",default=db_name,optional"`
	Zone    string `json:",default=+08:00,optional"`        // 格式 "+08:00" mysql 连接使用的
	TimeLoc string `json:",default=Asia/Shanghai,optional"` // 格式 Asia/Shanghai  mysql日期格式转成  time.Time使用的的时区
	Charset string `json:",default=utf8mb4,optional"`
}
