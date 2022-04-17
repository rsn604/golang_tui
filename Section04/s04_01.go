//
// s04_01.go
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
	header.AddItem(btnT, 0, 1, true)
	header.AddItem(btnS, 0, 1, true)
	header.AddItem(btnQ, 0, 1, true)

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
