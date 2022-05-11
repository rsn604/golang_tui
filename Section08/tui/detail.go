package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// -------------------------------------------------
//  detail body
// -------------------------------------------------
func (self *Detail) detailBody(pages *tview.Pages, header *tview.Flex, footer *tview.Flex, common *Common) *tview.Flex {
	var focusPrimitives []tview.Primitive
	body := tview.NewFlex().SetDirection(tview.FlexRow)
	btnCategory := myButton(common.category)
	editField01 := myEdit("フィールド01", common.cols)
	editField02 := myEdit("フィールド02", common.cols)
	focusPrimitives = append(focusPrimitives, btnCategory)
	focusPrimitives = append(focusPrimitives, editField01)
	focusPrimitives = append(focusPrimitives, editField02)
	editNote := tview.NewTextView().SetText("FirstLine\nNextLine")

	body.AddItem(myLabel("ID"), 1, 0, false)
	body.AddItem(myLabel(fmt.Sprintf("%d", common.selectedItem-1)), 1, 0, false)

	body.AddItem(myLabel("Category"), 1, 0, false)
	body.AddItem(btnCategory, 1, 0, false)
	body.AddItem(myLabel("Field01"), 1, 0, false)
	body.AddItem(editField01, 1, 0, false)
	body.AddItem(myLabel("Field02"), 1, 0, false)
	body.AddItem(editField02, 1, 0, false)
	body.AddItem(myLabel("Note"), 1, 0, false)
	body.AddItem(editNote, 0, 1, true)

	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {

		case tcell.KeyTab:
			mySetFocus(self.app, focusPrimitives, false)
			return nil
		case tcell.KeyBacktab:
			if btnCategory.HasFocus() {
				myFlexFocus(self.app, header, false)
				return nil
			} else {
				mySetFocus(self.app, focusPrimitives, true)
				return nil
			}
		case tcell.KeyDown:
			mySetFocus(self.app, focusPrimitives, false)
			return nil
		case tcell.KeyUp:
			if btnCategory.HasFocus() {
				myFlexFocus(self.app, header, false)
				return nil
			} else {
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
			}
		}
		return event
	})

	return body
}

// -------------------------------------------------
//  format screen
// -------------------------------------------------
func (self *Detail) doformat(common *Common) tview.Primitive {
	pages := tview.NewPages()

	header := tview.NewFlex()
	btnR := myButton("<R>").SetSelectedFunc(func() {
		NewMainList().run(self.app, common)
	})
	header.AddItem(btnR, 0, 1, true)
	btnQ := myButton("<Q>").SetSelectedFunc(func() {
		self.exit()
	})
	header.AddItem(btnQ, 0, 1, true)

	footer := tview.NewFlex()
	body := self.detailBody(pages, header, footer, common)
	detail := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, true).
		AddItem(body, 0, 1, false).
		AddItem(footer, 1, 0, false)
	pages.AddPage("detail", detail, true, true)
	return pages
}
