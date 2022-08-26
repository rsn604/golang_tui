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

func (self *BoltManager) GetDbNames() ([]string, error) {
	var dbNames []string
	err := self.GetDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LISTDB)).Bucket([]byte(METATABLE))
		err := b.ForEach(func(k, v []byte) error {
			dbNames = append(dbNames, string(k))
			return nil
		})
		return err
	})
	return dbNames, err
}

func (self *BoltManager) GetCategoryList(dbName string) ([]string, error) {
	var categoryList []string
	err := self.GetDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LISTDB)).Bucket([]byte(METATABLE))
		v := b.Get([]byte([]byte(dbName)))
		s := strings.Split(string(v), ",")
		categoryList = s[2:]
		return nil
	})
	return categoryList, err
}

// ------------------------------------------------------
func (self *BoltManager) contains(field string, search string) bool {
	if search == "" {
		return true
	}
	return strings.Contains(strings.ToLower(field), strings.ToLower(search))
}

func (self *BoltManager) SearchDB(dbName string, category string, search string, from_rec int, count_rec int) *ListDB {
	var listdb *ListDB
	_ = self.GetDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LISTDB)).Bucket([]byte(METATABLE))
		v := b.Get([]byte([]byte(dbName)))
		s := strings.Split(string(v), ",")
		listdb = new(ListDB)
		listdb.DbName = dbName
		listdb.FieldName01 = s[0]
		listdb.FieldName02 = s[1]
		listdb.CategoryList = strings.Join(s[2:], ",")

		b = tx.Bucket([]byte(LISTDB)).Bucket([]byte(dbName))
		c := b.Cursor()
		var listItem ListItem
		recCount := 1
		targetCount := 1
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if count_rec > 0 && targetCount > count_rec {
				break
			}
			if err := json.Unmarshal(v, &listItem); err != nil {
				return err
			}

			// check category
			if category != "" && category != listItem.Category {
				continue
			}

			// check string
			if !self.contains(listItem.Field01, search) && !self.contains(listItem.Field02, search) && !self.contains(listItem.Note, search) {
				continue
			}

			// add list within range
			if recCount >= from_rec {
				listdb.AddListData(listItem)
				targetCount++
			}
			// count up
			recCount++
		}
		return nil
	})
	return listdb
}

func (self *BoltManager) GetRecordCount(dbName string, category string, search string) int {
	dataCount := 0
	_ = self.GetDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LISTDB)).Bucket([]byte(dbName))
		c := b.Cursor()
		var listItem ListItem

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := json.Unmarshal(v, &listItem); err != nil {
				return err
			}

			// check category
			if category != "" && category != listItem.Category {
				continue
			}

			// check string
			if !self.contains(listItem.Field01, search) && !self.contains(listItem.Field02, search) && !self.contains(listItem.Note, search) {
				continue
			}

			// count up
			dataCount++
		}
		return nil
	})
	return dataCount
}

// ------------------------------------------------------
func (self *BoltManager) Delete(dbName string, id int) error {
	err := self.GetDb().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LISTDB)).Bucket([]byte(dbName))
		b.Delete(self.itob(id))
		return nil
	})
	return err
}

func (self *BoltManager) Update(dbName string, id int, listItem *ListItem) (*ListDB, error) {
	err := self.GetDb().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LISTDB)).Bucket([]byte(dbName))
		buf, err := json.Marshal(listItem)
		if err != nil {
			return err
		}
		b.Put(self.itob(id), buf)
		return nil
	})
	return nil, err
}

func (self *BoltManager) Insert(dbName string, listItem *ListItem) (int, error) {
	var id int
	err := self.GetDb().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LISTDB)).Bucket([]byte(dbName))
		seq, _ := b.NextSequence()
		id = int(seq)

		buf, err := json.Marshal(listItem)
		if err != nil {
			return err
		}
		b.Put(self.itob(id), buf)
		return nil
	})
	return id, err
}
