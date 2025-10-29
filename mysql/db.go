package mysql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Db struct {
	conn sqlx.SqlConn
}

func (s *Db) Insert(tableName string, data any) (sql.Result, error) {

	typeOf := reflect.TypeOf(data)
	valueOf := reflect.ValueOf(data)

	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}

	if valueOf.Kind() != reflect.Struct {
		return nil, fmt.Errorf("insert expected a struct or struct ptr, got %s", typeOf.Kind())
	}

	// 构造插入语句
	var columns []string
	var placeholders []string
	var values []any

	for i := 0; i < valueOf.NumField(); i++ {
		field := typeOf.Field(i)
		value := valueOf.Field(i)

		columnName := field.Tag.Get("db")
		if columnName == "" {
			columnName = field.Name
		}
		columns = append(columns, columnName)
		placeholders = append(placeholders, "?")
		values = append(values, value.Interface())
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ","),
		strings.Join(placeholders, ","))

	return s.Exec(query, values...)

}

func (s *Db) InsertAndGetId(tableName string, data any) (int64, error) {
	result, err := s.Insert(tableName, data)
	if err != nil {
		return 0, err
	}
	insertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (s *Db) Update(tableName string, data any, conditions map[string]any) (int64, error) {

	typeOf := reflect.TypeOf(data)
	valueOf := reflect.ValueOf(data)

	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}

	if valueOf.Kind() != reflect.Struct {
		return 0, fmt.Errorf("update expected a struct or struct ptr, got %s", typeOf.Kind())
	}

	var args []any

	// 构建 SET 子句
	var setClauses []string
	for i := 0; i < valueOf.NumField(); i++ {
		value := valueOf.Field(i)
		field := typeOf.Field(i)
		if value.CanInterface() {

			columnName := field.Tag.Get("db")
			if columnName == "" {
				columnName = field.Name
			}

			setClauses = append(setClauses, fmt.Sprintf("%s = ?", columnName))
			args = append(args, value.Interface())
		}
	}
	setClause := strings.Join(setClauses, ", ")

	// 构建 WHERE 子句
	var whereClauses []string
	for key, value := range conditions {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " " + strings.Join(whereClauses, " AND ")
	}

	// 构建 SQL 语句
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, setClause, whereClause)

	result, err := s.conn.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()

}

func (s *Db) Delete(tableName string, conditions map[string]any) (int64, error) {

	var args []any

	// 构建 WHERE 子句
	var whereClauses []string
	for key, value := range conditions {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " " + strings.Join(whereClauses, " AND ")
	}

	// 构建 SQL 语句
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, whereClause)

	// 设置参数并执行更新
	result, err := s.conn.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s *Db) GetRow(v any, query string, args ...any) error {
	return s.conn.QueryRow(v, query, args...)
}

func (s *Db) GetData(v any, query string, args ...any) error {
	return s.conn.QueryRows(v, query, args...)
}
func (s *Db) Exec(query string, args ...any) (sql.Result, error) {
	result, err := s.conn.Exec(query, args...)
	if err == nil {
		return result, nil
	} else {
		return nil, fmt.Errorf("执行SQL错误 %s sql:%s", err.Error(), query)
	}
}

var dbWrite *Db
var dbRead *Db

var onceWrite sync.Once
var onceRead sync.Once

func GetDbWrite() *Db {
	onceWrite.Do(func() {
		dbWrite = &Db{
			conn: connectionWrite,
		}
	})
	return dbWrite
}

func GetDbRead() *Db {
	onceRead.Do(func() {
		dbRead = &Db{
			conn: connectionRead,
		}
	})
	return dbRead
}
