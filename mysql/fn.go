package mysql

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var instance = sync.Map{}

// Init 配置
// name 配置名称 后面通过这个名称获取客户端
// cfg 配置参数
func Init(name string, cfg Cfg) {

	sqlx.DisableLog()
	sqlx.DisableStmtLog()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=false&loc=UTC&time_zone=%s",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Charset,
		url.QueryEscape(fmt.Sprintf("'%s'", cfg.Zone)),
	)
	instance.Store(name, &Db{conn: sqlx.NewMysql(dsn)})

}

func Get(name string) *Db {
	conn, ok := instance.Load(name)
	if !ok {
		panic(fmt.Errorf("mysql connection %s not found", name))
	}
	return conn.(*Db)
}
