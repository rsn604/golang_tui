package mytview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ------------------------------------------------
// Dialog
// ------------------------------------------------
type MyDialog struct {
	*tview.Box
	frame  *tview.Frame
	flex   *tview.Flex
	width  int
	height int
	ratioX int
	ratioY int
	border bool
	done   func(buttonLabel string, input string)
}

func NewMyDialog(width int, height int, ratioX int, ratioY int) *MyDialog {
	m := &MyDialog{
		Box:    tview.NewBox(),
		width:  width,
		height: height,
		ratioX: ratioX,
		ratioY: ratioY,
		border: true,
	}
	return m
}

func (m *MyDialog) SetFlex(flex *tview.Flex) *MyDialog {
	m.flex = flex
	m.frame = tview.NewFrame(m.flex).SetBorders(0, 0, 0, 0, 2, 2)

	//m.frame.SetBorder(true).
	m.frame.SetBorder(false).
		SetBackgroundColor(tcell.ColorBlack).
		SetBorderPadding(1, 1, 1, 1)
	return m
}

func (m *MyDialog) SetParm(label string, input string) *MyDialog {
	m.done(label, input)
	return m
}

func (m *MyDialog) SetBorder(border bool) *MyDialog {
	m.border = border
	return m
}

func (m *MyDialog) SetDoneFunc(handler func(buttonLabel string, input string)) *MyDialog {
	m.done = handler
	return m
}

func (m *MyDialog) Draw(screen tcell.Screen) {
	screenWidth, screenHeight := screen.Size()
	x := (screenWidth - m.width) / m.ratioX
	y := (screenHeight - m.height) / m.ratioY
	m.SetRect(x, y, m.width, m.height)
	m.frame.SetRect(x, y, m.width, m.height)
	m.frame.Clear()
	// ----------------------------------------------------------------
	// FrameのBorderだと、漢字2バイト目を壊すので、自前で内部に枠線を書く。
	// ----------------------------------------------------------------
	if m.border {
		m.frame.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			leftX := x + 1
			rightX := x + width
			topY := y + 1
			bottomY := y + height

			screen.SetContent(leftX, y+1, tview.BoxDrawingsLightDownAndRight, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			for cx := leftX + 1; cx < x+width; cx++ {
				screen.SetContent(cx, topY, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}

			screen.SetContent(x+width, y+1, tview.BoxDrawingsLightDownAndLeft, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			for cx := leftX + 1; cx < x+width; cx++ {
				screen.SetContent(cx, bottomY, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}

			screen.SetContent(leftX, y+height, tview.BoxDrawingsLightUpAndRight, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			for cy := topY + 1; cy < y+height; cy++ {
				screen.SetContent(leftX, cy, tview.BoxDrawingsLightVertical, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}

			screen.SetContent(x+width, y+height, tview.BoxDrawingsLightUpAndLeft, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			for cy := topY + 1; cy < y+height; cy++ {
				screen.SetContent(rightX, cy, tview.BoxDrawingsLightVertical, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}
			return x + 1, topY + 1, width - 2, height - 2
		})
	}
	m.frame.Draw(screen)
}

// ===========================================================
/*
func (m *MyDialog) SetFocus(index int) *MyDialog {
	m.flex.SetFocus(index)
	return m
}
*/

func (m *MyDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(m.flex)
}

func (m *MyDialog) HasFocus() bool {
	return m.flex.HasFocus()
}

func (m *MyDialog) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		// Pass mouse events on to the form.
		consumed, capture = m.flex.MouseHandler()(action, event, setFocus)
		if !consumed && action == tview.MouseLeftClick && m.InRect(event.Position()) {
			setFocus(m)
			consumed = true
		}
		return
	})
}

func (m *MyDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if m.flex.HasFocus() {
			if handler := m.flex.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
	})
}
