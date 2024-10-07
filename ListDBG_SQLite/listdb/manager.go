package listdb

// ----------------------------------------------------------------------
type Manager interface {
	Connect(databaseName string, connectString string) error
	GetDbNames() ([]string, error)
	GetCategoryList(dbName string) ([]string, error)
	SearchDB(dbName string, category string, search string, from_rec int, count_rec int) *ListDB
	//GetDB(dbName string,  sql_value string,  from_rec int ,  count_rec int)(ListDB)

	GetRecordCount2(dbName string, sql_value string) int
	GetRecordCount(dbName string, category string, search string) int
	Delete(dbName string, id int) error
	Update(dbName string, id int, listItem *ListItem) (*ListDB, error)
	Insert(dbName string, listItem *ListItem) (int, error)

	ImportCSV(fname string) bool
	Define() error
	Close()
}

func GetManager(name string) Manager {
	if name == "http" {
		return nil
	} else {
		return new(ListDBManager)
	}
}
