package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	*Adapter
}

func CreateMySQLDialect(adapter *Adapter) *MySQL{
	dialect := &MySQL{adapter}
	dialect.delimiters = []byte{0x60, 0x60}
	return dialect
}

func (mysql *MySQL) CreateAddress() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&columnsWithAlias=true",
		*mysql.user,
		*mysql.pass,
		*mysql.uri,
		*mysql.dbname,
		*mysql.charset,
		)
}

func (mysql *MySQL) CreateConnection(max uint16) error {
	return mysql.connect(mysql.CreateAddress(), max)
}

func (mysql *MySQL) GetQuery(builder *QueryBuilder) string {
	return ""
}