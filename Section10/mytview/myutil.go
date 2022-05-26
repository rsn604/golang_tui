package mytview

import (
	"regexp"
	"sort"
	"github.com/gdamore/tcell/v2"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

const (
	AlignLeft = iota
	AlignCenter
	AlignRight
)

// Common regular expressions.
var (
	colorPattern     = regexp.MustCompile(`\[([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([lbidrus]+|\-)?)?)?\]`)
	regionPattern    = regexp.MustCompile(`\["([a-zA-Z0-9_,;: \-\.]*)"\]`)
	escapePattern    = regexp.MustCompile(`\[([a-zA-Z0-9_,;: \-\."#]+)\[(\[*)\]`)
	nonEscapePattern = regexp.MustCompile(`(\[[a-zA-Z0-9_,;: \-\."#]+\[*)\]`)
	boundaryPattern  = regexp.MustCompile(`(([,\.\-:;!\?&#+]|\n)[ \t\f\r]*|([ \t\f\r]+))`)
	spacePattern     = regexp.MustCompile(`\s+`)
)

// Positions of substrings in regular expressions.

const (
	colorForegroundPos = 1
	colorBackgroundPos = 3
	colorFlagPos       = 5
)

func styleFromTag(fgColor, bgColor, attributes string, tagSubstrings []string) (newFgColor, newBgColor, newAttributes string) {
	if tagSubstrings[colorForegroundPos] != "" {
		color := tagSubstrings[colorForegroundPos]
		if color == "-" {
			fgColor = "-"
		} else if color != "" {
			fgColor = color
		}
	}

	if tagSubstrings[colorBackgroundPos-1] != "" {
		color := tagSubstrings[colorBackgroundPos]
		if color == "-" {
			bgColor = "-"
		} else if color != "" {
			bgColor = color
		}
	}

	if tagSubstrings[colorFlagPos-1] != "" {
		flags := tagSubstrings[colorFlagPos]
		if flags == "-" {
			attributes = "-"
		} else if flags != "" {
			attributes = flags
		}
	}

	return fgColor, bgColor, attributes
}

func overlayStyle(style tcell.Style, fgColor, bgColor, attributes string) tcell.Style {
	_, _, defAttr := style.Decompose()

	if fgColor != "" && fgColor != "-" {
		style = style.Foreground(tcell.GetColor(fgColor))
	}

	if bgColor != "" && bgColor != "-" {
		style = style.Background(tcell.GetColor(bgColor))
	}

	if attributes == "-" {
		style = style.Bold(defAttr&tcell.AttrBold > 0).
			Italic(defAttr&tcell.AttrItalic > 0).
			Blink(defAttr&tcell.AttrBlink > 0).
			Reverse(defAttr&tcell.AttrReverse > 0).
			Underline(defAttr&tcell.AttrUnderline > 0).
			Dim(defAttr&tcell.AttrDim > 0)
	} else if attributes != "" {
		style = style.Normal()
		for _, flag := range attributes {
			switch flag {
			case 'l':
				style = style.Blink(true)
			case 'b':
				style = style.Bold(true)
			case 'i':
				style = style.Italic(true)
			case 'd':
				style = style.Dim(true)
			case 'r':
				style = style.Reverse(true)
			case 'u':
				style = style.Underline(true)
			case 's':
				style = style.StrikeThrough(true)
			}
		}
	}

	return style
}

func decomposeString(text string, findColors, findRegions bool) (colorIndices [][]int, colors [][]string, regionIndices [][]int, regions [][]string, escapeIndices [][]int, stripped string, width int) {
	// Shortcut for the trivial case.
	if !findColors && !findRegions {
		return nil, nil, nil, nil, nil, text, stringWidth(text)
	}

	// Get positions of any tags.
	if findColors {
		colorIndices = colorPattern.FindAllStringIndex(text, -1)
		colors = colorPattern.FindAllStringSubmatch(text, -1)
	}
	if findRegions {
		regionIndices = regionPattern.FindAllStringIndex(text, -1)
		regions = regionPattern.FindAllStringSubmatch(text, -1)
	}
	escapeIndices = escapePattern.FindAllStringIndex(text, -1)

	// Because the color pattern detects empty tags, we need to filter them out.
	for i := len(colorIndices) - 1; i >= 0; i-- {
		if colorIndices[i][1]-colorIndices[i][0] == 2 {
			colorIndices = append(colorIndices[:i], colorIndices[i+1:]...)
			colors = append(colors[:i], colors[i+1:]...)
		}
	}

	// Make a (sorted) list of all tags.
	allIndices := make([][3]int, 0, len(colorIndices)+len(regionIndices)+len(escapeIndices))
	for indexType, index := range [][][]int{colorIndices, regionIndices, escapeIndices} {
		for _, tag := range index {
			allIndices = append(allIndices, [3]int{tag[0], tag[1], indexType})
		}
	}
	sort.Slice(allIndices, func(i int, j int) bool {
		return allIndices[i][0] < allIndices[j][0]
	})

	// Remove the tags from the original string.
	var from int
	buf := make([]byte, 0, len(text))
	for _, indices := range allIndices {
		if indices[2] == 2 { // Escape sequences are not simply removed.
			buf = append(buf, []byte(text[from:indices[1]-2])...)
			buf = append(buf, ']')
			from = indices[1]
		} else {
			buf = append(buf, []byte(text[from:indices[0]])...)
			from = indices[1]
		}
	}
	buf = append(buf, text[from:]...)
	stripped = string(buf)

	// Get the width of the stripped string.
	width = stringWidth(stripped)

	return
}

func stringWidth(text string) (width int) {
	g := uniseg.NewGraphemes(text)
	for g.Next() {
		var chWidth int
		for _, r := range g.Runes() {
			chWidth = runewidth.RuneWidth(r)
			if chWidth > 0 {
				break // Our best guess at this point is to use the width of the first non-zero-width rune.
			}
		}
		width += chWidth
	}
	return
}

func iterateString(text string, callback func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool) bool {
	var screenPos int

	gr := uniseg.NewGraphemes(text)
	for gr.Next() {
		r := gr.Runes()
		from, to := gr.Positions()
		width := stringWidth(gr.Str())
		var comb []rune
		if len(r) > 1 {
			comb = r[1:]
		}

		if callback(r[0], comb, from, to-from, screenPos, width) {
			return true
		}

		screenPos += width
	}

	return false
}

func stripTags(text string) string {
	stripped := colorPattern.ReplaceAllStringFunc(text, func(match string) string {
		if len(match) > 2 {
			return ""
		}
		return match
	})
	return escapePattern.ReplaceAllString(stripped, `[$1$2]`)
}
