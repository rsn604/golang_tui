package listdb

import (
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	//"bytes"
	"io"
	"strings"
)

// -------------------------------------------------
type BoltManager struct {
	Db *bolt.DB
}

const (
	LISTDB    = "ListDB"
	METATABLE = "MetaTable"
)

// -------------------------------------------------
func (self *BoltManager) GetDb() *bolt.DB {
	return self.Db
}

func (self *BoltManager) Connect(databaseName string, connectString string) error {
	var err error
	self.Db, err = bolt.Open(connectString, 0600, nil)
	if err != nil {
		return err
	}
	return nil
}

func (self *BoltManager) Close() {
	self.GetDb().Close()
}

func (self *BoltManager) Define() error {
	err := self.GetDb().Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(LISTDB))
		root, err := tx.CreateBucketIfNotExists([]byte(LISTDB))
		if err != nil {
			return fmt.Errorf("D.ER could not create root bucket: %v", err)
		}

		root.DeleteBucket([]byte(METATABLE))
		_, err = root.CreateBucketIfNotExists([]byte(METATABLE))
		if err != nil {
			return fmt.Errorf("D.ER could not create weight bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("D.ER could not set up buckets, %v", err)
	}
	return nil
}

// ---------------------------------------
func (self *BoltManager) setMetaTable(tx *bolt.Tx, fields []string) error {
	point := len(fields)
	for i, c := range fields {
		if c == "" {
			point = i
			break
		}
	}

	err := tx.Bucket([]byte(LISTDB)).Bucket([]byte(METATABLE)).Put([]byte(fields[0]), []byte(strings.Join(fields[1:point], ",")))
	if err != nil {
		return fmt.Errorf("D.ER could not set config: %v", err)
	}
	return nil
}

func (self *BoltManager) defineList(tx *bolt.Tx, dbName string) error {
	root := tx.Bucket([]byte(LISTDB))
	_, err := root.CreateBucketIfNotExists([]byte(dbName))
	if err != nil {
		return fmt.Errorf("D.ER could not create root bucket: %v", err)
	}
	return nil
}

func (self *BoltManager) itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func (self *BoltManager) addList(tx *bolt.Tx, dbName string, fields []string) error {
	listItem := new(ListItem)
	b := tx.Bucket([]byte(LISTDB)).Bucket([]byte(dbName))
	id, _ := b.NextSequence()
	listItem.ID = int(id)
	listItem.Category = fields[0]
	listItem.Field01 = fields[1]
	listItem.Field02 = fields[2]
	listItem.Note = strings.Join(fields[3:], "\n")

	buf, err := json.Marshal(listItem)
	if err != nil {
		return err
	}
	b.Put(self.itob(listItem.ID), buf)
	return nil
}

func (self *BoltManager) ImportCSV(fname string) bool {
	var fp *os.File
	var err error

	fp, err = os.Open(fname)
	if err != nil {
		return false
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	reader.Comma = ','
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1 // Nocheck fields count
	var firstTime = true
	var dbName string

	tx, err := self.GetDb().Begin(true)
	if err != nil {
		return false
	}
	defer tx.Rollback()

	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return false
		}
		if len(fields) == 0 {
			continue
		}
		if firstTime == true {
			err = self.setMetaTable(tx, fields)
			if err != nil {
				return false
			}
			dbName = fields[0]
			err = self.defineList(tx, dbName)
			if err != nil {
				return false
			}
			firstTime = false
		} else {
			err = self.addList(tx, dbName, fields)
			if err != nil {
				return false
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return false
	}
	return true
}
