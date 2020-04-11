package main // unit: TxtWind

// interface uses: Video

const (
	MAX_TEXT_WINDOW_LINES   = 1024
	MAX_RESOURCE_DATA_FILES = 24
)

type (
	TTextWindowLine  string
	TTextWindowState struct {
		Selectable     bool
		LineCount      int16
		LinePos        int16
		Lines          [MAX_TEXT_WINDOW_LINES]*TTextWindowLine
		Hyperlink      string
		Title          TTextWindowLine
		LoadedFilename string
		ScreenCopy     [25]string
	}
	TResourceDataHeader struct {
		EntryCount int16
		Name       [MAX_RESOURCE_DATA_FILES]string
		FileOffset [MAX_RESOURCE_DATA_FILES]int32
	}
)

var (
	TextWindowX, TextWindowY          int16
	TextWindowWidth, TextWindowHeight int16
	TextWindowStrInnerEmpty           TVideoLine
	TextWindowStrText                 TVideoLine
	TextWindowStrInnerLine            TVideoLine
	TextWindowStrTop                  TVideoLine
	TextWindowStrBottom               TVideoLine
	TextWindowStrSep                  TVideoLine
	TextWindowStrInnerSep             TVideoLine
	TextWindowStrInnerArrows          TVideoLine
	TextWindowRejected                bool
	ResourceDataFileName              string
	ResourceDataHeader                TResourceDataHeader
	OrderPrintId                      *string
)

// implementation uses: Crt, Input, Printer

func UpCaseString(input string) (UpCaseString string) {
	var i int16
	for i = 1; i <= Length(input); i++ {
		input[i] = UpCase(input[i])
	}
	UpCaseString = input
	return
}

func TextWindowInitState(state *TTextWindowState) {
	state.LineCount = 0
	state.LinePos = 1
	state.LoadedFilename = ""

}

func TextWindowDrawTitle(color int16, title TTextWindowLine) {
	VideoWriteText(TextWindowX+2, TextWindowY+1, color, TextWindowStrInnerEmpty)
	VideoWriteText(TextWindowX+((TextWindowWidth-Length(title))/2), TextWindowY+1, color, title)
}

func TextWindowDrawOpen(state *TTextWindowState) {
	var ix, iy int16
	for iy = 1; iy <= (TextWindowHeight + 1); iy++ {
		VideoMove(TextWindowX, iy+TextWindowY-1, TextWindowWidth, &state.ScreenCopy[iy-1], false)
	}
	for iy = (TextWindowHeight / 2); iy >= 0; iy-- {
		VideoWriteText(TextWindowX, TextWindowY+iy+1, 0x0F, TextWindowStrText)
		VideoWriteText(TextWindowX, TextWindowY+TextWindowHeight-iy-1, 0x0F, TextWindowStrText)
		VideoWriteText(TextWindowX, TextWindowY+iy, 0x0F, TextWindowStrTop)
		VideoWriteText(TextWindowX, TextWindowY+TextWindowHeight-iy, 0x0F, TextWindowStrBottom)
		Delay(25)
	}
	VideoWriteText(TextWindowX, TextWindowY+2, 0x0F, TextWindowStrSep)
	TextWindowDrawTitle(0x1E, state.Title)

}

func TextWindowDrawClose(state *TTextWindowState) {
	var (
		ix, iy     int16
		unk1, unk2 int16
	)
	for iy = 0; iy <= (TextWindowHeight / 2); iy++ {
		VideoWriteText(TextWindowX, TextWindowY+iy, 0x0F, TextWindowStrTop)
		VideoWriteText(TextWindowX, TextWindowY+TextWindowHeight-iy, 0x0F, TextWindowStrBottom)
		Delay(18)
		VideoMove(TextWindowX, TextWindowY+iy, TextWindowWidth, &state.ScreenCopy[(iy+1)-1], true)
		VideoMove(TextWindowX, TextWindowY+TextWindowHeight-iy, TextWindowWidth, &state.ScreenCopy[((TextWindowHeight-iy)+1)-1], true)
	}

}

func TextWindowDrawLine(state *TTextWindowState, lpos int16, withoutFormatting, viewingFile bool) {
	var (
		lineY                        int16
		textOffset, textColor, textX int16
	)
	lineY = ((TextWindowY + lpos) - state.LinePos) + (TextWindowHeight / 2) + 1
	if lpos == state.LinePos {
		VideoWriteText(TextWindowX+2, lineY, 0x1C, TextWindowStrInnerArrows)
	} else {
		VideoWriteText(TextWindowX+2, lineY, 0x1E, TextWindowStrInnerEmpty)
	}
	if (lpos > 0) && (lpos <= state.LineCount) {
		if withoutFormatting {
			VideoWriteText(TextWindowX+4, lineY, 0x1E, state.Lines[lpos-1])
		} else {
			textOffset = 1
			textColor = 0x1E
			textX = TextWindowX + 4
			if Length(state.Lines[lpos-1]) > 0 {
				switch state.Lines[lpos-1][0] {
				case '!':
					textOffset = Pos(';', state.Lines[lpos-1]) + 1
					VideoWriteText(textX+2, lineY, 0x1D, '\x10')
					textX = textX + 5
					textColor = 0x1F
				case ':':
					textOffset = Pos(';', state.Lines[lpos-1]) + 1
					textColor = 0x1F
				case '$':
					textOffset = 2
					textColor = 0x1F
					textX = (textX - 4) + ((TextWindowWidth - Length(state.Lines[lpos-1])) / 2)
				}
			}
			if textOffset > 0 {
				VideoWriteText(textX, lineY, textColor, Copy(state.Lines[lpos-1], textOffset, Length(state.Lines[lpos-1])-textOffset+1))
			}
		}
	} else if (lpos == 0) || (lpos == (state.LineCount + 1)) {
		VideoWriteText(TextWindowX+2, lineY, 0x1E, TextWindowStrInnerSep)
	} else if (lpos == -4) && viewingFile {
		VideoWriteText(TextWindowX+2, lineY, 0x1A, "   Use            to view text,")
		VideoWriteText(TextWindowX+2+7, lineY, 0x1F, "\x18 \x19, Enter")
	} else if (lpos == -3) && viewingFile {
		VideoWriteText(TextWindowX+2+1, lineY, 0x1A, "                 to print.")
		VideoWriteText(TextWindowX+2+12, lineY, 0x1F, "Alt-P")
	}

}

func TextWindowDraw(state *TTextWindowState, withoutFormatting, viewingFile bool) {
	var (
		i    int16
		unk1 int16
	)
	for i = 0; i <= (TextWindowHeight - 4); i++ {
		TextWindowDrawLine(state, state.LinePos-(TextWindowHeight/2)+i+2, withoutFormatting, viewingFile)
	}
	TextWindowDrawTitle(0x1E, state.Title)
}

func TextWindowAppend(state *TTextWindowState, line TTextWindowLine) {
	state.LineCount = state.LineCount + 1
	New(state.Lines[state.LineCount-1])
	state.Lines[state.LineCount-1] = line

}

func TextWindowFree(state *TTextWindowState) {
	for state.LineCount > 0 {
		Dispose(state.Lines[state.LineCount-1])
		state.LineCount = state.LineCount - 1
	}
	state.LoadedFilename = ""

}

func TextWindowPrint(state *TTextWindowState) {
	var (
		iLine, iChar int16
		line         string
	)
	Rewrite(Lst)
	for iLine = 1; iLine <= state.LineCount; iLine++ {
		line = state.Lines[iLine-1]
		if Length(line) > 0 {
			switch line[1] {
			case '$':
				line = Delete(line, 1, 1)
				for iChar = ((80 - Length(line)) / 2); iChar >= 1; iChar-- {
					line = ' ' + line
				}
			case '!', ':':
				iChar = Pos(';', line)
				if iChar > 0 {
					line = Delete(line, 1, iChar)
				} else {
					line = ""
				}
			default:
				line = "          " + line
			}
		}
		WriteLn(Lst, line)
		if IOResult != 0 {
			Close(Lst)
			return
		}
	}
	if state.LoadedFilename == "ORDER.HLP" {
		WriteLn(Lst, OrderPrintId)
	}
	Write(Lst, Chr(12))
	Close(Lst)

}

func TextWindowSelect(state *TTextWindowState, hyperlinkAsSelect, viewingFile bool) {
	var (
		newLinePos   int16
		unk1         int16
		iLine, iChar int16
		pointerStr   string
	)
	TextWindowRejected = false
	state.Hyperlink = ""
	TextWindowDraw(state, false, viewingFile)
	for {
		InputUpdate()
		newLinePos = state.LinePos
		if InputDeltaY != 0 {
			newLinePos = newLinePos + InputDeltaY
		} else if InputShiftPressed || (InputKeyPressed == KEY_ENTER) {
			InputShiftAccepted = true
			if (state.Lines[state.LinePos-1][0]) == '!' {
				pointerStr = Copy(state.Lines[state.LinePos-1], 2, Length(state.Lines[state.LinePos-1])-1)
				if Pos(';', pointerStr) > 0 {
					pointerStr = Copy(pointerStr, 1, Pos(';', pointerStr)-1)
				}
				if pointerStr[0] == '-' {
					pointerStr = Delete(pointerStr, 1, 1)
					TextWindowFree(state)
					TextWindowOpenFile(pointerStr, state)
					if state.LineCount == 0 {
						return
					} else {
						viewingFile = true
						newLinePos = state.LinePos
						TextWindowDraw(state, false, viewingFile)
						InputKeyPressed = '\x00'
						InputShiftPressed = false
					}
				} else {
					if hyperlinkAsSelect {
						state.Hyperlink = pointerStr
					} else {
						pointerStr = ':' + pointerStr
						for iLine = 1; iLine <= state.LineCount; iLine++ {
							if Length(pointerStr) > Length(state.Lines[iLine-1]) {
							} else {
								for iChar = 1; iChar <= Length(pointerStr); iChar++ {
									if UpCase(pointerStr[iChar-1]) != UpCase(state.Lines[iLine-1][iChar-1]) {
										goto LabelNotMatched
									}
								}
								newLinePos = iLine
								InputKeyPressed = '\x00'
								InputShiftPressed = false
								goto LabelMatched
							LabelNotMatched:
							}
						}
					}
				}
			}
		} else {
			if InputKeyPressed == KEY_PAGE_UP {
				newLinePos = state.LinePos - TextWindowHeight + 4
			} else if InputKeyPressed == KEY_PAGE_DOWN {
				newLinePos = state.LinePos + TextWindowHeight - 4
			} else if InputKeyPressed == KEY_ALT_P {
				TextWindowPrint(state)
			}

		}

	LabelMatched:
		if newLinePos < 1 {
			newLinePos = 1
		} else if newLinePos > state.LineCount {
			newLinePos = state.LineCount
		}

		if newLinePos != state.LinePos {
			state.LinePos = newLinePos
			TextWindowDraw(state, false, viewingFile)
			if (state.Lines[state.LinePos-1][0]) == '!' {
				if hyperlinkAsSelect {
					TextWindowDrawTitle(0x1E, "\xaePress ENTER to select this\xaf")
				} else {
					TextWindowDrawTitle(0x1E, "\xaePress ENTER for more info\xaf")
				}
			}
		}
		if InputJoystickMoved {
			Delay(35)
		}
		if (InputKeyPressed == KEY_ESCAPE) || (InputKeyPressed == KEY_ENTER) || InputShiftPressed {
			break
		}
	}
	if InputKeyPressed == KEY_ESCAPE {
		InputKeyPressed = '\x00'
		TextWindowRejected = true
	}

}

func TextWindowEdit(state *TTextWindowState) {
	var (
		newLinePos int16
		insertMode bool
		charPos    int16
		i          int16
	)
	DeleteCurrLine := func() {
		var i int16
		if state.LineCount > 1 {
			Dispose(state.Lines[state.LinePos-1])
			for i = (state.LinePos + 1); i <= state.LineCount; i++ {
				state.Lines[(i-1)-1] = state.Lines[i-1]
			}
			state.LineCount = state.LineCount - 1
			if state.LinePos > state.LineCount {
				newLinePos = state.LineCount
			} else {
				TextWindowDraw(state, true, false)
			}
		} else {
			state.Lines[0] = ""
		}

	}

	if state.LineCount == 0 {
		TextWindowAppend(state, "")
	}
	insertMode = true
	state.LinePos = 1
	charPos = 1
	TextWindowDraw(state, true, false)
	for {
		if insertMode {
			VideoWriteText(77, 14, 0x1E, "on ")
		} else {
			VideoWriteText(77, 14, 0x1E, "off")
		}
		if charPos >= (Length(state.Lines[state.LinePos-1]) + 1) {
			charPos = Length(state.Lines[state.LinePos-1]) + 1
			VideoWriteText(charPos+TextWindowX+3, TextWindowY+(TextWindowHeight/2)+1, 0x70, ' ')
		} else {
			VideoWriteText(charPos+TextWindowX+3, TextWindowY+(TextWindowHeight/2)+1, 0x70, state.Lines[state.LinePos-1][charPos-1])
		}
		InputReadWaitKey()
		newLinePos = state.LinePos
		switch InputKeyPressed {
		case KEY_UP:
			newLinePos = state.LinePos - 1
		case KEY_DOWN:
			newLinePos = state.LinePos + 1
		case KEY_PAGE_UP:
			newLinePos = state.LinePos - TextWindowHeight + 4
		case KEY_PAGE_DOWN:
			newLinePos = state.LinePos + TextWindowHeight - 4
		case KEY_RIGHT:
			charPos = charPos + 1
			if charPos > (Length(state.Lines[state.LinePos-1]) + 1) {
				charPos = 1
				newLinePos = state.LinePos + 1
			}
		case KEY_LEFT:
			charPos = charPos - 1
			if charPos < 1 {
				charPos = TextWindowWidth
				newLinePos = state.LinePos - 1
			}
		case KEY_ENTER:
			if state.LineCount < MAX_TEXT_WINDOW_LINES {
				for i = state.LineCount; i >= (state.LinePos + 1); i-- {
					state.Lines[(i+1)-1] = state.Lines[i-1]
				}
				New(state.Lines[(state.LinePos+1)-1])
				state.Lines[(state.LinePos+1)-1] = Copy(state.Lines[state.LinePos-1], charPos, Length(state.Lines[state.LinePos-1])-charPos+1)
				state.Lines[state.LinePos-1] = Copy(state.Lines[state.LinePos-1], 1, charPos-1)
				newLinePos = state.LinePos + 1
				charPos = 1
				state.LineCount = state.LineCount + 1
			}
		case KEY_BACKSPACE:
			if charPos > 1 {
				state.Lines[state.LinePos-1] = Copy(state.Lines[state.LinePos-1], 1, charPos-2) + Copy(state.Lines[state.LinePos-1], charPos, Length(state.Lines[state.LinePos-1])-charPos+1)
				charPos = charPos - 1
			} else if Length(state.Lines[state.LinePos-1]) == 0 {
				DeleteCurrLine()
				newLinePos = state.LinePos - 1
				charPos = TextWindowWidth
			}

		case KEY_INSERT:
			insertMode = !insertMode
		case KEY_DELETE:
			state.Lines[state.LinePos-1] = Copy(state.Lines[state.LinePos-1], 1, charPos-1) + Copy(state.Lines[state.LinePos-1], charPos+1, Length(state.Lines[state.LinePos-1])-charPos)
		case KEY_CTRL_Y:
			DeleteCurrLine()
		default:
			if (InputKeyPressed >= ' ') && (charPos < (TextWindowWidth - 7)) {
				if !insertMode {
					state.Lines[state.LinePos-1] = Copy(state.Lines[state.LinePos-1], 1, charPos-1) + InputKeyPressed + Copy(state.Lines[state.LinePos-1], charPos+1, Length(state.Lines[state.LinePos-1])-charPos)
					charPos = charPos + 1
				} else {
					if Length(state.Lines[state.LinePos-1]) < (TextWindowWidth - 8) {
						state.Lines[state.LinePos-1] = Copy(state.Lines[state.LinePos-1], 1, charPos-1) + InputKeyPressed + Copy(state.Lines[state.LinePos-1], charPos, Length(state.Lines[state.LinePos-1])-charPos+1)
						charPos = charPos + 1
					}
				}
			}
		}
		if newLinePos < 1 {
			newLinePos = 1
		} else if newLinePos > state.LineCount {
			newLinePos = state.LineCount
		}

		if newLinePos != state.LinePos {
			state.LinePos = newLinePos
			TextWindowDraw(state, true, false)
		} else {
			TextWindowDrawLine(state, state.LinePos, true, false)
		}
		if InputKeyPressed == KEY_ESCAPE {
			break
		}
	}
	if Length(state.Lines[state.LineCount-1]) == 0 {
		Dispose(state.Lines[state.LineCount-1])
		state.LineCount = state.LineCount - 1
	}

}

func TextWindowOpenFile(filename TTextWindowLine, state *TTextWindowState) {
	var (
		f        FILE
		tf       text
		i        int16
		entryPos int16
		retVal   bool
		line     *string
		lineLen  byte
	)
	retVal = true
	for i = 1; i <= Length(filename); i++ {
		retVal = retVal && (filename[i-1] != '.')
	}
	if retVal {
		filename = filename + ".HLP"
	}
	if filename[0] == '*' {
		filename = Copy(filename, 2, Length(filename)-1)
		entryPos = -1
	} else {
		entryPos = 0
	}
	TextWindowInitState(state)
	state.LoadedFilename = UpCaseString(filename)
	if ResourceDataHeader.EntryCount == 0 {
		Assign(f, ResourceDataFileName)
		Reset(f, 1)
		if IOResult == 0 {
			BlockRead(f, ResourceDataHeader, SizeOf(ResourceDataHeader))
		}
		if IOResult != 0 {
			ResourceDataHeader.EntryCount = -1
		}
		Close(f)
	}
	if entryPos == 0 {
		for i = 1; i <= ResourceDataHeader.EntryCount; i++ {
			if UpCaseString(ResourceDataHeader.Name[i-1]) == UpCaseString(filename) {
				entryPos = i
			}
		}
	}
	if entryPos <= 0 {
		Assign(tf, filename)
		Reset(tf)
		for (IOResult == 0) && (!Eof(tf)) {
			Inc(state.LineCount)
			New(state.Lines[state.LineCount-1])
			ReadLn(tf, state.Lines[state.LineCount-1])
		}
		Close(tf)
	} else {
		Assign(f, ResourceDataFilename)
		Reset(f, 1)
		Seek(f, ResourceDataHeader.FileOffset[entryPos-1])
		if IOResult == 0 {
			retVal = true
			for (IOResult == 0) && retVal {
				Inc(state.LineCount)
				New(state.Lines[state.LineCount-1])
				BlockRead(f, state.Lines[state.LineCount-1], 1)
				line = Ptr(Seg(state.Lines[state.LineCount-1]), Ofs(state.Lines[state.LineCount-1])+1)
				lineLen = Ord(state.Lines[state.LineCount-1][-1])
				if lineLen == 0 {
					state.Lines[state.LineCount-1] = ""
				} else {
					BlockRead(f, line, Ord(state.Lines[state.LineCount-1][-1]))
				}
				if state.Lines[state.LineCount-1] == '@' {
					retVal = false
					state.Lines[state.LineCount-1] = ""
				}
			}
			Close(f)
		}
	}

}

func TextWindowSaveFile(filename TTextWindowLine, state *TTextWindowState) {
	var (
		f text
		i int16
	)
	Assign(f, filename)
	Rewrite(f)
	if IOResult != 0 {
		return
	}
	for i = 1; i <= state.LineCount; i++ {
		WriteLn(f, state.Lines[i-1])
		if IOResult != 0 {
			return
		}
	}
	Close(f)

}

func TextWindowDisplayFile(filename, title string) {
	var state TTextWindowState
	state.Title = title
	TextWindowOpenFile(filename, &state)
	state.Selectable = false
	if state.LineCount > 0 {
		TextWindowDrawOpen(&state)
		TextWindowSelect(&state, false, true)
		TextWindowDrawClose(&state)
	}
	TextWindowFree(&state)
}

func TextWindowInit(x, y, width, height int16) {
	var i int16
	TextWindowX = x
	TextWindowWidth = width
	TextWindowY = y
	TextWindowHeight = height
	TextWindowStrInnerEmpty = ""
	TextWindowStrInnerLine = ""
	for i = 1; i <= (TextWindowWidth - 5); i++ {
		TextWindowStrInnerEmpty = TextWindowStrInnerEmpty + ' '
		TextWindowStrInnerLine = TextWindowStrInnerLine + 'Í'
	}
	TextWindowStrTop = "\xc6\xd1" + TextWindowStrInnerLine + 'Ñ' + 'µ'
	TextWindowStrBottom = "\xc6\xcf" + TextWindowStrInnerLine + 'Ï' + 'µ'
	TextWindowStrSep = " \xc6" + TextWindowStrInnerLine + 'µ' + ' '
	TextWindowStrText = " \xb3" + TextWindowStrInnerEmpty + '³' + ' '
	TextWindowStrInnerArrows = TextWindowStrInnerEmpty
	TextWindowStrInnerArrows[0] = '¯'
	TextWindowStrInnerArrows[Length(TextWindowStrInnerArrows)-1] = '®'
	TextWindowStrInnerSep = TextWindowStrInnerEmpty
	for i = 1; i <= (TextWindowWidth / 5); i++ {
		TextWindowStrInnerSep[(i*5+((TextWindowWidth%5)/2))-1] = '\a'
	}
}

func init() {
	ResourceDataFileName = ""
	ResourceDataHeader.EntryCount = 0
}
