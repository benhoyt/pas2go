package main // unit: Editor

// interface uses: GameVars, TxtWind

// implementation uses: Dos, Crt, Video, Sounds, Input, Elements, Oop, Game

type TDrawMode uint8

const (
	DrawingOff TDrawMode = iota + 1
	DrawingOn
	TextEntry
)

const NeighborBoardStrs [4]string = [...]string{"       Board \x18", "       Board \x19", "       Board \x1b", "       Board \x1a"}

func EditorAppendBoard() {
	if World.BoardCount < MAX_BOARD {
		BoardClose()
		World.BoardCount = World.BoardCount + 1
		World.Info.CurrentBoard = World.BoardCount
		World.BoardLen[World.BoardCount] = 0
		BoardCreate()
		TransitionDrawToBoard()
		for {
			PopupPromptString("Room's Title:", Board.Name)
			if Length(Board.Name) != 0 {
				break
			}
		}
		TransitionDrawToBoard()
	}
}

func EditorLoop() {
	var (
		selectedCategory           int16
		elemMenuColor              int16
		wasModified                bool
		editorExitRequested        bool
		drawMode                   TDrawMode
		cursorX, cursorY           int16
		cursorPattern, cursorColor int16
		i, iElem                   int16
		canModify                  bool
		unk1                       [50]byte
		copiedStat                 TStat
		copiedHasStat              bool
		copiedTile                 TTile
		copiedX, copiedY           int16
		cursorBlinker              int16
	)
	EditorDrawSidebar := func() {
		var (
			i         int16
			copiedChr byte
		)
		SidebarClear()
		SidebarClearLine(1)
		VideoWriteText(61, 0, 0x1F, "     - - - -       ")
		VideoWriteText(62, 1, 0x70, "  ZZT Editor   ")
		VideoWriteText(61, 2, 0x1F, "     - - - -       ")
		VideoWriteText(61, 4, 0x70, " L ")
		VideoWriteText(64, 4, 0x1F, " Load")
		VideoWriteText(61, 5, 0x30, " S ")
		VideoWriteText(64, 5, 0x1F, " Save")
		VideoWriteText(70, 4, 0x70, " H ")
		VideoWriteText(73, 4, 0x1E, " Help")
		VideoWriteText(70, 5, 0x30, " Q ")
		VideoWriteText(73, 5, 0x1F, " Quit")
		VideoWriteText(61, 7, 0x70, " B ")
		VideoWriteText(65, 7, 0x1F, " Switch boards")
		VideoWriteText(61, 8, 0x30, " I ")
		VideoWriteText(65, 8, 0x1F, " Board Info")
		VideoWriteText(61, 10, 0x70, "  f1   ")
		VideoWriteText(68, 10, 0x1F, " Item")
		VideoWriteText(61, 11, 0x30, "  f2   ")
		VideoWriteText(68, 11, 0x1F, " Creature")
		VideoWriteText(61, 12, 0x70, "  f3   ")
		VideoWriteText(68, 12, 0x1F, " Terrain")
		VideoWriteText(61, 13, 0x30, "  f4   ")
		VideoWriteText(68, 13, 0x1F, " Enter text")
		VideoWriteText(61, 15, 0x70, " Space ")
		VideoWriteText(68, 15, 0x1F, " Plot")
		VideoWriteText(61, 16, 0x30, "  Tab  ")
		VideoWriteText(68, 16, 0x1F, " Draw mode")
		VideoWriteText(61, 18, 0x70, " P ")
		VideoWriteText(64, 18, 0x1F, " Pattern")
		VideoWriteText(61, 19, 0x30, " C ")
		VideoWriteText(64, 19, 0x1F, " Color:")
		for i = 9; i <= 15; i++ {
			VideoWriteText(61+i, 22, i, 'Û')
		}
		for i = 1; i <= EditorPatternCount; i++ {
			VideoWriteText(61+i, 22, 0x0F, ElementDefs[EditorPatterns[i+1]].Character)
		}
		if ElementDefs[copiedTile.Element].HasDrawProc {
			ElementDefs[copiedTile.Element].DrawProc(copiedX, copiedY, copiedChr)
		} else {
			copiedChr = Ord(ElementDefs[copiedTile.Element].Character)
		}
		VideoWriteText(62+EditorPatternCount, 22, copiedTile.Color, Chr(copiedChr))
		VideoWriteText(61, 24, 0x1F, " Mode:")
	}

	EditorDrawTileAndNeighborsAt := func(x, y int16) {
		var i, ix, iy int16
		BoardDrawTile(x, y)
		for i = 0; i <= 3; i++ {
			ix = x + NeighborDeltaX[i]
			iy = y + NeighborDeltaY[i]
			if (ix >= 1) && (ix <= BOARD_WIDTH) && (iy >= 1) && (iy <= BOARD_HEIGHT) {
				BoardDrawTile(ix, iy)
			}
		}
	}

	EditorUpdateSidebar := func() {
		if drawMode == DrawingOn {
			VideoWriteText(68, 24, 0x9E, "Drawing on ")
		} else if drawMode == TextEntry {
			VideoWriteText(68, 24, 0x9E, "Text entry ")
		} else if drawMode == DrawingOff {
			VideoWriteText(68, 24, 0x1E, "Drawing off")
		}

		VideoWriteText(72, 19, 0x1E, ColorNames[(cursorColor-8)+1])
		VideoWriteText(61+cursorPattern, 21, 0x1F, '\x1f')
		VideoWriteText(61+cursorColor, 21, 0x1F, '\x1f')
	}

	EditorDrawRefresh := func() {
		var boardNumStr string
		BoardDrawBorder()
		EditorDrawSidebar()
		boardNumStr = fmt.Sprint(World.Info.CurrentBoard)
		TransitionDrawToBoard()
		if Length(Board.Name) != 0 {
			VideoWriteText((59-Length(Board.Name))/2, 0, 0x70, ' '+Board.Name+' ')
		} else {
			VideoWriteText(26, 0, 0x70, " Untitled ")
		}
	}

	EditorSetAndCopyTile := func(x, y, element, color byte) {
		Board.Tiles[x][y].Element = element
		Board.Tiles[x][y].Color = color
		copiedTile = Board.Tiles[x][y]
		copiedHasStat = false
		copiedX = x
		copiedY = y
		EditorDrawTileAndNeighborsAt(x, y)
	}

	EditorAskSaveChanged := func() {
		InputKeyPressed = '\x00'
		if wasModified {
			if SidebarPromptYesNo("Save first? ", true) {
				if InputKeyPressed != KEY_ESCAPE {
					GameWorldSave("Save world", LoadedGameFileName, ".ZZT")
				}
			}
		}
		World.Info.Name = LoadedGameFileName
	}

	EditorPrepareModifyTile := func(x, y int16) (EditorPrepareModifyTile bool) {
		wasModified = true
		EditorPrepareModifyTile = BoardPrepareTileForPlacement(x, y)
		EditorDrawTileAndNeighborsAt(x, y)
		return
	}

	EditorPrepareModifyStatAtCursor := func() (EditorPrepareModifyStatAtCursor bool) {
		if Board.StatCount < MAX_STAT {
			EditorPrepareModifyStatAtCursor = EditorPrepareModifyTile(cursorX, cursorY)
		} else {
			EditorPrepareModifyStatAtCursor = false
		}
		return
	}

	EditorPlaceTile := func(x, y int16) {
		tile := &Board.Tiles[x][y]
		if cursorPattern <= EditorPatternCount {
			if EditorPrepareModifyTile(x, y) {
				tile.Element = EditorPatterns[cursorPattern+1]
				tile.Color = cursorColor
			}
		} else if copiedHasStat {
			if EditorPrepareModifyStatAtCursor {
				AddStat(x, y, copiedTile.Element, copiedTile.Color, copiedStat.Cycle, copiedStat)
			}
		} else {
			if EditorPrepareModifyTile(x, y) {
				Board.Tiles[x][y] = copiedTile
			}
		}

		EditorDrawTileAndNeighborsAt(x, y)

	}

	EditorEditBoardInfo := func() {
		var (
			state         TTextWindowState
			i             int16
			numStr        string
			exitRequested bool
		)
		BoolToString := func(val bool) (BoolToString string) {
			if val {
				BoolToString = "Yes"
			} else {
				BoolToString = "No "
			}
			return
		}

		state.Title = "Board Information"
		TextWindowDrawOpen(state)
		state.LinePos = 1
		state.LineCount = 9
		state.Selectable = true
		exitRequested = false
		for i = 1; i <= state.LineCount; i++ {
			New(state.Lines[i+1])
		}
		for {
			state.Selectable = true
			state.LineCount = 10
			for i = 1; i <= state.LineCount; i++ {
				New(state.Lines[i+1])
			}
			state.Lines[2] = "         Title: " + Board.Name
			numStr = fmt.Sprint(Board.Info.MaxShots)
			state.Lines[3] = "      Can fire: " + numStr + " shots."
			state.Lines[4] = " Board is dark: " + BoolToString(Board.Info.IsDark)
			for i = 4; i <= 7; i++ {
				state.Lines[i+1] = NeighborBoardStrs[i-4] + ": " + EditorGetBoardName(Board.Info.NeighborBoards[i-4], true)
			}
			state.Lines[9] = "Re-enter when zapped: " + BoolToString(Board.Info.ReenterWhenZapped)
			numStr = fmt.Sprint(Board.Info.TimeLimitSec)
			state.Lines[10] = "  Time limit, 0=None: " + numStr + " sec."
			state.Lines[11] = "          Quit!"
			TextWindowSelect(state, false, false)
			if (InputKeyPressed == KEY_ENTER) && (state.LinePos >= 1) && (state.LinePos <= 8) {
				wasModified = true
			}
			if InputKeyPressed == KEY_ENTER {
				switch state.LinePos {
				case 1:
					PopupPromptString("New title for board:", Board.Name)
					exitRequested = true
					TextWindowDrawClose(state)
				case 2:
					numStr = fmt.Sprint(Board.Info.MaxShots)
					SidebarPromptString("Maximum shots?", "", numStr, PROMPT_NUMERIC)
					if Length(numStr) != 0 {
						Val(numStr, Board.Info.MaxShots, i)
					}
					EditorDrawSidebar()
				case 3:
					Board.Info.IsDark = !Board.Info.IsDark
				case 4, 5, 6, 7:
					Board.Info.NeighborBoards[state.LinePos-4] = EditorSelectBoard(NeighborBoardStrs[state.LinePos-4], Board.Info.NeighborBoards[state.LinePos-4], true)
					if Board.Info.NeighborBoards[state.LinePos-4] > World.BoardCount {
						EditorAppendBoard()
					}
					exitRequested = true
				case 8:
					Board.Info.ReenterWhenZapped = !Board.Info.ReenterWhenZapped
				case 9:
					numStr = fmt.Sprint(Board.Info.TimeLimitSec)
					SidebarPromptString("Time limit?", " Sec", numStr, PROMPT_NUMERIC)
					if Length(numStr) != 0 {
						Val(numStr, Board.Info.TimeLimitSec, i)
					}
					EditorDrawSidebar()
				case 10:
					exitRequested = true
					TextWindowDrawClose(state)
				}
			} else {
				exitRequested = true
				TextWindowDrawClose(state)
			}
			if exitRequested {
				break
			}
		}
		TextWindowFree(state)
	}

	EditorEditStatText := func(statId int16, prompt string) {
		var (
			state        TTextWindowState
			iLine, iChar int16
			unk1         [52]byte
			dataChar     byte
			dataPtr      uintptr
		)
		stat := &Board.Stats[statId]
		state.Title = prompt
		TextWindowDrawOpen(state)
		state.Selectable = false
		CopyStatDataToTextWindow(statId, state)
		if stat.DataLen > 0 {
			FreeMem(stat.Data, stat.DataLen)
			stat.DataLen = 0
		}
		EditorOpenEditTextWindow(state)
		for iLine = 1; iLine <= state.LineCount; iLine++ {
			stat.DataLen = stat.DataLen + Length(state.Lines[iLine+1]) + 1
		}
		GetMem(stat.Data, stat.DataLen)
		dataPtr = stat.Data
		for iLine = 1; iLine <= state.LineCount; iLine++ {
			for iChar = 1; iChar <= Length(state.Lines[iLine+1]); iChar++ {
				dataChar = state.Lines[iLine+1][iChar]
				Move(dataChar, dataPtr, 1)
				AdvancePointer(dataPtr, 1)
			}
			dataChar = '\r'
			Move(dataChar, dataPtr, 1)
			AdvancePointer(dataPtr, 1)
		}
		TextWindowFree(state)
		TextWindowDrawClose(state)
		InputKeyPressed = '\x00'

	}

	EditorEditStat := func(statId int16) {
		var (
			element       byte
			i             int16
			categoryName  string
			selectedBoard byte
			iy            int16
			promptByte    byte
		)
		EditorEditStatSettings := func(selected bool) {
			stat := &Board.Stats[statId]
			InputKeyPressed = '\x00'
			iy = 9
			if Length(ElementDefs[element].Param1Name) != 0 {
				if Length(ElementDefs[element].ParamTextName) == 0 {
					SidebarPromptSlider(selected, 63, iy, ElementDefs[element].Param1Name, stat.P1)
				} else {
					if stat.P1 == 0 {
						stat.P1 = World.EditorStatSettings[element].P1
					}
					BoardDrawTile(stat.X, stat.Y)
					SidebarPromptCharacter(selected, 63, iy, ElementDefs[element].Param1Name, stat.P1)
					BoardDrawTile(stat.X, stat.Y)
				}
				if selected {
					World.EditorStatSettings[element].P1 = stat.P1
				}
				iy = iy + 4
			}
			if (InputKeyPressed != KEY_ESCAPE) && (Length(ElementDefs[element].ParamTextName) != 0) {
				if selected {
					EditorEditStatText(statId, ElementDefs[element].ParamTextName)
				}
			}
			if (InputKeyPressed != KEY_ESCAPE) && (Length(ElementDefs[element].Param2Name) != 0) {
				promptByte = (stat.P2 % 0x80)
				SidebarPromptSlider(selected, 63, iy, ElementDefs[element].Param2Name, promptByte)
				if selected {
					stat.P2 = (stat.P2 && 0x80) + promptByte
					World.EditorStatSettings[element].P2 = stat.P2
				}
				iy = iy + 4
			}
			if (InputKeyPressed != KEY_ESCAPE) && (Length(ElementDefs[element].ParamBulletTypeName) != 0) {
				promptByte = (stat.P2) / 0x80
				SidebarPromptChoice(selected, iy, ElementDefs[element].ParamBulletTypeName, "Bullets Stars", promptByte)
				if selected {
					stat.P2 = (stat.P2 % 0x80) + (promptByte * 0x80)
					World.EditorStatSettings[element].P2 = stat.P2
				}
				iy = iy + 4
			}
			if (InputKeyPressed != KEY_ESCAPE) && (Length(ElementDefs[element].ParamDirName) != 0) {
				SidebarPromptDirection(selected, iy, ElementDefs[element].ParamDirName, stat.StepX, stat.StepY)
				if selected {
					World.EditorStatSettings[element].StepX = stat.StepX
					World.EditorStatSettings[element].StepY = stat.StepY
				}
				iy = iy + 4
			}
			if (InputKeyPressed != KEY_ESCAPE) && (Length(ElementDefs[element].ParamBoardName) != 0) {
				if selected {
					selectedBoard = EditorSelectBoard(ElementDefs[element].ParamBoardName, stat.P3, true)
					if selectedBoard != 0 {
						stat.P3 = selectedBoard
						World.EditorStatSettings[element].P3 = World.Info.CurrentBoard
						if stat.P3 > World.BoardCount {
							EditorAppendBoard()
							copiedHasStat = false
							copiedTile.Element = 0
							copiedTile.Color = 0x0F
						}
						World.EditorStatSettings[element].P3 = stat.P3
					} else {
						InputKeyPressed = KEY_ESCAPE
					}
					iy = iy + 4
				} else {
					VideoWriteText(63, iy, 0x1F, "Room: "+Copy(EditorGetBoardName(stat.P3, true), 1, 10))
				}
			}

		}

		stat := &Board.Stats[statId]
		SidebarClear()
		element = Board.Tiles[stat.X][stat.Y].Element
		wasModified = true
		categoryName = ""
		for i = 0; i <= element; i++ {
			if (ElementDefs[i].EditorCategory == ElementDefs[element].EditorCategory) && (Length(ElementDefs[i].CategoryName) != 0) {
				categoryName = ElementDefs[i].CategoryName
			}
		}
		VideoWriteText(64, 6, 0x1E, categoryName)
		VideoWriteText(64, 7, 0x1F, ElementDefs[element].Name)
		EditorEditStatSettings(false)
		EditorEditStatSettings(true)
		if InputKeyPressed != KEY_ESCAPE {
			copiedHasStat = true
			copiedStat = Board.Stats[statId]
			copiedTile = Board.Tiles[stat.X][stat.Y]
			copiedX = stat.X
			copiedY = stat.Y
		}

	}

	EditorTransferBoard := func() {
		var (
			i byte
			f FILE
		)
		i = 1
		SidebarPromptChoice(true, 3, "Transfer board:", "Import Export", i)
		if InputKeyPressed != KEY_ESCAPE {
			if i == 0 {
				SidebarPromptString("Import board", ".BRD", SavedBoardFileName, PROMPT_ALPHANUM)
				if (InputKeyPressed != KEY_ESCAPE) && (Length(SavedBoardFileName) != 0) {
					Assign(f, SavedBoardFileName+".BRD")
					Reset(f, 1)
					if DisplayIOError {
						goto TransferEnd
					}
					BoardClose()
					FreeMem(World.BoardData[World.Info.CurrentBoard], World.BoardLen[World.Info.CurrentBoard])
					BlockRead(f, World.BoardLen[World.Info.CurrentBoard], 2)
					if !DisplayIOError {
						GetMem(World.BoardData[World.Info.CurrentBoard], World.BoardLen[World.Info.CurrentBoard])
						BlockRead(f, World.BoardData[World.Info.CurrentBoard], World.BoardLen[World.Info.CurrentBoard])
					}
					if DisplayIOError {
						World.BoardLen[World.Info.CurrentBoard] = 0
						BoardCreate()
						EditorDrawRefresh()
					} else {
						BoardOpen(World.Info.CurrentBoard)
						EditorDrawRefresh()
						for i = 0; i <= 3; i++ {
							Board.Info.NeighborBoards[i] = 0
						}
					}
				}
			} else if i == 1 {
				SidebarPromptString("Export board", ".BRD", SavedBoardFileName, PROMPT_ALPHANUM)
				if (InputKeyPressed != KEY_ESCAPE) && (Length(SavedBoardFileName) != 0) {
					Assign(f, SavedBoardFileName+".BRD")
					Rewrite(f, 1)
					if DisplayIOError {
						goto TransferEnd
					}
					BoardClose()
					BlockWrite(f, World.BoardLen[World.Info.CurrentBoard], 2)
					BlockWrite(f, World.BoardData[World.Info.CurrentBoard], World.BoardLen[World.Info.CurrentBoard])
					BoardOpen(World.Info.CurrentBoard)
					if DisplayIOError {
					} else {
						Close(f)
					}
				}
			}

		}
	TransferEnd:
		EditorDrawSidebar()

	}

	EditorFloodFill := func(x, y int16, from TTile) {
		var (
			i              int16
			tileAt         TTile
			toFill, filled byte
			xPosition      [256]int16
			yPosition      [256]int16
		)
		toFill = 1
		filled = 0
		for toFill != filled {
			tileAt = Board.Tiles[x][y]
			EditorPlaceTile(x, y)
			if (Board.Tiles[x][y].Element != tileAt.Element) || (Board.Tiles[x][y].Color != tileAt.Color) {
				for i = 0; i <= 3; i++ {
					tile := &Board.Tiles[x+NeighborDeltaX[i]][y+NeighborDeltaY[i]]
					if (tile.Element == from.Element) && ((from.Element == 0) || (tile.Color == from.Color)) {
						xPosition[toFill] = x + NeighborDeltaX[i]
						yPosition[toFill] = y + NeighborDeltaY[i]
						toFill = toFill + 1
					}

				}
			}
			filled = filled + 1
			x = xPosition[filled]
			y = yPosition[filled]
		}
	}

	if World.Info.IsSave || (WorldGetFlagPosition("SECRET") >= 0) {
		WorldUnload()
		WorldCreate()
	}
	InitElementsEditor()
	CurrentTick = 0
	wasModified = false
	cursorX = 30
	cursorY = 12
	drawMode = DrawingOff
	cursorPattern = 1
	cursorColor = 0x0E
	cursorBlinker = 0
	copiedHasStat = false
	copiedTile.Element = 0
	copiedTile.Color = 0x0F
	if World.Info.CurrentBoard != 0 {
		BoardChange(World.Info.CurrentBoard)
	}
	EditorDrawRefresh()
	if World.BoardCount == 0 {
		EditorAppendBoard()
	}
	editorExitRequested = false
	for {
		if drawMode == DrawingOn {
			EditorPlaceTile(cursorX, cursorY)
		}
		InputUpdate()
		if (InputKeyPressed == '\x00') && (InputDeltaX == 0) && (InputDeltaY == 0) && !InputShiftPressed {
			if SoundHasTimeElapsed(TickTimeCounter, 15) {
				cursorBlinker = (cursorBlinker + 1) % 3
			}
			if cursorBlinker == 0 {
				BoardDrawTile(cursorX, cursorY)
			} else {
				VideoWriteText(cursorX-1, cursorY-1, 0x0F, 'Å')
			}
			EditorUpdateSidebar()
		} else {
			BoardDrawTile(cursorX, cursorY)
		}
		if drawMode == TextEntry {
			if (InputKeyPressed >= ' ') && (InputKeyPressed < '\u0080') {
				if EditorPrepareModifyTile(cursorX, cursorY) {
					Board.Tiles[cursorX][cursorY].Element = (cursorColor - 9) + E_TEXT_MIN
					Board.Tiles[cursorX][cursorY].Color = Ord(InputKeyPressed)
					EditorDrawTileAndNeighborsAt(cursorX, cursorY)
					InputDeltaX = 1
					InputDeltaY = 0
				}
				InputKeyPressed = '\x00'
			} else if (InputKeyPressed == KEY_BACKSPACE) && (cursorX > 1) && EditorPrepareModifyTile(cursorX-1, cursorY) {
				cursorX = cursorX - 1
			} else if (InputKeyPressed == KEY_ENTER) || (InputKeyPressed == KEY_ESCAPE) {
				drawMode = DrawingOff
				InputKeyPressed = '\x00'
			}

		}
		tile := &Board.Tiles[cursorX][cursorY]
		if InputShiftPressed || (InputKeyPressed == ' ') {
			InputShiftAccepted = true
			if (tile.Element == 0) || (ElementDefs[tile.Element].PlaceableOnTop && copiedHasStat && (cursorPattern > EditorPatternCount)) || (InputDeltaX != 0) || (InputDeltaY != 0) {
				EditorPlaceTile(cursorX, cursorY)
			} else {
				canModify = EditorPrepareModifyTile(cursorX, cursorY)
				if canModify {
					Board.Tiles[cursorX][cursorY].Element = 0
				}
			}
		}
		if (InputDeltaX != 0) || (InputDeltaY != 0) {
			cursorX = cursorX + InputDeltaX
			if cursorX < 1 {
				cursorX = 1
			}
			if cursorX > BOARD_WIDTH {
				cursorX = BOARD_WIDTH
			}
			cursorY = cursorY + InputDeltaY
			if cursorY < 1 {
				cursorY = 1
			}
			if cursorY > BOARD_HEIGHT {
				cursorY = BOARD_HEIGHT
			}
			VideoWriteText(cursorX-1, cursorY-1, 0x0F, 'Å')
			if (InputKeyPressed == '\x00') && InputJoystickEnabled {
				Delay(70)
			}
			InputShiftAccepted = false
		}
		switch UpCase(InputKeyPressed) {
		case '`':
			EditorDrawRefresh()
		case 'P':
			VideoWriteText(62, 21, 0x1F, "       ")
			if cursorPattern <= EditorPatternCount {
				cursorPattern = cursorPattern + 1
			} else {
				cursorPattern = 1
			}
		case 'C':
			VideoWriteText(72, 19, 0x1E, "       ")
			VideoWriteText(69, 21, 0x1F, "        ")
			if (cursorColor % 0x10) != 0x0F {
				cursorColor = cursorColor + 1
			} else {
				cursorColor = ((cursorColor / 0x10) * 0x10) + 9
			}
		case 'L':
			EditorAskSaveChanged()
			if (InputKeyPressed != KEY_ESCAPE) && GameWorldLoad(".ZZT") {
				if World.Info.IsSave || (WorldGetFlagPosition("SECRET") >= 0) {
					if !DebugEnabled {
						SidebarClearLine(3)
						SidebarClearLine(4)
						SidebarClearLine(5)
						VideoWriteText(63, 4, 0x1E, "Can not edit")
						if World.Info.IsSave {
							VideoWriteText(63, 5, 0x1E, "a saved game!")
						} else {
							VideoWriteText(63, 5, 0x1E, "  "+World.Info.Name+'!')
						}
						PauseOnError()
						WorldUnload()
						WorldCreate()
					}
				}
				wasModified = false
				EditorDrawRefresh()
			}
			EditorDrawSidebar()
		case 'S':
			GameWorldSave("Save world:", LoadedGameFileName, ".ZZT")
			if InputKeyPressed != KEY_ESCAPE {
				wasModified = false
			}
			EditorDrawSidebar()
		case 'Z':
			if SidebarPromptYesNo("Clear board? ", false) {
				for i = Board.StatCount; i >= 1; i-- {
					RemoveStat(i)
				}
				BoardCreate()
				EditorDrawRefresh()
			} else {
				EditorDrawSidebar()
			}
		case 'N':
			if SidebarPromptYesNo("Make new world? ", false) && (InputKeyPressed != KEY_ESCAPE) {
				EditorAskSaveChanged()
				if InputKeyPressed != KEY_ESCAPE {
					WorldUnload()
					WorldCreate()
					EditorDrawRefresh()
					wasModified = false
				}
			}
			EditorDrawSidebar()
		case 'Q', KEY_ESCAPE:
			editorExitRequested = true
		case 'B':
			i = EditorSelectBoard("Switch boards", World.Info.CurrentBoard, false)
			if InputKeyPressed != KEY_ESCAPE {
				if i > World.BoardCount {
					if SidebarPromptYesNo("Add new board? ", false) {
						EditorAppendBoard()
					}
				}
				BoardChange(i)
				EditorDrawRefresh()
			}
			EditorDrawSidebar()
		case '?':
			GameDebugPrompt()
			EditorDrawSidebar()
		case KEY_TAB:
			if drawMode == DrawingOff {
				drawMode = DrawingOn
			} else {
				drawMode = DrawingOff
			}
		case KEY_F1, KEY_F2, KEY_F3:
			VideoWriteText(cursorX-1, cursorY-1, 0x0F, 'Å')
			for i = 3; i <= 20; i++ {
				SidebarClearLine(i)
			}
			switch InputKeyPressed {
			case KEY_F1:
				selectedCategory = CATEGORY_ITEM
			case KEY_F2:
				selectedCategory = CATEGORY_CREATURE
			case KEY_F3:
				selectedCategory = CATEGORY_TERRAIN
			}
			i = 3
			for iElem = 0; iElem <= MAX_ELEMENT; iElem++ {
				if ElementDefs[iElem].EditorCategory == selectedCategory {
					if Length(ElementDefs[iElem].CategoryName) != 0 {
						i = i + 1
						VideoWriteText(65, i, 0x1E, ElementDefs[iElem].CategoryName)
						i = i + 1
					}
					VideoWriteText(61, i, ((i%2)<<6)+0x30, ' '+ElementDefs[iElem].EditorShortcut+' ')
					VideoWriteText(65, i, 0x1F, ElementDefs[iElem].Name)
					if ElementDefs[iElem].Color == COLOR_CHOICE_ON_BLACK {
						elemMenuColor = (cursorColor % 0x10) + 0x10
					} else if ElementDefs[iElem].Color == COLOR_WHITE_ON_CHOICE {
						elemMenuColor = (cursorColor * 0x10) - 0x71
					} else if ElementDefs[iElem].Color == COLOR_CHOICE_ON_CHOICE {
						elemMenuColor = ((cursorColor - 8) * 0x11) + 8
					} else if (ElementDefs[iElem].Color && 0x70) == 0x00 {
						elemMenuColor = (ElementDefs[iElem].Color % 0x10) + 0x10
					} else {
						elemMenuColor = ElementDefs[iElem].Color
					}

					VideoWriteText(78, i, elemMenuColor, ElementDefs[iElem].Character)
					i = i + 1
				}
			}
			InputReadWaitKey()
			for iElem = 1; iElem <= MAX_ELEMENT; iElem++ {
				if (ElementDefs[iElem].EditorCategory == selectedCategory) && (ElementDefs[iElem].EditorShortcut == UpCase(InputKeyPressed)) {
					if iElem == E_PLAYER {
						if EditorPrepareModifyTile(cursorX, cursorY) {
							MoveStat(0, cursorX, cursorY)
						}
					} else {
						if ElementDefs[iElem].Color == COLOR_CHOICE_ON_BLACK {
							elemMenuColor = cursorColor
						} else if ElementDefs[iElem].Color == COLOR_WHITE_ON_CHOICE {
							elemMenuColor = (cursorColor * 0x10) - 0x71
						} else if ElementDefs[iElem].Color == COLOR_CHOICE_ON_CHOICE {
							elemMenuColor = ((cursorColor - 8) * 0x11) + 8
						} else {
							elemMenuColor = ElementDefs[iElem].Color
						}

						if ElementDefs[iElem].Cycle == -1 {
							if EditorPrepareModifyTile(cursorX, cursorY) {
								EditorSetAndCopyTile(cursorX, cursorY, iElem, elemMenuColor)
							}
						} else {
							if EditorPrepareModifyStatAtCursor {
								AddStat(cursorX, cursorY, iElem, elemMenuColor, ElementDefs[iElem].Cycle, StatTemplateDefault)
								stat := &Board.Stats[Board.StatCount]
								if Length(ElementDefs[iElem].Param1Name) != 0 {
									stat.P1 = World.EditorStatSettings[iElem].P1
								}
								if Length(ElementDefs[iElem].Param2Name) != 0 {
									stat.P2 = World.EditorStatSettings[iElem].P2
								}
								if Length(ElementDefs[iElem].ParamDirName) != 0 {
									stat.StepX = World.EditorStatSettings[iElem].StepX
									stat.StepY = World.EditorStatSettings[iElem].StepY
								}
								if Length(ElementDefs[iElem].ParamBoardName) != 0 {
									stat.P3 = World.EditorStatSettings[iElem].P3
								}

								EditorEditStat(Board.StatCount)
								if InputKeyPressed == KEY_ESCAPE {
									RemoveStat(Board.StatCount)
								}
							}
						}
					}
				}
			}
			EditorDrawSidebar()
		case KEY_F4:
			if drawMode != TextEntry {
				drawMode = TextEntry
			} else {
				drawMode = DrawingOff
			}
		case 'H':
			TextWindowDisplayFile("editor.hlp", "World editor help")
		case 'X':
			EditorFloodFill(cursorX, cursorY, Board.Tiles[cursorX][cursorY])
		case '!':
			EditorEditHelpFile()
			EditorDrawSidebar()
		case 'T':
			EditorTransferBoard()
		case KEY_ENTER:
			if GetStatIdAt(cursorX, cursorY) >= 0 {
				EditorEditStat(GetStatIdAt(cursorX, cursorY))
				EditorDrawSidebar()
			} else {
				copiedHasStat = false
				copiedTile = Board.Tiles[cursorX][cursorY]
			}
		case 'I':
			EditorEditBoardInfo()
			TransitionDrawToBoard()
		}

		if editorExitRequested {
			EditorAskSaveChanged()
			if InputKeyPressed == KEY_ESCAPE {
				editorExitRequested = false
				EditorDrawSidebar()
			}
		}
		if editorExitRequested {
			break
		}
	}
	InputKeyPressed = '\x00'
	InitElementsGame()
}

func HighScoresLoad() {
	var (
		f FILE
		i int16
	)
	Assign(f, World.Info.Name+".HI")
	Reset(f)
	if IOResult == 0 {
		Read(f, HighScoreList)
	}
	Close(f)
	if IOResult != 0 {
		for i = 1; i <= 30; i++ {
			HighScoreList[i+1].Name = ""
			HighScoreList[i+1].Score = -1
		}
	}
}

func HighScoresSave() {
	var f FILE
	Assign(f, World.Info.Name+".HI")
	Rewrite(f)
	Write(f, HighScoreList)
	Close(f)
	if DisplayIOError {
	} else {
	}
}

func HighScoresInitTextWindow(state *TTextWindowState) {
	var (
		i        int16
		scoreStr string
	)
	TextWindowInitState(state)
	TextWindowAppend(state, "Score  Name")
	TextWindowAppend(state, "-----  ----------------------------------")
	for i = 1; i <= HIGH_SCORE_COUNT; i++ {
		if Length(HighScoreList[i+1].Name) != 0 {
			scoreStr = fmt.Sprintf("%5v", HighScoreList[i+1].Score)
			TextWindowAppend(state, scoreStr+"  "+HighScoreList[i+1].Name)
		}
	}
}

func HighScoresDisplay(linePos int16) {
	var state TTextWindowState
	state.LinePos = linePos
	HighScoresInitTextWindow(state)
	if state.LineCount > 2 {
		state.Title = "High scores for " + World.Info.Name
		TextWindowDrawOpen(state)
		TextWindowSelect(state, false, true)
		TextWindowDrawClose(state)
	}
	TextWindowFree(state)
}

func EditorOpenEditTextWindow(state *TTextWindowState) {
	SidebarClear()
	VideoWriteText(61, 4, 0x30, " Return ")
	VideoWriteText(64, 5, 0x1F, " Insert line")
	VideoWriteText(61, 7, 0x70, " Ctrl-Y ")
	VideoWriteText(64, 8, 0x1F, " Delete line")
	VideoWriteText(61, 10, 0x30, " Cursor keys ")
	VideoWriteText(64, 11, 0x1F, " Move cursor")
	VideoWriteText(61, 13, 0x70, " Insert ")
	VideoWriteText(64, 14, 0x1F, " Insert mode: ")
	VideoWriteText(61, 16, 0x30, " Delete ")
	VideoWriteText(64, 17, 0x1F, " Delete char")
	VideoWriteText(61, 19, 0x70, " Escape ")
	VideoWriteText(64, 20, 0x1F, " Exit editor")
	TextWindowEdit(state)
}

func EditorEditHelpFile() {
	var (
		textWindow TTextWindowState
		filename   string
	)
	filename = ""
	SidebarPromptString("File to edit", ".HLP", filename, PROMPT_ALPHANUM)
	if Length(filename) != 0 {
		TextWindowOpenFile('*'+filename+".HLP", textWindow)
		textWindow.Title = "Editing " + filename
		TextWindowDrawOpen(textWindow)
		EditorOpenEditTextWindow(textWindow)
		TextWindowSaveFile(filename+".HLP", textWindow)
		TextWindowFree(textWindow)
		TextWindowDrawClose(textWindow)
	}
}

func HighScoresAdd(score int16) {
	var (
		textWindow TTextWindowState
		name       string
		i, listPos int16
	)
	listPos = 1
	for (listPos <= 30) && (score < HighScoreList[listPos+1].Score) {
		listPos = listPos + 1
	}
	if (listPos <= 30) && (score > 0) {
		for i = 29; i >= listPos; i-- {
			HighScoreList[(i+1)+1] = HighScoreList[i+1]
		}
		HighScoreList[listPos+1].Score = score
		HighScoreList[listPos+1].Name = "-- You! --"
		HighScoresInitTextWindow(textWindow)
		textWindow.LinePos = listPos
		textWindow.Title = "New high score for " + World.Info.Name
		TextWindowDrawOpen(textWindow)
		TextWindowDraw(textWindow, false, false)
		name = ""
		PopupPromptString("Congratulations!  Enter your name:", name)
		HighScoreList[listPos+1].Name = name
		HighScoresSave()
		TextWindowDrawClose(textWindow)
		TransitionDrawToBoard()
		TextWindowFree(textWindow)
	}
}

func EditorGetBoardName(boardId int16, titleScreenIsNone bool) (EditorGetBoardName TString50) {
	var (
		boardData  uintptr
		copiedName string
	)
	if (boardId == 0) && titleScreenIsNone {
		EditorGetBoardName = "None"
	} else if boardId == World.Info.CurrentBoard {
		EditorGetBoardName = Board.Name
	} else {
		boardData = World.BoardData[boardId]
		Move(boardData, copiedName, SizeOf(copiedName))
		EditorGetBoardName = copiedName
	}

	return
}

func EditorSelectBoard(title string, currentBoard int16, titleScreenIsNone bool) (EditorSelectBoard int16) {
	var (
		unk1       string
		i          int16
		unk2       int16
		textWindow TTextWindowState
	)
	textWindow.Title = title
	textWindow.LinePos = currentBoard + 1
	textWindow.Selectable = true
	textWindow.LineCount = 0
	for i = 0; i <= World.BoardCount; i++ {
		TextWindowAppend(textWindow, EditorGetBoardName(i, titleScreenIsNone))
	}
	TextWindowAppend(textWindow, "Add new board")
	TextWindowDrawOpen(textWindow)
	TextWindowSelect(textWindow, false, false)
	TextWindowDrawClose(textWindow)
	TextWindowFree(textWindow)
	if InputKeyPressed == KEY_ESCAPE {
		EditorSelectBoard = 0
	} else {
		EditorSelectBoard = textWindow.LinePos - 1
	}
	return
}

func init() {
}
