package listdb

// ----------------------------------------------------------------------
type Manager interface {
	Connect(databaseName string, connectString string) error
	GetDbNames() ([]string, error)
	GetCategoryList(dbName string) ([]string, error)
	SearchDB(dbName string, category string, search string, from_rec int, count_rec int) *ListDB
	GetRecordCount(dbName string, category string, search string) int
	Insert(dbName string, listItem *ListItem) (int, error)
	Update(dbName string, id int, listItem *ListItem) (*ListDB, error)
	Delete(dbName string, id int) error
	ImportCSV(fname string) bool
	Define() error
	Close()
}

func GetManager(name string) Manager {
	if name == "BOLT" {
		return new(BoltManager)
	} else {
		//return new(ListDBManager)
	}

	return nil
}
