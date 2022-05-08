package tui

import (
	"github.com/rivo/tview"
)

// ------------------------------------------------------
type MyApplication struct {
	app *tview.Application
}

func (self *MyApplication) exit() {
	self.app.Stop()
}

func (self *MyApplication) display(primitive tview.Primitive) {
	self.app.SetRoot(primitive, true)
}

func (self *MyApplication) run(app *tview.Application, primitive tview.Primitive) {
	if app == nil {
		self.app = tview.NewApplication()
		if err := self.app.SetRoot(primitive, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}
	} else {
		self.app = app
		self.display(primitive)
	}
}

// ------------------------------------------------------
type MainList struct {
	MyApplication
}

var mainList *MainList = nil

func NewMainList() *MainList {
	if mainList == nil {
		abstract := MyApplication{}
		mainList = &MainList{MyApplication: abstract}
	}
	return mainList
}

func (self *MainList) display(common *Common) {
	self.MyApplication.display(self.doformat(common))
}

func (self *MainList) run(app *tview.Application, common *Common) {
	self.MyApplication.run(app, self.doformat(common))
}

func (self *MainList) Init(databaseName string, connectString string, cols int, rows int) {
	common := NewCommon()
	common.reset()
	common.databaseName = databaseName
	common.connectString = connectString
	common.cols = cols
	common.rows = rows
	self.run(nil, common)
}
