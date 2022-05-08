package tui

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

func getListData(common *Common) []string {
	if len(listData) > common.from+common.rows-2 {
		return listData[common.from-1 : common.from+common.rows-2]
	} else {
		return listData[common.from-1:]
	}
}

// -----------------------------------------------
func mySetFocus(app *tview.Application, elements []tview.Primitive, reverse bool) {
	for i, el := range elements {
		if !el.HasFocus() {
			continue
		}

		if reverse {
			i = i - 1
			if i < 0 {
				i = len(elements) - 1
			}
		} else {
			i = i + 1
			i = i % len(elements)
		}

		app.SetFocus(elements[i])
		return
	}
	app.SetFocus(elements[0])
}

func myFlexFocus(app *tview.Application, flex *tview.Flex, reverse bool) {
	count := flex.GetItemCount()
	current := -1
	firstFocus := -1
	for i := 0; i < count; i++ {
		if flex.GetItem(i) == nil {
			continue
		}
		if firstFocus == -1 {
			firstFocus = i
		}
		if flex.GetItem(i).HasFocus() {
			current = i
			break
		}
	}
	if current == -1 {
		app.SetFocus(flex.GetItem(firstFocus))
		return
	}

	if reverse {
		for i := current - 1; i >= 0; i-- {
			if flex.GetItem(i) == nil {
				continue
			}
			app.SetFocus(flex.GetItem(i))
			return
		}
	} else {
		for i := current + 1; i < count; i++ {
			if flex.GetItem(i) == nil {
				continue
			}
			app.SetFocus(flex.GetItem(i))
			return
		}
	}
	app.SetFocus(flex.GetItem(firstFocus))
}

func myButton(label string) *tview.Button {
	button := tview.NewButton(label)
	button.Box = tview.NewBox().SetBackgroundColor(tcell.ColorBlack)
	button.SetLabelColor(tcell.ColorYellow).SetLabelColorActivated(tcell.ColorBlack).SetBackgroundColorActivated(tcell.ColorYellow)
	return button
}
