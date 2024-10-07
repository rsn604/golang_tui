package tui

import (
	//"listdbg/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"listdbg/mytview"
)

// ----------------------------------------
// Set focus
// ----------------------------------------
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

// ----------------------------------------
// Set focus in Flex
// ----------------------------------------
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
		//app.SetFocus(flex.GetItem(0))
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
	//app.SetFocus(flex.GetItem(0))
	app.SetFocus(flex.GetItem(firstFocus))
}

// ----------------------------------------
// Button
// ----------------------------------------
func myButton(label string) *tview.Button {
	button := tview.NewButton(label)
	button.Box = tview.NewBox().SetBackgroundColor(tcell.ColorBlack)
	button.SetLabelColor(tcell.ColorYellow).SetLabelColorActivated(tcell.ColorBlack).SetBackgroundColorActivated(tcell.ColorYellow)
	return button
}

// ----------------------------------------
// TextView
// ----------------------------------------
// func myLabel(label string) *tview.TextView{
func myLabel(label string) *mytview.MyTextView {
	//return tview.NewTextView().SetTextColor(tcell.ColorAqua).SetText(label).SetTextAlign(tview.AlignLeft)
	return mytview.NewMyTextView().SetTextColor(tcell.ColorAqua).SetText(label).SetTextAlign(tview.AlignLeft)

}

// ----------------------------------------
// InputField
// ----------------------------------------
func myEdit(text string, rows int) *tview.InputField {
	return tview.NewInputField().SetFieldWidth(rows).SetText(text).SetFieldTextColor(tcell.ColorWhite).SetFieldBackgroundColor(tcell.ColorBlack).SetFieldStyle(tcell.StyleDefault.Underline(true))

}

// ----------------------------------------
// SelectBox
// ----------------------------------------
func MySelectBox(items []string, width int, height int, ratioX int, ratioY int, current int) *mytview.MyDialog {
	dialog := mytview.NewMyDialog(width, height, ratioX, ratioY)

	list := tview.NewList().ShowSecondaryText(false).SetSelectedTextColor(tcell.ColorWhite).SetSelectedBackgroundColor(tcell.ColorAqua)
	for _, s := range items {
		list.AddItem(s, "", 0, nil)
	}

	list.SetCurrentItem(current).SetSelectedFunc(func(index int, s string, secondary string, code rune) {
		dialog.SetParm("OK", s)
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(list, 0, 1, true)
	dialog.SetFlex(flex)
	return dialog
}

// ----------------------------------------------
// InputDialog
// ----------------------------------------------
func MyInputDialog(app *tview.Application, search string, width int, height int, ratioX int, ratioY int) *mytview.MyDialog {
	dialog := mytview.NewMyDialog(width, height, ratioX, ratioY)

	input := tview.NewInputField().
		SetLabel("Search : ").SetText(search).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorYellow).
		SetFieldStyle(tcell.StyleDefault.Underline(true)).
		SetFieldWidth(20)
	btnOK := myButton("OK")
	btnCancel := myButton("Cancel")

	btnOK.SetSelectedFunc(func() {
		dialog.SetParm(btnOK.GetLabel(), input.GetText())
	})

	btnCancel.SetSelectedFunc(func() {
		dialog.SetParm(btnCancel.GetLabel(), "")
	})

	buttons := tview.NewFlex().AddItem(nil, 0, 1, false).AddItem(btnOK, 10, 0, true).AddItem(btnCancel, 10, 0, true).AddItem(nil, 0, 1, false)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 1, 0, false).
		AddItem(input, 1, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(buttons, 1, 0, true)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			if !input.HasFocus() {
				myFlexFocus(app, buttons, false)
				return nil
			}
		case tcell.KeyLeft:
			if !input.HasFocus() {
				myFlexFocus(app, buttons, true)
				return nil
			}
		case tcell.KeyUp:
			app.SetFocus(input)
			return nil
		case tcell.KeyDown:
			app.SetFocus(btnOK)
			return nil
		case tcell.KeyEnter:
			if btnCancel.HasFocus() {
				dialog.SetParm(btnCancel.GetLabel(), input.GetText())
			} else {
				dialog.SetParm(btnOK.GetLabel(), input.GetText())
			}
			return nil
		case tcell.KeyEscape:
			dialog.SetParm(btnCancel.GetLabel(), "")
			return nil
		}
		return event
	})

	dialog.SetFlex(flex)
	return dialog
}

// ----------------------------------------------
// MessageBox
// ----------------------------------------------
func MyMessageBox(app *tview.Application, message string, width int, height int, ratioX int, ratioY int) *mytview.MyDialog {
	dialog := mytview.NewMyDialog(width, height, ratioX, ratioY)

	msg := myLabel(message)
	btnYes := myButton("Yes")
	btnNo := myButton("No")

	btnYes.SetSelectedFunc(func() {
		dialog.SetParm(btnYes.GetLabel(), "")
	})

	btnNo.SetSelectedFunc(func() {
		dialog.SetParm(btnNo.GetLabel(), "")
	})
	primitives := []tview.Primitive{
		btnYes,
		btnNo,
	}

	buttons := tview.NewFlex().AddItem(nil, 0, 1, false).AddItem(btnYes, 10, 0, true).AddItem(btnNo, 10, 0, true).AddItem(nil, 0, 1, false)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 1, 0, false).
		AddItem(msg, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(buttons, 1, 0, true)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			mySetFocus(app, primitives, false)
			return nil
		case tcell.KeyLeft:
			mySetFocus(app, primitives, true)
			return nil
		case tcell.KeyEnter:
			if btnNo.HasFocus() {
				dialog.SetParm(btnNo.GetLabel(), "")
			} else {
				dialog.SetParm(btnYes.GetLabel(), "")
			}
			return nil
		case tcell.KeyEscape:
			dialog.SetParm(btnNo.GetLabel(), "")
			return nil
		}
		return event
	})

	dialog.SetFlex(flex)
	return dialog
}
