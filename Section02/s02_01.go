//
// s02_01.go
//
package main
import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)
func main() {
	app := tview.NewApplication()
	inputField := tview.NewInputField().
		SetLabel("文字を入力: ").
		SetFieldWidth(10).
		SetDoneFunc(func(key tcell.Key) {
			app.Stop()
		})
	if err := app.SetRoot(inputField, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
