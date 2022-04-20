//
// s05_02.go
//
package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var listData = createListData()

func createListData() []string {
	var listData []string
	var s string
	for i := 0; i < 100; i++ {
		s = "テストデータ" + fmt.Sprintf("%d", i)
		listData = append(listData, s)
	}
	return listData
}

func getListData() []string {
	return listData
}

func myButton(label string) *tview.Button {
	button := tview.NewButton(label)
	button.SetBackgroundColor(tcell.ColorBlack)
	button.SetLabelColor(tcell.ColorYellow).SetLabelColorActivated(tcell.ColorBlack).SetBackgroundColorActivated(tcell.ColorYellow)
	return button
}

func setList(list *tview.List) *tview.List {
	list.Clear()
	items := getListData()
	for _, item := range items {
		list.AddItem(item, "", 0, nil)
	}
	return list
}

func doformat(app *tview.Application) tview.Primitive {
	pages := tview.NewPages()
	footer := tview.NewTextView().SetText("これはフッター")
	header := tview.NewFlex()
	btnQ := myButton("<Q>").SetSelectedFunc(func() {
		app.Stop()
	})
	header.AddItem(btnQ, 6, 0, true)

	list := tview.NewList().ShowSecondaryText(false).SetSelectedTextColor(tcell.ColorWhite).SetSelectedBackgroundColor(tcell.ColorAqua).SetSelectedFocusOnly(true)
	list = setList(list)

	body := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(list, 0, 1, true)

	main := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, true).
		AddItem(body, 0, 1, true).
		AddItem(footer, 1, 0, false)
	pages.AddPage("main", main, true, true)

	header.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown:
			app.SetFocus(list)
			return nil
		}
		return event
	})
	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			if list.GetCurrentItem() == 0 {
				app.SetFocus(btnQ)
				return nil
			}
		}
		return event
	})
	pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				app.Stop()
			}
		}
		return event
	})
	return pages
}

func main() {
	app := tview.NewApplication()
	pages := doformat(app)
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
