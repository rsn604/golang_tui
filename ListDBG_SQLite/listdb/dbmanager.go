package listdb

import (
	"errors"
)

type ListDBManager struct {
	Db Database
}

func (listDBManager *ListDBManager) Connect(databaseName string, connectString string) error {
	var err error
	//var db listdb.Database
	var db Database

	if databaseName == "SQLITE3" {
		db = new(DB_SQLite3)
		err = db.Connect(connectString)

	} else if databaseName == "MYSQL" {
		db = new(DB_MySQL)
		err = db.Connect(connectString)
	} else {
		return errors.New("D.ER Database name error")
	}
	if err != nil {
		return err
	}
	listDBManager.Db = db
	return nil
}

func (listDBManager *ListDBManager) GetDatabase() Database {
	return listDBManager.Db
}

func (listDBManager *ListDBManager) GetDbNames() ([]string, error) {
	return GetDbNames(listDBManager.GetDatabase())
}

func (listDBManager *ListDBManager) GetCategoryList(dbName string) ([]string, error) {
	return GetCategoryList(listDBManager.GetDatabase(), dbName)
}

func (listDBManager *ListDBManager) Close() {
	listDBManager.GetDatabase().Close()
}

func (listDBManager *ListDBManager) ImportCSV(fname string) bool {
	return ImportCSV(listDBManager.GetDatabase(), fname)
}

func (listDBManager *ListDBManager) Define() error {
	return Define(listDBManager.GetDatabase())
}

func (listDBManager *ListDBManager) SearchDB(dbName string, category string, search string, from_rec int, count_rec int) *ListDB {
	return SearchDB(listDBManager.GetDatabase(), dbName, category, search, from_rec, count_rec)
}

func (listDBManager *ListDBManager) GetDB(dbName string, sql_value string, from_rec int, count_rec int) *ListDB {
	return GetDB(listDBManager.GetDatabase(), dbName, sql_value, from_rec, count_rec)
}

func (listDBManager *ListDBManager) GetRecordCount2(dbName string, sql_value string) int {
	return GetRecordCount(listDBManager.GetDatabase(), dbName, sql_value)
}

func (listDBManager *ListDBManager) GetRecordCount(dbName string, category string, search string) int {
	return GetRecordCount(listDBManager.GetDatabase(), dbName, CreateSQL(category, search))
}

// ----------------------------------------------------------------------
func (listDBManager *ListDBManager) Delete(dbName string, id int) error {
	_, err := DeleteData(listDBManager.GetDatabase(), dbName, id)
	return err
}

func (listDBManager *ListDBManager) Update(dbName string, id int, listItem *ListItem) (*ListDB, error) {
	var fields []string
	fields = append(fields, listItem.Category)
	fields = append(fields, listItem.Field01)
	fields = append(fields, listItem.Field02)
	fields = append(fields, listItem.Note)
	_, err := UpdateData(listDBManager.GetDatabase(), dbName, id, fields)
	listDB := GetDataById(listDBManager.GetDatabase(), dbName, id)
	return listDB, err
}

func (listDBManager *ListDBManager) Insert(dbName string, listItem *ListItem) (int, error) {
	var fields []string
	fields = append(fields, listItem.Category)
	fields = append(fields, listItem.Field01)
	fields = append(fields, listItem.Field02)
	fields = append(fields, listItem.Note)
	return InsertData(listDBManager.GetDatabase(), dbName, fields)
}
