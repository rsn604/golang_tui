package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

func setList(list *tview.List, from int) *tview.List {
	list.Clear()
	item := getListData(from)
	for i := 0; i < rows-2; i++ {
		list.AddItem(item[i], "", 0, nil)
	}
	return list
}

func nextPage(list *tview.List, from int) (*tview.List, int) {
	if from+(rows-2) < len(listData) {
		from = from + (rows - 2)
		list = setList(list, from)
	}
	return list, from
}

func priorPage(list *tview.List, from int) (*tview.List, int) {
	if from > 1 {
		from = from - (rows - 2)
		list = setList(list, from)
	}
	return list, from
}

func doformat(app *tview.Application, from int) tview.Primitive {
	list := tview.NewList().ShowSecondaryText(false).SetSelectedTextColor(tcell.ColorWhite).SetSelectedBackgroundColor(tcell.ColorAqua).SetSelectedFocusOnly(true)
	list = setList(list, from)

	pages := tview.NewPages()
	footer := tview.NewTextView().SetText("これはフッター")
	header := tview.NewFlex()
	btnN := myButton("<N>").SetSelectedFunc(func() {
		list, from = nextPage(list, from)
	})
	btnP := myButton("<P>").SetSelectedFunc(func() {
		list, from = priorPage(list, from)
	})
	btnQ := myButton("<Q>").SetSelectedFunc(func() {
		app.Stop()
	})
	buttons := []tview.Primitive{
		btnN,
		btnP,
		btnQ,
	}
	header.AddItem(btnN, 0, 1, true)
	header.AddItem(btnP, 0, 1, true)
	header.AddItem(btnQ, 0, 1, true)

	body := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(list, 0, 1, true)

	main := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, true).
		AddItem(body, 0, 1, true).
		AddItem(footer, 1, 0, false)
	pages.AddPage("main", main, true, true)

	header.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			mySetFocus(app, buttons, false)
		case tcell.KeyLeft:
			mySetFocus(app, buttons, true)
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
				mySetFocus(app, buttons, true)
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
			case 'n':
				list, from = nextPage(list, from)
			case 'p':
				list, from = priorPage(list, from)
			}
		}
		return event
	})

	list.SetSelectedFunc(func(index int, s string, secondary string, code rune) {
		footer.SetText("pos:" + strconv.Itoa(index) + " data:" + s)
	})

	return pages
}

func main() {
	app := tview.NewApplication()
	pages := doformat(app, 1)
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
