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
		Lines          [MAX_TEXT_WINDOW_LINES - 1 + 1]*TTextWindowLine
		Hyperlink      string
		Title          TTextWindowLine
		LoadedFilename string
		ScreenCopy     [25 - 1 + 1]string
	}
	TResourceDataHeader struct {
		EntryCount int16
		Name       [MAX_RESOURCE_DATA_FILES - 1 + 1]string
		FileOffset [MAX_RESOURCE_DATA_FILES - 1 + 1]int32
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
	var (
		i int16
	)
	for i := 1; i <= Length(input); i++ {
		input[i] = UpCase(input[i])
	}
	UpCaseString = input
	return
}

func TextWindowInitState(state *TTextWindowState) {
	// WITH temp = state
	LineCount = 0
	LinePos = 1
	LoadedFilename = ""

}

func TextWindowDrawTitle(color int16, title TTextWindowLine) {
	VideoWriteText(TextWindowX+2, TextWindowY+1, color, TextWindowStrInnerEmpty)
	VideoWriteText(TextWindowX+((TextWindowWidth-Length(title))/2), TextWindowY+1, color, title)
}

func TextWindowDrawOpen(state *TTextWindowState) {
	var (
		ix, iy int16
	)
	// WITH temp = state
	for iy := 1; iy <= (TextWindowHeight + 1); iy++ {
		VideoMove(TextWindowX, iy+TextWindowY-1, TextWindowWidth, *ScreenCopy[iy], false)
	}
	for iy := (TextWindowHeight / 2); iy >= 0; iy-- {
		VideoWriteText(TextWindowX, TextWindowY+iy+1, 0x0F, TextWindowStrText)
		VideoWriteText(TextWindowX, TextWindowY+TextWindowHeight-iy-1, 0x0F, TextWindowStrText)
		VideoWriteText(TextWindowX, TextWindowY+iy, 0x0F, TextWindowStrTop)
		VideoWriteText(TextWindowX, TextWindowY+TextWindowHeight-iy, 0x0F, TextWindowStrBottom)
		Delay(25)
	}
	VideoWriteText(TextWindowX, TextWindowY+2, 0x0F, TextWindowStrSep)
	TextWindowDrawTitle(0x1E, Title)

}

func TextWindowDrawClose(state *TTextWindowState) {
	var (
		ix, iy     int16
		unk1, unk2 int16
	)
	// WITH temp = state
	for iy := 0; iy <= (TextWindowHeight / 2); iy++ {
		VideoWriteText(TextWindowX, TextWindowY+iy, 0x0F, TextWindowStrTop)
		VideoWriteText(TextWindowX, TextWindowY+TextWindowHeight-iy, 0x0F, TextWindowStrBottom)
		Delay(18)
		VideoMove(TextWindowX, TextWindowY+iy, TextWindowWidth, *ScreenCopy[iy+1], true)
		VideoMove(TextWindowX, TextWindowY+TextWindowHeight-iy, TextWindowWidth, *ScreenCopy[(TextWindowHeight-iy)+1], true)
	}

}

func TextWindowDrawLine(state *TTextWindowState, lpos int16, withoutFormatting, viewingFile bool) {
	var (
		lineY                        int16
		textOffset, textColor, textX int16
	)
	// WITH temp = state
	lineY = ((TextWindowY + lpos) - LinePos) + (TextWindowHeight / 2) + 1
	if lpos == LinePos {
		VideoWriteText(TextWindowX+2, lineY, 0x1C, TextWindowStrInnerArrows)
	} else {
		VideoWriteText(TextWindowX+2, lineY, 0x1E, TextWindowStrInnerEmpty)
	}
	if (lpos > 0) && (lpos <= LineCount) {
		if withoutFormatting {
			VideoWriteText(TextWindowX+4, lineY, 0x1E, Lines[lpos])
		} else {
			textOffset = 1
			textColor = 0x1E
			textX = TextWindowX + 4
			if Length(state.Lines[lpos]) > 0 {
				switch state.Lines[lpos][1] {
				case '!':
					textOffset = Pos(';', Lines[lpos]) + 1
					VideoWriteText(textX+2, lineY, 0x1D, '\x10')
					textX = textX + 5
					textColor = 0x1F
				case ':':
					textOffset = Pos(';', Lines[lpos]) + 1
					textColor = 0x1F
				case '$':
					textOffset = 2
					textColor = 0x1F
					textX = (textX - 4) + ((TextWindowWidth - Length(Lines[lpos])) / 2)
				}
			}
			if textOffset > 0 {
				VideoWriteText(textX, lineY, textColor, Copy(Lines[lpos], textOffset, Length(Lines[lpos])-textOffset+1))
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
	for i := 0; i <= (TextWindowHeight - 4); i++ {
		TextWindowDrawLine(state, state.LinePos-(TextWindowHeight/2)+i+2, withoutFormatting, viewingFile)
	}
	TextWindowDrawTitle(0x1E, state.Title)
}

func TextWindowAppend(state *TTextWindowState, line TTextWindowLine) {
	// WITH temp = state
	LineCount = LineCount + 1
	New(Lines[LineCount])
	Lines[LineCount] = line

}

func TextWindowFree(state *TTextWindowState) {
	// WITH temp = state
	for LineCount > 0 {
		Dispose(Lines[LineCount])
		LineCount = LineCount - 1
	}
	LoadedFilename = ""

}

func TextWindowPrint(state *TTextWindowState) {
	var (
		iLine, iChar int16
		line         string
	)
	// WITH temp = state
	Rewrite(Lst)
	for iLine := 1; iLine <= LineCount; iLine++ {
		line = Lines[iLine]
		if Length(line) > 0 {
			switch line[1] {
			case '$':
				Delete(line, 1, 1)
				for iChar := ((80 - Length(line)) / 2); iChar >= 1; iChar-- {
					line = ' ' + line
				}
			case '!', ':':
				iChar = Pos(';', line)
				if iChar > 0 {
					Delete(line, 1, iChar)
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
			exit()
		}
	}
	if LoadedFilename == "ORDER.HLP" {
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
	// WITH temp = state
	TextWindowRejected = false
	Hyperlink = ""
	TextWindowDraw(state, false, viewingFile)
	for {
		InputUpdate()
		newLinePos = LinePos
		if InputDeltaY != 0 {
			newLinePos = newLinePos + InputDeltaY
		} else if InputShiftPressed || (InputKeyPressed == KEY_ENTER) {
			InputShiftAccepted = true
			if (Lines[LinePos][1]) == '!' {
				pointerStr = Copy(Lines[LinePos], 2, Length(Lines[LinePos])-1)
				if Pos(';', pointerStr) > 0 {
					pointerStr = Copy(pointerStr, 1, Pos(';', pointerStr)-1)
				}
				if pointerStr[1] == '-' {
					Delete(pointerStr, 1, 1)
					TextWindowFree(state)
					TextWindowOpenFile(pointerStr, state)
					if state.LineCount == 0 {
						exit()
					} else {
						viewingFile = true
						newLinePos = LinePos
						TextWindowDraw(state, false, viewingFile)
						InputKeyPressed = '\x00'
						InputShiftPressed = false
					}
				} else {
					if hyperlinkAsSelect {
						Hyperlink = pointerStr
					} else {
						pointerStr = ':' + pointerStr
						for iLine := 1; iLine <= LineCount; iLine++ {
							if Length(pointerStr) > Length(Lines[iLine]) {
							} else {
								for iChar := 1; iChar <= Length(pointerStr); iChar++ {
									if UpCase(pointerStr[iChar]) != UpCase(Lines[iLine][iChar]) {
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
				newLinePos = LinePos - TextWindowHeight + 4
			} else if InputKeyPressed == KEY_PAGE_DOWN {
				newLinePos = LinePos + TextWindowHeight - 4
			} else if InputKeyPressed == KEY_ALT_P {
				TextWindowPrint(state)
			}

		}

	LabelMatched:
		if newLinePos < 1 {
			newLinePos = 1
		} else if newLinePos > state.LineCount {
			newLinePos = LineCount
		}

		if newLinePos != LinePos {
			LinePos = newLinePos
			TextWindowDraw(state, false, viewingFile)
			if (Lines[LinePos][1]) == '!' {
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
		var (
			i int16
		)
		// WITH temp = state
		if LineCount > 1 {
			Dispose(Lines[LinePos])
			for i := (LinePos + 1); i <= LineCount; i++ {
				Lines[i-1] = Lines[i]
			}
			LineCount = LineCount - 1
			if LinePos > LineCount {
				newLinePos = LineCount
			} else {
				TextWindowDraw(state, true, false)
			}
		} else {
			Lines[1] = ""
		}

	}

	// WITH temp = state
	if LineCount == 0 {
		TextWindowAppend(state, "")
	}
	insertMode = true
	LinePos = 1
	charPos = 1
	TextWindowDraw(state, true, false)
	for {
		if insertMode {
			VideoWriteText(77, 14, 0x1E, "on ")
		} else {
			VideoWriteText(77, 14, 0x1E, "off")
		}
		if charPos >= (Length(Lines[LinePos]) + 1) {
			charPos = Length(Lines[LinePos]) + 1
			VideoWriteText(charPos+TextWindowX+3, TextWindowY+(TextWindowHeight/2)+1, 0x70, ' ')
		} else {
			VideoWriteText(charPos+TextWindowX+3, TextWindowY+(TextWindowHeight/2)+1, 0x70, Lines[LinePos][charPos])
		}
		InputReadWaitKey()
		newLinePos = LinePos
		switch InputKeyPressed {
		case KEY_UP:
			newLinePos = LinePos - 1
		case KEY_DOWN:
			newLinePos = LinePos + 1
		case KEY_PAGE_UP:
			newLinePos = LinePos - TextWindowHeight + 4
		case KEY_PAGE_DOWN:
			newLinePos = LinePos + TextWindowHeight - 4
		case KEY_RIGHT:
			charPos = charPos + 1
			if charPos > (Length(Lines[LinePos]) + 1) {
				charPos = 1
				newLinePos = LinePos + 1
			}
		case KEY_LEFT:
			charPos = charPos - 1
			if charPos < 1 {
				charPos = TextWindowWidth
				newLinePos = LinePos - 1
			}
		case KEY_ENTER:
			if LineCount < MAX_TEXT_WINDOW_LINES {
				for i := LineCount; i >= (LinePos + 1); i-- {
					Lines[i+1] = Lines[i]
				}
				New(Lines[LinePos+1])
				Lines[LinePos+1] = Copy(Lines[LinePos], charPos, Length(Lines[LinePos])-charPos+1)
				Lines[LinePos] = Copy(Lines[LinePos], 1, charPos-1)
				newLinePos = LinePos + 1
				charPos = 1
				LineCount = LineCount + 1
			}
		case KEY_BACKSPACE:
			if charPos > 1 {
				Lines[LinePos] = Copy(Lines[LinePos], 1, charPos-2) + Copy(Lines[LinePos], charPos, Length(Lines[LinePos])-charPos+1)
				charPos = charPos - 1
			} else if Length(Lines[LinePos]) == 0 {
				DeleteCurrLine()
				newLinePos = LinePos - 1
				charPos = TextWindowWidth
			}

		case KEY_INSERT:
			insertMode = !insertMode
		case KEY_DELETE:
			Lines[LinePos] = Copy(Lines[LinePos], 1, charPos-1) + Copy(Lines[LinePos], charPos+1, Length(Lines[LinePos])-charPos)
		case KEY_CTRL_Y:
			DeleteCurrLine()
		default:
			if (InputKeyPressed >= ' ') && (charPos < (TextWindowWidth - 7)) {
				if !insertMode {
					Lines[LinePos] = Copy(Lines[LinePos], 1, charPos-1) + InputKeyPressed + Copy(Lines[LinePos], charPos+1, Length(Lines[LinePos])-charPos)
					charPos = charPos + 1
				} else {
					if Length(Lines[LinePos]) < (TextWindowWidth - 8) {
						Lines[LinePos] = Copy(Lines[LinePos], 1, charPos-1) + InputKeyPressed + Copy(Lines[LinePos], charPos, Length(Lines[LinePos])-charPos+1)
						charPos = charPos + 1
					}
				}
			}
		}
		if newLinePos < 1 {
			newLinePos = 1
		} else if newLinePos > LineCount {
			newLinePos = LineCount
		}

		if newLinePos != LinePos {
			LinePos = newLinePos
			TextWindowDraw(state, true, false)
		} else {
			TextWindowDrawLine(state, LinePos, true, false)
		}
		if InputKeyPressed == KEY_ESCAPE {
			break
		}
	}
	if Length(Lines[LineCount]) == 0 {
		Dispose(Lines[LineCount])
		LineCount = LineCount - 1
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
	// WITH temp = state
	retVal = true
	for i := 1; i <= Length(filename); i++ {
		retVal = retVal && (filename[i] != '.')
	}
	if retVal {
		filename = filename + ".HLP"
	}
	if filename[1] == '*' {
		filename = Copy(filename, 2, Length(filename)-1)
		entryPos = -1
	} else {
		entryPos = 0
	}
	TextWindowInitState(state)
	LoadedFilename = UpCaseString(filename)
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
		for i := 1; i <= ResourceDataHeader.EntryCount; i++ {
			if UpCaseString(ResourceDataHeader.Name[i]) == UpCaseString(filename) {
				entryPos = i
			}
		}
	}
	if entryPos <= 0 {
		Assign(tf, filename)
		Reset(tf)
		for (IOResult == 0) && (!Eof(tf)) {
			Inc(LineCount)
			New(Lines[LineCount])
			ReadLn(tf, Lines[LineCount])
		}
		Close(tf)
	} else {
		Assign(f, ResourceDataFilename)
		Reset(f, 1)
		Seek(f, ResourceDataHeader.FileOffset[entryPos])
		if IOResult == 0 {
			retVal = true
			for (IOResult == 0) && retVal {
				Inc(LineCount)
				New(Lines[LineCount])
				BlockRead(f, Lines[LineCount], 1)
				line = Ptr(Seg(Lines[LineCount]), Ofs(Lines[LineCount])+1)
				lineLen = Ord(Lines[LineCount][0])
				if lineLen == 0 {
					Lines[LineCount] = ""
				} else {
					BlockRead(f, line, Ord(Lines[LineCount][0]))
				}
				if Lines[LineCount] == '@' {
					retVal = false
					Lines[LineCount] = ""
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
	// WITH temp = state
	Assign(f, filename)
	Rewrite(f)
	if IOResult != 0 {
		exit()
	}
	for i := 1; i <= LineCount; i++ {
		WriteLn(f, Lines[i])
		if IOResult != 0 {
			exit()
		}
	}
	Close(f)

}

func TextWindowDisplayFile(filename, title string) {
	var (
		state TTextWindowState
	)
	state.Title = title
	TextWindowOpenFile(filename, state)
	state.Selectable = false
	if state.LineCount > 0 {
		TextWindowDrawOpen(state)
		TextWindowSelect(state, false, true)
		TextWindowDrawClose(state)
	}
	TextWindowFree(state)
}

func TextWindowInit(x, y, width, height int16) {
	var (
		i int16
	)
	TextWindowX = x
	TextWindowWidth = width
	TextWindowY = y
	TextWindowHeight = height
	TextWindowStrInnerEmpty = ""
	TextWindowStrInnerLine = ""
	for i := 1; i <= (TextWindowWidth - 5); i++ {
		TextWindowStrInnerEmpty = TextWindowStrInnerEmpty + ' '
		TextWindowStrInnerLine = TextWindowStrInnerLine + 'Í'
	}
	TextWindowStrTop = "\xc6\xd1" + TextWindowStrInnerLine + 'Ñ' + 'µ'
	TextWindowStrBottom = "\xc6\xcf" + TextWindowStrInnerLine + 'Ï' + 'µ'
	TextWindowStrSep = " \xc6" + TextWindowStrInnerLine + 'µ' + ' '
	TextWindowStrText = " \xb3" + TextWindowStrInnerEmpty + '³' + ' '
	TextWindowStrInnerArrows = TextWindowStrInnerEmpty
	TextWindowStrInnerArrows[1] = '¯'
	TextWindowStrInnerArrows[Length(TextWindowStrInnerArrows)] = '®'
	TextWindowStrInnerSep = TextWindowStrInnerEmpty
	for i := 1; i <= (TextWindowWidth / 5); i++ {
		TextWindowStrInnerSep[i*5+((TextWindowWidth%5)/2)] = '\a'
	}
}

func init() {
	ResourceDataFileName = ""
	ResourceDataHeader.EntryCount = 0
}
