package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
	"github.com/toqueteos/webbrowser"
	"listdbg/listdb"
	"listdbg/mytview"
	"os/exec"
	"strings"
)

var isLastItem bool

// -------------------------------------------------
//
//	Paging
//
// -------------------------------------------------
func (self *Detail) firstPage(common *Common) bool {
	return common.selectedItem == 1
}

func (self *Detail) lastPage() bool {
	return isLastItem
}

func (self *Detail) nextPage(common *Common) {
	if !self.lastPage() {
		common.selectedItem += 1
		self.display(common)
	}
}

func (self *Detail) priorPage(common *Common) {
	if !self.firstPage(common) {
		common.selectedItem -= 1
	}
	if common.selectedItem < 1 {
		common.selectedItem = 1
	}
	self.display(common)
}

// -------------------------------------------------
//
//	Category
//
// -------------------------------------------------
func (self *Detail) getCategoryList(pages *tview.Pages, common *Common, currentCategory string, btnCategory *tview.Button) {
	manager := listdb.GetManager(common.databaseName)
	err := manager.Connect(common.databaseName, common.connectString)
	if err != nil {
		panic(err)
	}
	categoryList, _ := manager.GetCategoryList(common.tableName)
	manager.Close()
	current := 0
	for i, c := range categoryList {
		if c == currentCategory {
			current = i
			break
		}
	}

	category := MySelectBox(categoryList, 30, 10, 2, 3, current, true).
		SetDoneFunc(func(buttonLabel string, inputString string) {
			if buttonLabel == "OK" {
				btnCategory.SetLabel(self.setButtonLabel(inputString, common))
			}
			pages.RemovePage("category")
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
//	Set button label
//
// -------------------------------------------------
func (self *Detail) setButtonLabel(label string, common *Common) string {
	flen := runewidth.StringWidth(label)
	s := label
	if common.cols > flen {
		s = s + strings.Repeat(" ", common.cols-flen)
	}
	return s
}

// -------------------------------------------------
//
//	Set field to invoke Browser
//
// -------------------------------------------------
func (self *Detail) getIntentPath(command string) string {
	path, err := exec.LookPath(command)
	if err != nil {
		return ""
	}
	return path
}

func (self *Detail) setFieldIntent(fieldName string, editField *tview.InputField, common *Common) tview.Primitive {
	if self.getIntentPath("am") == "" {
		return myButton(self.setButtonLabel(fieldName+" -> Browser(GooleMaps)", common)).SetSelectedFunc(func() {
			webbrowser.Open("http://maps.google.co.jp/maps?q=" + editField.GetText())
		})
	} else {
		return myButton(self.setButtonLabel(fieldName+" -> Intent(GooleMaps)", common)).SetSelectedFunc(func() {
			parm := strings.Split("start -a android.intent.action.VIEW -d geo:0,0?q="+editField.GetText(), " ")
			_ = exec.Command("am", parm...).Run()
		})

	}
}

// -------------------------------------------------
//
//	Execute update, delete, insert
//
// -------------------------------------------------
func (self *Detail) execute(pages *tview.Pages, listItem *listdb.ListItem, common *Common, msg string) {
	flag := msg[:3]
	execute := MyMessageBox(self.app, msg, 30, 7, 2, 3, true).
		SetDoneFunc(func(buttonLabel string, inputString string) {
			if buttonLabel == "Yes" {
				manager := self.connectDB(common)
				if flag == "Upd" {
					manager.Update(common.tableName, listItem.ID, listItem)
				} else if flag == "Del" {
					manager.Delete(common.tableName, listItem.ID)

				} else if flag == "Ins" {
					manager.Insert(common.tableName, listItem)
				}
				manager.Close()
			}
			pages.RemovePage("execute")
			if flag == "Del" || flag == "Ins" {
				NewMainList().run(self.app, common)
			} else {
				self.display(common)
			}
		})

	execute.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			pages.RemovePage("execute")
		}
		return event
	})
	pages.AddPage("execute", execute, false, true)
}

// -------------------------------------------------
//
//	Update
//
// -------------------------------------------------
func (self *Detail) update(pages *tview.Pages, listItem *listdb.ListItem, common *Common) {
	self.execute(pages, listItem, common, "Update record ?")
}

// -------------------------------------------------
//
//	Delete
//
// -------------------------------------------------
func (self *Detail) delete(pages *tview.Pages, listItem *listdb.ListItem, common *Common) {
	self.execute(pages, listItem, common, "Delete record ?")
}

// -------------------------------------------------
//
//	Insert
//
// -------------------------------------------------
func (self *Detail) insert(pages *tview.Pages, listItem *listdb.ListItem, common *Common) {
	self.execute(pages, listItem, common, "Insert record ?")
}

// -------------------------------------------------
//
//	Connect DB
//
// -------------------------------------------------
func (self *Detail) connectDB(common *Common) listdb.Manager {
	manager := listdb.GetManager(common.databaseName)
	err := manager.Connect(common.databaseName, common.connectString)
	if err != nil {
		panic(err)
	}
	return manager
}

// -------------------------------------------------
//
//	ListItem for update
//
// -------------------------------------------------
func (self *Detail) createListItem(id int, category string, field01 string, field02 string, note string) *listdb.ListItem {
	listItem := new(listdb.ListItem)
	listItem.ID = id
	listItem.Category = category
	listItem.Field01 = field01
	listItem.Field02 = field02
	listItem.Note = note
	return listItem
}

// -------------------------------------------------
//
//	detail body
//
// -------------------------------------------------
func (self *Detail) detailBody(pages *tview.Pages, header *tview.Flex, footer *tview.Flex, common *Common) *tview.Flex {
	isLastItem = false
	manager := self.connectDB(common)
	listdb := manager.SearchDB(common.tableName, common.category, common.search, common.selectedItem, 2)
	listdata := listdb.GetListData()
	if len(listdata) < 2 {
		isLastItem = true
	}
	manager.Close()
	//--------------------------------------------

	var focusPrimitives []tview.Primitive
	body := tview.NewFlex().SetDirection(tview.FlexRow)
	btnCategory := myButton(self.setButtonLabel(listdata[0].Category, common))
	btnCategory.SetSelectedFunc(func() {
		self.getCategoryList(pages, common, listdata[0].Category, btnCategory)
	})
	focusPrimitives = append(focusPrimitives, btnCategory)
	editField01 := myEdit(listdata[0].Field01, common.cols)
	editField02 := myEdit(listdata[0].Field02, common.cols)

	var lblField01 tview.Primitive
	if strings.ToLower(listdb.FieldName01) == "address" || strings.ToLower(listdb.FieldName01) == "map" {
		lblField01 = self.setFieldIntent(listdb.FieldName01, editField01, common)
		focusPrimitives = append(focusPrimitives, lblField01)
	} else {
		lblField01 = myLabel(listdb.FieldName01)
	}
	focusPrimitives = append(focusPrimitives, editField01)

	var lblField02 tview.Primitive
	if strings.ToLower(listdb.FieldName02) == "address" || strings.ToLower(listdb.FieldName02) == "map" {
		lblField02 = self.setFieldIntent(listdb.FieldName02, editField02, common)
		focusPrimitives = append(focusPrimitives, lblField02)
	} else {
		lblField02 = myLabel(listdb.FieldName02)
	}
	focusPrimitives = append(focusPrimitives, editField02)

	editNote := mytview.NewMyTextArea().SetText(listdata[0].Note)
	//editNote := tview.NewTextView().SetText(listdata[0].Note)

	focusPrimitives = append(focusPrimitives, editNote)

	body.AddItem(myLabel("ID"), 1, 0, false)
	body.AddItem(myLabel(fmt.Sprintf("%d", listdata[0].ID)), 1, 0, false)
	body.AddItem(myLabel("Category"), 1, 0, false)
	body.AddItem(btnCategory, 1, 0, false)
	body.AddItem(lblField01, 1, 0, false)
	body.AddItem(editField01, 1, 0, false)
	body.AddItem(lblField02, 1, 0, false)
	body.AddItem(editField02, 1, 0, false)
	body.AddItem(myLabel("Note"), 1, 0, false)
	body.AddItem(editNote, 0, 1, false)

	// ------------------------------
	// Footer
	// ------------------------------
	btnU := myButton("<U>").SetSelectedFunc(func() {
		listItem := self.createListItem(listdata[0].ID, btnCategory.GetLabel(), editField01.GetText(), editField02.GetText(), editNote.GetText())
		self.update(pages, listItem, common)
	})
	btnD := myButton("<D>").SetSelectedFunc(func() {
		listItem := self.createListItem(listdata[0].ID, btnCategory.GetLabel(), editField01.GetText(), editField02.GetText(), editNote.GetText())
		self.delete(pages, listItem, common)
	})
	btnI := myButton("<I>").SetSelectedFunc(func() {
		listItem := self.createListItem(listdata[0].ID, btnCategory.GetLabel(), editField01.GetText(), editField02.GetText(), editNote.GetText())
		self.insert(pages, listItem, common)
	})
	footer.AddItem(btnU, 0, 1, true)
	footer.AddItem(btnD, 0, 1, true)
	footer.AddItem(btnI, 0, 1, true)

	// ------------------------------
	// InputCapture on Body
	// ------------------------------
	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {

		case tcell.KeyTab:
			mySetFocus(self.app, focusPrimitives, false)
			return nil

		case tcell.KeyBacktab:
			if btnCategory.HasFocus() {
				myFlexFocus(self.app, header, false)
			} else {
				mySetFocus(self.app, focusPrimitives, true)
				return nil
			}

		case tcell.KeyDown:
			if !editNote.HasFocus() {
				mySetFocus(self.app, focusPrimitives, false)
				return nil
			} else if editNote.LastLine() {
				myFlexFocus(self.app, footer, false)
				//mySetFocus(self.app, focusPrimitives, false)
				return nil
			}

		case tcell.KeyUp:
			if btnCategory.HasFocus() {
				myFlexFocus(self.app, header, false)
				return nil
			} else if !editNote.HasFocus() {
				mySetFocus(self.app, focusPrimitives, true)
				return nil
			} else if editNote.FirstLine() {
				mySetFocus(self.app, focusPrimitives, true)
				return nil
			}

		}
		return event
	})

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
			mySetFocus(self.app, focusPrimitives, false)
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				self.exit()
			case 'r', 'R':
				NewMainList().run(self.app, common)
			case 'n', 'N':
				if common.tableName != "" {
					self.nextPage(common)
				}
			case 'p', 'P':
				if common.tableName != "" {
					self.priorPage(common)
				}
			}

		}
		return event
	})

	// ------------------------------
	// InputCapture on Footer
	// ------------------------------
	footer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight, tcell.KeyTab:
			myFlexFocus(self.app, footer, false)
			return nil
		case tcell.KeyLeft, tcell.KeyBacktab:
			myFlexFocus(self.app, footer, true)
			return nil
		case tcell.KeyDown:
			myFlexFocus(self.app, header, true)
			return nil
		case tcell.KeyUp:
			mySetFocus(self.app, focusPrimitives, false)
			return nil

		}
		return event
	})

	return body
}

// -------------------------------------------------
//
//	format screen
//
// -------------------------------------------------
func (self *Detail) doformat(common *Common) tview.Primitive {
	pages := tview.NewPages()

	// ------------------------------
	//  Header
	// ------------------------------
	header := tview.NewFlex()

	// Return to MainList
	btnR := myButton("<R>").SetSelectedFunc(func() {
		NewMainList().run(self.app, common)
	})
	header.AddItem(btnR, 0, 1, true)

	if !self.lastPage() {
		btnN := myButton("<N>").SetSelectedFunc(func() {
			self.nextPage(common)
		})
		header.AddItem(btnN, 0, 1, true)
	} else {
		header.AddItem(nil, 0, 1, true)
	}
	if !self.firstPage(common) {
		btnP := myButton("<P>").SetSelectedFunc(func() {
			self.priorPage(common)
		})
		header.AddItem(btnP, 0, 1, true)
	} else {
		header.AddItem(nil, 0, 1, true)
	}
	btnQ := myButton("<Q>").SetSelectedFunc(func() {
		self.exit()
	})
	header.AddItem(btnQ, 0, 1, true)

	// ------------------------------
	//  Body
	// ------------------------------
	footer := tview.NewFlex()
	body := self.detailBody(pages, header, footer, common)

	detail := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, true).
		AddItem(body, 0, 1, false).
		AddItem(footer, 1, 0, false)

	pages.AddPage("detail", detail, true, true)
	return pages
}
