//
// mywidget.go
//
package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var cols, rows = GetScreenSize()

func GetScreenSize() (int, int) {
	s, _ := tcell.NewScreen()
	s.Init()
	cols, rows := s.Size()
	s.Fini()
	return cols, rows
}

var listData = createListData()

func createListData() []string {
	var listData []string
	var s string
	for i := 0; i < 100; i++ {
		s = "テストデータ" + fmt.Sprintf("%d", i)
		for j := len(s); j < cols; j++ {
			s = s + " "
		}
		listData = append(listData, s)
	}
	return listData
}

func getListData(from int) []string {
	return listData[from-1 : from+rows-2]
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

func myButton(label string) *tview.Button {
	button := tview.NewButton(label)
	button.Box = tview.NewBox().SetBackgroundColor(tcell.ColorBlack)
	button.SetLabelColor(tcell.ColorYellow).SetLabelColorActivated(tcell.ColorBlack).SetBackgroundColorActivated(tcell.ColorYellow)
	return button
}
