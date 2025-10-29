package mysql

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/dawnco/cool/env"
	"github.com/stretchr/testify/assert"
)

/*

CREATE TABLE `test` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `val` int DEFAULT NULL,
  `money` decimal(10,2) DEFAULT NULL,
  `sn` bigint DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

*/

func TestDb(t *testing.T) {

	var rc, wc *Cfg

	if env.Get("MYSQL_READ_HOST", "") != "" {
		rc = &Cfg{
			Host:    env.Get("MYSQL_READ_HOST", ""),
			Port:    env.Get("MYSQL_READ_PORT", 3306),
			User:    env.Get("MYSQL_READ_USER", "root"),
			Pass:    env.Get("MYSQL_READ_PASS", "root"),
			Name:    env.Get("MYSQL_READ_NAME", "test"),
			Zone:    env.Get("MYSQL_READ_ZONE", "+08:00"),
			Charset: env.Get("MYSQL_READ_CHAR", "utf8mb4"),
		}
	}
	if env.Get("MYSQL_WRITE_HOST", "") != "" {
		wc = &Cfg{
			Host:    env.Get("MYSQL_WRITE_HOST", ""),
			Port:    env.Get("MYSQL_READ_PORT", 3306),
			User:    env.Get("MYSQL_WRITE_USER", "root"),
			Pass:    env.Get("MYSQL_WRITE_PASS", "root"),
			Name:    env.Get("MYSQL_WRITE_NAME", "test"),
			Zone:    env.Get("MYSQL_WRITE_ZONE", "+08:00"),
			Charset: env.Get("MYSQL_WRITE_CHAR", "utf8mb4"),
		}
	}

	if env.Get("MYSQL_READ_HOST", "") == "" {
		t.Errorf("环境变量为空")
		return
	}
	if env.Get("MYSQL_WRITE_HOST", "") == "" {
		t.Errorf("环境变量为空")
		return
	}

	Init(rc, wc)

	db := GetDbWrite()

	type row struct {
		Name  string  `db:"name"`
		Val   int     "db:\"`val`\""
		Money float64 `db:"money"`
		Sn    int64   `db:"sn"`
	}

	type rowGet struct {
		row
		ID int `db:"id"`
	}

	insertRow := row{
		Name:  "n1",
		Val:   9,
		Money: 3.12,
		Sn:    123,
	}

	insertRow.Name = "n1" + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	_, err := db.Insert("test", &insertRow)
	assert.Equal(t, nil, err)

	insertRow.Name = "nn1"

	id, err := db.InsertAndGetId("test", insertRow)
	assert.Equal(t, nil, err)

	if id <= 0 {
		t.Error("自增ID不能为小于0")
	}

	findRow := rowGet{}
	fmt.Println(err)
	db.GetRow(&findRow, "SELECT * FROM test WHERE name = ?", insertRow.Name)
	assert.Equal(t, insertRow.Name, findRow.Name)

	updateRow := row{
		Name:  "nUpdate",
		Val:   9,
		Money: 3.12,
		Sn:    123,
	}
	_, err2 := db.Update("test", updateRow, map[string]any{"name": "nn1"})
	assert.Equal(t, nil, err2)

	_, err3 := db.Delete("test", map[string]any{
		"name": "nUpdate",
	})
	assert.Equal(t, nil, err3)
}
