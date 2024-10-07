package main

import (
	"fmt"
	"io/ioutil"
	"listdbg/listdb"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetFilesFromDir(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, GetFilesFromDir(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}
	return paths
}

func loadDB(manager listdb.Manager, databaseName string, connectString string, csvdir string) {
	var retCode bool
	var err error
	start := time.Now()

	err = manager.Connect(databaseName, connectString)
	if err != nil {
		panic(err)
	}

	if err = manager.Define(); err != nil {
		panic(err)
	}

	csvfiles := GetFilesFromDir(csvdir)
	for _, csvfile := range csvfiles {
		pos := strings.LastIndex(csvfile, ".")
		if csvfile[pos:] == ".csv" {
			retCode = manager.ImportCSV(csvfile)
			fmt.Printf("FileName:%s, RetCode:%t\n", csvfile, retCode)
		}
	}
	manager.Close()
	end := time.Now()
	fmt.Printf("%fsec\n", (end.Sub(start)).Seconds())
}

func main() {
	var databaseName string
	var connectString string
	var csvdir string
	if len(os.Args) == 4 {
		databaseName = os.Args[1]
		connectString = os.Args[2]
		csvdir = os.Args[3]
	} else {
		databaseName = "SQLITE3"
		connectString = "./db/ListDB.sqlite3"
		csvdir = "./csv"
	}
	var manager = listdb.GetManager(databaseName)
	loadDB(manager, databaseName, connectString, csvdir)
}
