package db

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var connectionRead sqlx.SqlConn
var connectionWrite sqlx.SqlConn

func initConnection(cRead *Cfg, cWrite *Cfg) {

	sqlx.DisableLog()
	sqlx.DisableStmtLog()

	if cRead != nil {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=false&loc=UTC&time_zone=%s",
			cRead.User,
			cRead.Pass,
			cRead.Host,
			cRead.Port,
			cRead.Name,
			cRead.Charset,
			url.QueryEscape(fmt.Sprintf("'%s'", cRead.Zone)),
		)
		connectionRead = sqlx.NewMysql(dsn)
	}

	if cWrite != nil {
		dsn2 := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=false&loc=UTC&time_zone=%s",
			cWrite.User,
			cWrite.Pass,
			cWrite.Host,
			cWrite.Port,
			cWrite.Name,
			cWrite.Charset,
			url.QueryEscape(fmt.Sprintf("'%s'", cWrite.Zone)),
		)
		connectionWrite = sqlx.NewMysql(dsn2)
	}
}

func Init(cRead *Cfg, cWrite *Cfg) {
	sync.OnceFunc(func() {
		initConnection(cRead, cWrite)
	})()
}
