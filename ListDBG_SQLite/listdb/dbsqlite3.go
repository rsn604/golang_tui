package listdb

import (
	"database/sql"
	"strconv"
)

// ----------------------------------------------------------------------
// Implement db_interface
// ----------------------------------------------------------------------
type DB_SQLite3 struct {
	Db *sql.DB
}

func (sqlite3 *DB_SQLite3) GetDb() *sql.DB {
	return sqlite3.Db
}

func (sqlite3 *DB_SQLite3) GetDropMetaSql() string {
	return "drop table if exists MetaTable"
}

func (sqlite3 *DB_SQLite3) GetDefineSql() string {
	return "CREATE TABLE MetaTable (id INTEGER PRIMARY KEY AUTOINCREMENT, DBName text, fieldName1 text, fieldName2 text, categoryList text)"
}

func (sqlite3 *DB_SQLite3) GetInsertMetaSql() string {
	return "INSERT INTO MetaTable (DBName, fieldName1, fieldName2, categoryList) VALUES (?, ?, ?, ?)"
}

func (sqlite3 *DB_SQLite3) GetUpdateMetaSql(dbName string) string {
	return "UPDATE MetaTable set fieldName1=?, fieldName2=?, categoryList=? where DBName='" + dbName + "'"
}

func (sqlite3 *DB_SQLite3) GetDeleteMetaSql(dbName string) string {
	return "DELETE FROM MetaTable where DBName = '" + dbName + "'"
}

func (sqlite3 *DB_SQLite3) GetDefineListSql(dbName string) string {
	return "CREATE TABLE '" + dbName + "' (id INTEGER  PRIMARY KEY AUTOINCREMENT, category text, field01 text, field02 text, note text)"
}

func (sqlite3 *DB_SQLite3) GetInsertSql(dbName string) string {
	return "INSERT INTO '" + dbName + "' (category, field01, field02, note) VALUES (?, ?, ?, ?)"
}

func (sqlite3 *DB_SQLite3) GetUpdateSql(dbName string) string {
	return "UPDATE '" + dbName + "' set category=?, field01=?, field02=?, note=? where id="
}

func (sqlite3 *DB_SQLite3) GetDeleteSql(dbName string) string {
	return "DELETE FROM '" + dbName + "' where id="
}

func (sqlite3 *DB_SQLite3) GetLimitSql(countRec int, fromRec int) string {
	return " limit " + strconv.Itoa(countRec) + " offset " + strconv.Itoa(fromRec-1)
}

func (sqlite3 *DB_SQLite3) GetDropSql(dbName string) string {
	return "drop table if exists '" + dbName + "'"
}

func (sqlite3 *DB_SQLite3) GetQname(dbName string) string {
	return "'" + dbName + "'"
}

func (sqlite3 *DB_SQLite3) GetLastRowId() string {
	//return .cursor.lastrowid
	return ""
}

func (sqlite3 *DB_SQLite3) GetRecordCountSql(dbName string) string {
	return "select count(*) from '" + dbName + "'"
}

func (sqlite3 *DB_SQLite3) Connect(connectString string) error {
	var err error
	// @@@ Change here for github.com/mattn/go-sqlite3"->"github.com/glebarez/go-sqlite"
	//sqlite3.Db , err = sql.Open("sqlite3", connectString)
	sqlite3.Db, err = sql.Open("sqlite", connectString)
	return err
}

func (sqlite3 *DB_SQLite3) Close() {
	sqlite3.Db.Close()
}
