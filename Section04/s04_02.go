//
// s04_02.go
//
package main
import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)
func myButton(label string) *tview.Button {
	button := tview.NewButton(label)
	button.SetBackgroundColor(tcell.ColorBlack)
	button.SetLabelColor(tcell.ColorYellow).SetLabelColorActivated(tcell.ColorBlack).SetBackgroundColorActivated(tcell.ColorYellow)
	return button
}

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

func doformat() {
	app := tview.NewApplication()
	header := tview.NewFlex()
	footer := tview.NewTextView().SetText("これはフッター")
	body := tview.NewTextView().SetText("Buttonテスト").SetTextAlign(tview.AlignCenter)

	btnT := myButton("<T>").SetSelectedFunc(func() {
		footer.SetText("Pushed T")
	})
	btnS := myButton("<S>").SetSelectedFunc(func() {
		footer.SetText("Pushed S")
	})
	btnQ := myButton("<Q>").SetSelectedFunc(func() {
		app.Stop()
	})

	buttons := []tview.Primitive{
		btnT,
		btnS,
		btnQ,
	}
	header.AddItem(btnT, 0, 1, true)
	header.AddItem(btnS, 0, 1, true)
	header.AddItem(btnQ, 0, 1, true)
	header.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			mySetFocus(app, buttons, false)
		case tcell.KeyLeft:
			mySetFocus(app, buttons, true)
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				app.Stop()
			}
		}
		body.SetText(string(event.Rune()))
		return event
	})

	main := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, true).
		AddItem(body, 0, 1, false).
		AddItem(footer, 1, 0, false)

	if err := app.SetRoot(main, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func main() {
	doformat()
}
