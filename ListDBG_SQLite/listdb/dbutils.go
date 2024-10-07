package listdb

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	_ "github.com/glebarez/go-sqlite"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

// ----------------------------------------------------------------------
func GetDbNames(db Database) ([]string, error) {
	rows, err := db.GetDb().Query("select DBName from MetaTable")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dbNames := []string{}
	for rows.Next() {
		dbName := ""
		err = rows.Scan(&dbName)
		if err != nil {
			return nil, err
		}
		dbNames = append(dbNames, dbName)
	}
	return dbNames, nil
}

func GetCategoryList(db Database, dbName string) ([]string, error) {
	var category string
	var categoryList []string
	if err := db.GetDb().QueryRow("select categoryList from MetaTable where DBName =? limit 1", dbName).Scan(&category); err != nil {
		return nil, err
	}
	categoryList = strings.Split(category, ",")
	return categoryList, nil
}

// ----------------------------------------------------------------------
func GetData(db Database, listDB *ListDB, sql_value string, from_rec int, count_rec int) *ListDB {
	qs_sql := "select * from " + db.GetQname(listDB.DbName) + " "

	if sql_value != "" {
		qs_sql = qs_sql + sql_value
	}
	if count_rec > 0 {
		qs_sql = qs_sql + db.GetLimitSql(count_rec, from_rec)
	}

	rows, err := db.GetDb().Query(qs_sql + " ")

	if err != nil {
		fmt.Println("D.ER GetData Error " + qs_sql)
		//return nil, err
		//return nil
	}
	defer rows.Close()

	var rec = 1
	var listItem ListItem
	for rows.Next() {
		if (count_rec > 0) && (rec > count_rec) {
			break
		}
		if err := rows.Scan(&listItem.ID, &listItem.Category, &listItem.Field01, &listItem.Field02, &listItem.Note); err != nil {
			fmt.Println("D.ER GetData Error")
			//return nil
		}

		listDB.AddListData(listItem)
		rec = rec + 1
	}

	return listDB
}

// func GetMetaTable(db Database, sql string)(ListDB, error){
func GetMetaTable(db Database, sql string) *ListDB {
	//fmt.Println("D.MG SQL:"+sql)
	listDB := new(ListDB)
	if err := db.GetDb().QueryRow(sql).Scan(&listDB.ID, &listDB.DbName, &listDB.FieldName01, &listDB.FieldName02, &listDB.CategoryList); err != nil {
		fmt.Println("D.ER GetMetaData Error " + sql)
		//return nil, err
	}
	//return listdb, nil
	//fmt.Println("ListDB.DbName:"+listDB.DbName)
	return listDB
}

func GetDataById(db Database, dbName string, id int) *ListDB {
	sql := "select * from MetaTable where DBName ='" + dbName + "'"

	listDB := GetMetaTable(db, sql)
	sql_Value := "where id = " + strconv.Itoa(id)
	return GetData(db, listDB, sql_Value, 1, 1)
}

func GetDB(db Database, dbName string, sql_value string, from_rec int, count_rec int) *ListDB {
	var sql = "select * from MetaTable where DBName ='" + dbName + "'"
	//fmt.Println("D.MG SQL:"+sql)
	listDB := GetMetaTable(db, sql)
	return GetData(db, listDB, sql_value, from_rec, count_rec)
}

func CreateSQL(category string, search string) string {
	sql := ""
	if search != "" {
		sql = " where (field01 like " + "\"%" + search + "%\""
		sql += " or field02 like " + "\"%" + search + "%\""
		sql += " or note like " + "\"%" + search + "%\")"
	}
	if category != "" {
		if sql == "" {
			sql = "where (category = '" + category + "')"
		} else {
			sql = sql + " and (category = '" + category + "')"
		}
	}
	return sql
}

func SearchDB(db Database, dbName string, category string, search string, from_rec int, count_rec int) *ListDB {
	return GetDB(db, dbName, CreateSQL(category, search), from_rec, count_rec)
}

// ----------------------------------------------------------------------
func DeleteData(db Database, dbName string, id int) (sql.Result, error) {
	result, err := db.GetDb().Exec(db.GetDeleteSql(dbName) + strconv.Itoa(id))
	return result, err
}

func UpdateData(db Database, dbName string, id int, fields []string) (sql.Result, error) {
	result, err := db.GetDb().Exec(db.GetUpdateSql(dbName)+strconv.Itoa(id), fields[0], fields[1], fields[2], strings.Join(fields[3:], "\n"))
	return result, err
}

func InsertData(db Database, dbName string, fields []string) (int, error) {
	//fmt.Println("InsertData --------------------------------")
	//fmt.Println(fields)
	result, err := db.GetDb().Exec(db.GetInsertSql(dbName), fields[0], fields[1], fields[2], strings.Join(fields[3:], "\n"))
	var lastRowId int64
	if err == nil {
		lastRowId, err = result.LastInsertId()
	}
	return int(lastRowId), err
}

// ----------------------------------------------------------------------
func InsertMeta(db Database, fields []string) bool {
	if len(fields) < 4 {
		return false
	}
	_, err := db.GetDb().Exec(db.GetDeleteMetaSql(fields[0]))
	_, err = db.GetDb().Exec(db.GetInsertMetaSql(), fields[0], fields[1], fields[2], strings.Join(fields[3:], ","))
	if err != nil {
		return false
	}
	return true
}

func DeleteList(db Database, dbName string) {
	//_, err := db.Exec(" drop table if exists ?", dbName)
	//db.Exec(" drop table if exists ?", dbName)
	db.GetDb().Exec(db.GetDropSql(dbName))
}

func DefineList(db Database, dbName string) error {
	DeleteList(db, dbName)
	_, err := db.GetDb().Exec(db.GetDefineListSql(dbName))
	return err
}

func ImportCSV(db Database, fname string) bool {
	var fp *os.File
	var err error

	fp, err = os.Open(fname)
	if err != nil {
		return false
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	reader.Comma = ','
	reader.LazyQuotes = true    // ダブルクオートを厳密にチェックしない！
	reader.FieldsPerRecord = -1 // Nocheck fields count
	var firstTime = true
	var dbName string
	var retCode = true

	tx, _ := db.GetDb().Begin()
	var stmt *sql.Stmt

	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return false
		}
		if len(fields) == 0 {
			fmt.Println(fields)
			continue
		}
		if firstTime == true {
			retCode = InsertMeta(db, fields)
			if retCode == false {
				fmt.Println("D.ER InsertMeta Error")
				return retCode
			}
			dbName = fields[0]
			err = DefineList(db, dbName)
			if err != nil {
				fmt.Println("D.ER DefineList Error")
				fmt.Println(err)
				return false
			}
			firstTime = false
			stmt, err = tx.Prepare(db.GetInsertSql(dbName))
			defer tx.Rollback()
		} else {

			//_, err = InsertData(db, dbName, fields)
			_, err = stmt.Exec(fields[0], fields[1], fields[2], strings.Join(fields[3:], "\n"))
			if err != nil {
				fmt.Println("D.ER Prepared and Insert Error")
				fmt.Println(err)
				return false
			}
		}
	}
	tx.Commit()
	return retCode
}

func Define(db Database) error {
	_, err := db.GetDb().Exec(db.GetDropMetaSql())
	_, err = db.GetDb().Exec(db.GetDefineSql())
	return err
}

func GetRecordCount(db Database, dbName string, sql_value string) int {
	var sql string
	sql = db.GetRecordCountSql(dbName)

	if sql_value != "" {
		sql = sql + sql_value
	}
	sql = sql + " limit 1"
	var row int
	if err := db.GetDb().QueryRow(sql).Scan(&row); err != nil {
		return 0
	}
	return row
}
