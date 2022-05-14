package main

import (
	"github.com/gdamore/tcell/v2"
	"listdbg/tui"
	"os"
)

func getScreenSize() (int, int) {
	s, _ := tcell.NewScreen()
	s.Init()
	cols, rows := s.Size()
	s.Fini()
	return cols, rows
}

// -------------------------------------------------
func main() {
	cols, rows := getScreenSize()
	if len(os.Args) == 3 {
		tui.NewMainList().Init(os.Args[1], os.Args[2], cols, rows)
	} else {
		tui.NewMainList().Init("", "", cols, rows)
	}
}
