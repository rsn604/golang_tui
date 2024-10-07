package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
	"listdbg/listdb"
	"sort"
	"strings"
)

var isLast bool

// -------------------------------------------------
//
//	Table
//
// -------------------------------------------------
func (self *MainList) getTable(pages *tview.Pages, common *Common) {
	manager := listdb.GetManager(common.databaseName)
	err := manager.Connect(common.databaseName, common.connectString)
	if err != nil {
		panic(err)
	}
	dbNames, _ := manager.GetDbNames()
	sort.Strings(dbNames)
	manager.Close()

	tables := MySelectBox(dbNames, 30, 20, 2, 3, 0).
		SetDoneFunc(func(buttonLabel string, inputString string) {
			if buttonLabel == "OK" {
				common.reset()
				common.tableName = inputString
				pages.RemovePage("table")
				self.display(common)
			}
		})

	tables.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			pages.RemovePage("table")
		}
		return event
	})
	pages.AddPage("table", tables, false, true)
}

// -------------------------------------------------
//
//	Category
//
// -------------------------------------------------
func (self *MainList) getCategory(pages *tview.Pages, common *Common) {
	manager := listdb.GetManager(common.databaseName)
	err := manager.Connect(common.databaseName, common.connectString)
	if err != nil {
		panic(err)
	}
	categoryList, _ := manager.GetCategoryList(common.tableName)
	manager.Close()
	categoryList = append([]string{"Clear"}, categoryList...)
	current := 0
	for i, c := range categoryList {
		if c == common.category {
			current = i
			break
		}
	}

	category := MySelectBox(categoryList, 30, 10, 2, 3, current).
		SetDoneFunc(func(buttonLabel string, inputString string) {
			if buttonLabel == "OK" {
				if inputString == "Clear" {
					common.category = ""
				} else {
					common.category = inputString
				}
				pages.RemovePage("category")
				common.resetPaging()
				self.display(common)
			}
		})

	category.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			pages.RemovePage("category")
		}
		return event
	})
	pages.AddPage("category", category, false, true)
}

// -------------------------------------------------
//
//	Search
//
// -------------------------------------------------
func (self *MainList) getSearch(pages *tview.Pages, common *Common) {
	search := MyInputDialog(self.app, common.search, 40, 8, 2, 4).
		SetDoneFunc(func(buttonLabel string, inputString string) {
			if buttonLabel == "OK" {
				common.search = inputString
			} else if buttonLabel == "Cancel" {
				common.search = ""
			}
			pages.RemovePage("search")
			common.resetPaging()
			self.display(common)
		})
	pages.AddPage("search", search, false, true)
}

// -------------------------------------------------
//
//	Paging
//
// -------------------------------------------------
func (self *MainList) firstPage(common *Common) bool {
	return common.from == 1
}

func (self *MainList) lastPage() bool {
	return isLast
}

func (self *MainList) nextPage(common *Common) {
	if !self.lastPage() {
		common.from += (common.rows - 2)
		self.display(common)
	}
}

func (self *MainList) priorPage(common *Common) {
	if !self.firstPage(common) {
		common.from -= (common.rows - 2)
	}
	if common.from < 1 {
		common.from = 1
	}

	self.display(common)
}

// -------------------------------------------------
//
//	List
//
// -------------------------------------------------
func (self *MainList) getFieldsInfo(listdata []listdb.ListItem, common *Common) (int, []int) {
	maxLength := 0
	var fieldsLength []int
	var flen int
	for _, s := range listdata {
		flen = runewidth.StringWidth(strings.TrimSpace(s.Field01))
		if flen > maxLength {
			maxLength = flen
		}
		fieldsLength = append(fieldsLength, flen)
	}
	if maxLength > common.cols {
		maxLength = common.cols
	}
	return maxLength, fieldsLength
}

func (self *MainList) setList(list *tview.List, common *Common) (*tview.List, int) {
	manager := listdb.GetManager(common.databaseName)
	err := manager.Connect(common.databaseName, common.connectString)
	if err != nil {
		panic(err)
	}
	listdb := manager.SearchDB(common.tableName, common.category, common.search, common.from, common.rows-1)
	recordCount := manager.GetRecordCount(common.tableName, common.category, common.search)
	manager.Close()
	listdata := listdb.GetListData()
	if len(listdata) <= common.rows-2 {
		isLast = true
	} else {
		isLast = false
	}

	maxLength, fieldsLength := self.getFieldsInfo(listdata, common)
	list.Clear()

	for i, s := range listdata {
		if i >= common.rows-2 {
			break
		}
		if maxLength < fieldsLength[i] {
			list.AddItem(strings.TrimSpace(s.Field01)+s.Field02, "", 0, nil)
		} else {
			list.AddItem(strings.TrimSpace(s.Field01)+strings.Repeat(" ", (maxLength-fieldsLength[i]+1))+s.Field02, "", 0, nil)
		}
	}
	return list, recordCount
}

// -------------------------------------------------
//
//	format screen
//
// -------------------------------------------------
func isDataExist(common *Common) bool {
	return common.tableName != ""
}

func (self *MainList) getStartRecord(common *Common) {
	if common.from >= common.selectedItem {
		return
	}
	r := common.rows - 2
	p := common.selectedItem / r
	if p*r == common.selectedItem {
		common.from = (p-1)*r + 1
	} else {
		common.from = p*r + 1
	}
}

func (self *MainList) doformat(common *Common) tview.Primitive {
	self.getStartRecord(common)
	pages := tview.NewPages()
	header := tview.NewFlex()

	btnT := myButton("<T>").SetSelectedFunc(func() {
		self.getTable(pages, common)
	})
	header.AddItem(btnT, 0, 1, true)

	var recordCount int
	list := tview.NewList().ShowSecondaryText(false).SetSelectedTextColor(tcell.ColorWhite).SetSelectedBackgroundColor(tcell.ColorAqua).SetSelectedFocusOnly(true)
	footer := tview.NewTextView()

	if isDataExist(common) {
		list, recordCount = self.setList(list, common)
		if common.selectedItem == 0 {
			list.SetCurrentItem(0)
		} else {
			list.SetCurrentItem(common.selectedItem - common.from)
			common.selectedItem = 0
		}

		if !self.lastPage() {
			btnN := myButton("<N>").SetSelectedFunc(func() {
				self.nextPage(common)
			})
			header.AddItem(btnN, 0, 1, true)
		} else {
			header.AddItem(nil, 0, 1, false)
		}

		if !self.firstPage(common) {
			btnP := myButton("<P>").SetSelectedFunc(func() {
				self.priorPage(common)
			})
			header.AddItem(btnP, 0, 1, true)
		} else {
			header.AddItem(nil, 0, 1, false)
		}

		btnC := myButton("<C>").SetSelectedFunc(func() {
			self.getCategory(pages, common)
		})
		header.AddItem(btnC, 0, 1, true)

		btnS := myButton("<S>").SetSelectedFunc(func() {
			self.getSearch(pages, common)
		})
		header.AddItem(btnS, 0, 1, true)
		//footer.SetTextColor(tcell.ColorYellow).SetText(fmt.Sprintf("%d(%d)/%d %d", common.from, common.selectedItem, common.rows, recordCount)+common.tableName)
		footer.SetTextColor(tcell.ColorYellow).SetText(fmt.Sprintf("%d/%d ", common.from, recordCount) + common.tableName)
		//.SetBackgroundColor(tcell.ColorBlue)
	} else {
		header.AddItem(nil, 0, 1, false)
		header.AddItem(nil, 0, 1, false)
	}

	btnQ := myButton("<Q>").SetSelectedFunc(func() {
		self.exit()
		//self.app.Stop()
	})
	header.AddItem(btnQ, 0, 1, true)

	body := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(list, 0, 1, true)

	main := tview.NewFlex().SetDirection(tview.FlexRow)
	if isDataExist(common) {
		main.AddItem(header, 1, 0, false)
	} else {
		main.AddItem(header, 1, 0, true)
	}
	main.AddItem(body, 0, 1, true).
		AddItem(footer, 1, 0, false)

	// ------------------------------
	// InputCapture on Header
	// ------------------------------
	header.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight, tcell.KeyTab:
			myFlexFocus(self.app, header, false)
			return nil
		case tcell.KeyLeft, tcell.KeyBacktab:
			myFlexFocus(self.app, header, true)
			return nil
		case tcell.KeyDown:
			//list.SetCurrentItem(0)
			self.app.SetFocus(list)
			return nil
		}
		return event
	})

	// ------------------------------
	// InputCapture on Body
	// ------------------------------
	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			if isDataExist(common) {
				if list.GetCurrentItem() == 0 {
					myFlexFocus(self.app, header, true)
					return nil
				}
			}
		}
		return event
	})

	// ------------------------------
	// InputCapture on Main(common)
	// ------------------------------
	main.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				self.exit()
				//self.app.Stop()
			case 't', 'T':
				self.getTable(pages, common)
			case 'n', 'N':
				if isDataExist(common) {
					self.nextPage(common)
				}
			case 'p', 'P':
				if isDataExist(common) {
					self.priorPage(common)
				}
			case 'c', 'C':
				if isDataExist(common) {
					self.getCategory(pages, common)
				}
			case 's', 'S':
				if isDataExist(common) {
					self.getSearch(pages, common)
				}
			}
		}
		return event
	})

	// ------------------------------
	// Go to Detail with Selected
	// ------------------------------
	list.SetSelectedFunc(func(index int, s string, secondary string, code rune) {
		common.selectedItem = index + common.from
		NewDetail().run(self.app, common)
		return
	})

	pages.AddPage("main", main, true, true)
	return pages
}
