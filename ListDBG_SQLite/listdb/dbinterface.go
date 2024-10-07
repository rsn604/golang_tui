package listdb

import (
	"database/sql"
)

// ----------------------------------------------------------------------
type Database interface {
	GetDb() *sql.DB
	GetDropMetaSql() string
	GetDefineSql() string
	GetInsertMetaSql() string
	GetUpdateMetaSql(string) string
	GetDeleteMetaSql(string) string
	GetDefineListSql(string) string
	GetInsertSql(string) string
	GetUpdateSql(string) string
	GetDeleteSql(string) string
	GetLimitSql(countRec int, fromRec int) string
	GetDropSql(string) string
	GetQname(string) string
	GetLastRowId() string
	GetRecordCountSql(string) string
	Connect(string) error
	Close()
}
