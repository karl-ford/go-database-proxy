package database

import (
	_ "github.com/denisenkom/go-mssqldb"
	"net/url"
)

type SQLServer struct {
	*Adapter
}

func CreateSQLServerDialect(adapter *Adapter) *SQLServer{
	dialect := &SQLServer{adapter}
	dialect.delimiters = []byte{0x5b, 0x5d}
	return dialect
}

func (server *SQLServer) CreateAddress() string {
	query := url.Values{}
	query.Add("database", *server.dbname)
	query.Add("connection timeout", "30")
	query.Add("dial timeout", "10")
	query.Add("app name", "data sync to mysql")
	query.Add("TrustServerCertificate", "true")

	uri := &url.URL{
		Scheme:     "sqlserver",
		User:       url.UserPassword(*server.user, *server.pass),
		Host:       *server.uri,
		RawQuery:   query.Encode(),
	}

	return uri.String()
}

func (server *SQLServer) CreateConnection(max uint16) error {
	return server.connect(server.CreateAddress(), max)
}

func (server *SQLServer) GetQuery(builder *QueryBuilder) string {
	return ""
}