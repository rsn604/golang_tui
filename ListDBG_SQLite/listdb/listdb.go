package listdb

import (
	"encoding/json"
)

type ListDB struct {
	ID            int        `json:"id"`
	DbName        string     `json:"dbName"`
	FieldName01   string     `json:"fieldName01"`
	FieldName02   string     `json:"fieldName02"`
	CategoryList  string     `json:"categoryList"`
	CurrentNumber int        `json:"currentNumber"`
	ListData      []ListItem `json:"listData"`
}

type ListItem struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Field01  string `json:"field01"`
	Field02  string `json:"field02"`
	Note     string `json:"note"`
}

func Json2Obj(jsonStr string) (*ListDB, error) {
	jsonBytes := ([]byte)(jsonStr)
	listDB := new(ListDB)
	if err := json.Unmarshal(jsonBytes, listDB); err != nil {
		return nil, err
	}
	return listDB, nil
}

func (listDB *ListDB) Obj2Json() ([]byte, error) {
	jsonBytes, err := json.Marshal(listDB)
	return jsonBytes, err
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
