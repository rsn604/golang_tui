package main

import (
	"listdbg/listdb"
	"fmt"
)

func printData(listDB *listdb.ListDB){
	fmt.Println("------------------------------------------------------------------")
    fmt.Printf("ID=%d, DbName=%s, FieldName01=%s, FieldName02=%s, CategoryList=%s\n", listDB.ID, listDB.DbName, listDB.FieldName01, listDB.FieldName02, listDB.CategoryList)
	for _, listItem := range listDB.ListData {	
		fmt.Printf("%d | %s | %s | %s | %s\n", listItem.ID, listItem.Category, listItem.Field01, listItem.Field02, listItem.Note)
	}
}

func printCount(count int){
	fmt.Printf("Count=%d\n", count)
}

func testdb(manager listdb.Manager, databaseName string, connectString string) {
	var err error 

	err = manager.Connect(databaseName, connectString) 
    if err != nil { panic(err) }
	var dbNames []string
	dbNames, err = manager.GetDbNames()
    for _, dbName := range dbNames {
		categoryList, _ := manager.GetCategoryList(dbName)
		fmt.Printf("table:%s categoryList:%v\n", dbName, categoryList)
	}

	listdb := manager.SearchDB(dbNames[0], "", "", 1, 3)
	printData(listdb)
	count := manager.GetRecordCount(dbNames[0], "", "")
	printCount(count)

	listdb = manager.SearchDB(dbNames[0], "", "ディクスン・カー", 1, 3)
	printData(listdb)
	count = manager.GetRecordCount(dbNames[0], "", "ディクスン・カー")
	printCount(count)

	listdb = manager.SearchDB(dbNames[0], "1958", "", 1, 3)
	printData(listdb)
	count = manager.GetRecordCount(dbNames[0], "1958", "")
	printCount(count)

	listdb = manager.SearchDB(dbNames[0], "1958", "ディクスン・カー", 1, 3)
	printData(listdb)
	count = manager.GetRecordCount(dbNames[0], "1958", "ディクスン・カー")
	printCount(count)

	manager.Close()
}

func main() {
	var manager = listdb.GetManager("BOLT")
	testdb(manager, "BOLT", "./db/ListDB.boltdb")
}
