package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

var isLast bool

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
func (self *MainList) setList(list *tview.List, common *Common) *tview.List {
	listdata := getListData(common)
	if len(listdata) <= common.rows-2 {
		isLast = true
	} else {
		isLast = false
	}
	list.Clear()
	for i, s := range listdata {
		if i >= common.rows-2 {
			break
		}
		list.AddItem(s, "", 0, nil)
	}
	return list
}
func (self *MainList) doformat(common *Common) tview.Primitive {
	pages := tview.NewPages()
	header := tview.NewFlex()
	footer := tview.NewTextView()

	list := tview.NewList().ShowSecondaryText(false).SetSelectedTextColor(tcell.ColorWhite).SetSelectedBackgroundColor(tcell.ColorAqua).SetSelectedFocusOnly(true)
	list = self.setList(list, common)
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

	footer.SetTextColor(tcell.ColorYellow).SetText("Footer")

	btnQ := myButton("<Q>").SetSelectedFunc(func() {
		self.exit()
	})
	header.AddItem(btnQ, 0, 1, true)

	body := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(list, 0, 1, true)

	main := tview.NewFlex().SetDirection(tview.FlexRow)

	main.AddItem(header, 1, 0, false)
	main.AddItem(body, 0, 1, true).
		AddItem(footer, 1, 0, false)

	header.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight, tcell.KeyTab:
			myFlexFocus(self.app, header, false)
			return nil
		case tcell.KeyLeft, tcell.KeyBacktab:
			myFlexFocus(self.app, header, true)
			return nil
		case tcell.KeyDown:
			self.app.SetFocus(list)
			return nil
		}
		return event
	})

	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			if list.GetCurrentItem() == 0 {
				myFlexFocus(self.app, header, true)
				return nil
			}
		}
		return event
	})

	main.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				self.exit()

			case 'n', 'N':
				self.nextPage(common)

			case 'p', 'P':
				self.priorPage(common)
			}
		}
		return event
	})

	list.SetSelectedFunc(func(index int, s string, secondary string, code rune) {
		footer.SetText("pos:" + strconv.Itoa(index) + " data:" + s)
	})

	pages.AddPage("main", main, true, true)
	return pages
}
