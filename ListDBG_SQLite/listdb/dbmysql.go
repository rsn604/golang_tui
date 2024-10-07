package listdb

import (
	"database/sql"
	"strconv"
)

// ----------------------------------------------------------------------
// Implement db_interface
// ----------------------------------------------------------------------
type DB_MySQL struct {
	Db *sql.DB
}

func (mysql *DB_MySQL) GetDb() *sql.DB {
	return mysql.Db
}

func (mysql *DB_MySQL) GetDropMetaSql() string {
	return "drop table if exists MetaTable"
}

func (mysql *DB_MySQL) GetDefineSql() string {
	return "CREATE TABLE MetaTable (id INTEGER PRIMARY KEY AUTO_INCREMENT , DBName text, fieldName1 text, fieldName2 text, categoryList text)"
}

func (mysql *DB_MySQL) GetInsertMetaSql() string {
	return "INSERT INTO MetaTable (DBName, fieldName1, fieldName2, categoryList) VALUES (?, ?, ?, ?)"
}

func (mysql *DB_MySQL) GetUpdateMetaSql(dbName string) string {
	return "UPDATE MetaTable set fieldName1=?, fieldName2=?, categoryList=? where DBName=`" + dbName + "`"
}

func (mysql *DB_MySQL) GetDeleteMetaSql(dbName string) string {
	return "DELETE FROM MetaTable where DBName = '" + dbName + "'"
}

func (mysql *DB_MySQL) GetDefineListSql(dbName string) string {
	return "CREATE TABLE `" + dbName + "` (id INTEGER  PRIMARY KEY AUTO_INCREMENT, category text, field01 text, field02 text, note text)"
}

func (mysql *DB_MySQL) GetInsertSql(dbName string) string {
	return "INSERT INTO `" + dbName + "` (category, field01, field02, note) VALUES (?, ?, ?, ?)"
}

func (mysql *DB_MySQL) GetUpdateSql(dbName string) string {
	return "UPDATE `" + dbName + "` set category=?, field01=?, field02=?, note=? where id="
}

func (mysql *DB_MySQL) GetDeleteSql(dbName string) string {
	return "DELETE FROM `" + dbName + "` where id="
}

func (mysql *DB_MySQL) GetLimitSql(countRec int, fromRec int) string {
	return " limit " + strconv.Itoa(countRec) + " offset " + strconv.Itoa(fromRec-1)
}

func (mysql *DB_MySQL) GetDropSql(dbName string) string {
	return "drop table if exists `" + dbName + "`"

}

func (mysql *DB_MySQL) GetQname(dbName string) string {
	return "`" + dbName + "`"
}

func (mysql *DB_MySQL) GetLastRowId() string {
	//return .cursor.lastrowid
	return ""
}

func (mysql *DB_MySQL) GetRecordCountSql(dbName string) string {
	return "select count(*) from `" + dbName + "`"
}

func (mysql *DB_MySQL) Connect(connectString string) error {
	var err error
	mysql.Db, err = sql.Open("mysql", connectString)
	return err
}

func (mysql *DB_MySQL) Close() {
	mysql.Db.Close()
}
