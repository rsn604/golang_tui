package listdb

type ListDB struct {
	ID            int
	DbName        string
	FieldName01   string
	FieldName02   string
	CategoryList  string
	CurrentNumber int
	ListData      []ListItem
}

type ListItem struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Field01  string `json:"field01"`
	Field02  string `json:"field02"`
	Note     string `json:"note"`
}

func (listdb *ListDB) AddListData(listItem ListItem) {
	listdb.ListData = append(listdb.ListData, listItem)
}

func (listdb *ListDB) GetListData() []ListItem {
	return listdb.ListData
}

func (listdb *ListDB) CountListData() int {
	return len(listdb.ListData)
}
