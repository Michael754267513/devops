/*
		Handpay ServiceMesh

           创建时间: 2020年11月25日15:55:24

	       少侠好武功,一起Giao起来
	  	 我说一Giao,你说Giao
		   一 Giao ？？？？

*/

package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"gitlab.handpay.com.cn/seed-sql/global"
)

type Mysql struct {
	DSN string `json:"dsn"`
	DB  *sql.DB
}

func (m *Mysql) Exec(sql string) (res interface{}, err error) {

	res, err = m.DB.Exec(sql)
	return
}

// 动态表格可以实现查询
func (m *Mysql) Query(sql string) (res interface{}, err error) {
	rows, err := m.DB.Query(sql)
	if err != nil {
		return
	}
	cols, err := rows.Columns()
	if err != nil {
		return
	}
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		rows.Scan(scans...)
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		result[i] = row
		i++
	}

	return result, err
}

func (m *Mysql) Connection(dsn string) (err error) {
	m.DSN = dsn
	m.DB, err = sql.Open("mysql", m.DSN)
	if err != nil {
		global.Log.WithError(err)
		return
	}
	err = m.DB.Ping()
	if err != nil {
		global.Log.WithError(err)
	}
	return
}

func (m *Mysql) Close() {
	m.DB.Close()
}
