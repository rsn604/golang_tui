package tui

type Common struct {
	databaseName  string
	connectString string
	tableName     string
	category      string
	search        string
	from          int
	selectedItem  int
	cols          int
	rows          int
}

func NewCommon() *Common {
	common := &Common{
		from:         1,
		selectedItem: 1,
	}
	return common
}

func (self *Common) resetPaging() {
	self.from = 1
	self.selectedItem = 1
}

func (self *Common) reset() {
	self.search = ""
	self.category = ""
	self.resetPaging()
}
