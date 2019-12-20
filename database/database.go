package database

import (
	"database/sql"
	"fmt"
	"strings"
)

type Adapter struct {
	dialect, uri, user, pass, dbname, charset *string
	maxIdle int
	conn *sql.DB
	delimiters []byte
}

type Dialect interface {
	CreateAddress() string
	CreateConnection(max uint16) error
	CreateQueryBuilder() *QueryBuilder
	CreateQueryBuilderWithTable(table *Table) *QueryBuilder
	GetQuery(builder *QueryBuilder) string
}

func CreateDatabase(dialect, host, user, pass, dbname, charset *string, port *uint16) Dialect {
	*dialect = strings.ToLower(*dialect)
	uri := fmt.Sprintf("%s:%d", *host, *port)
	adapter := &Adapter{
		dialect: dialect,
		uri:     &uri,
		user:    user,
		pass:    pass,
		dbname:  dbname,
		charset: charset,
		maxIdle: 0,
		conn:    nil,
	}
	switch *dialect {
	case "mysql":
		return CreateMySQLDialect(adapter)
	}

	return nil
}

func (adapter *Adapter) CreateQueryBuilder() *QueryBuilder {
	builder := CreateQueryBuilder()
	builder.delimiters = adapter.delimiters
	return builder
}

func (adapter *Adapter) CreateQueryBuilderWithTable(table *Table) *QueryBuilder {
	builder := CreateQueryBuilderForTable(table)
	builder.delimiters = adapter.delimiters
	return builder
}

func (adapter *Adapter) connect(dsn string, max uint16) error {
	conn, err := sql.Open(*adapter.dialect, dsn)
	if err != nil {
		return err
	}

	if max < 1 {
		max = 1
	} else if max > 100 {
		max = 100
	}

	adapter.maxIdle = int(max)
	conn.SetMaxIdleConns(adapter.maxIdle)
	adapter.conn = conn
	return nil
}