package main // unit: Oop

// interface uses: GameVars

// implementation uses: Sounds, TxtWind, Game, Elements

func OopError(statId int16, message string) {
	stat := &Board.Stats[statId]
	DisplayMessage(200, "ERR: "+message)
	SoundQueue(5, "P\n")
	stat.DataPos = -1

}

func OopReadChar(statId int16, position *int16) {
	stat := &Board.Stats[statId]
	if (position >= 0) && (position < stat.DataLen) {
		Move(&Ptr(Seg(Data), Ofs(Data)+position), OopChar, 1)
		Inc(position)
	} else {
		OopChar = '\x00'
	}

}

func OopReadWord(statId int16, position *int16) {
	OopWord = ""
	for {
		OopReadChar(statId, position)
		if OopChar != ' ' {
			break
		}
	}
	OopChar = UpCase(OopChar)
	if (OopChar < '0') || (OopChar > '9') {
		for ((OopChar >= 'A') && (OopChar <= 'Z')) || (OopChar == ':') || ((OopChar >= '0') && (OopChar <= '9')) || (OopChar == '_') {
			OopWord = OopWord + OopChar
			OopReadChar(statId, position)
			OopChar = UpCase(OopChar)
		}
	}
	if position > 0 {
		Dec(position)
	}
}

func OopReadValue(statId int16, position *int16) {
	var (
		s    string
		code int16
	)
	s = ""
	for {
		OopReadChar(statId, position)
		if OopChar != ' ' {
			break
		}
	}
	OopChar = UpCase(OopChar)
	for (OopChar >= '0') && (OopChar <= '9') {
		s = s + OopChar
		OopReadChar(statId, position)
		OopChar = UpCase(OopChar)
	}
	if position > 0 {
		position = position - 1
	}
	if Length(s) != 0 {
		Val(s, OopValue, code)
	} else {
		OopValue = -1
	}
}

func OopSkipLine(statId int16, position *int16) {
	for {
		OopReadChar(statId, position)
		if (OopChar == '\x00') || (OopChar == '\r') {
			break
		}
	}
}

func OopParseDirection(statId int16, position *int16, dx, dy *int16) (OopParseDirection bool) {
	stat := &Board.Stats[statId]
	OopParseDirection = true
	if (OopWord == 'N') || (OopWord == "NORTH") {
		dx = 0
		dy = -1
	} else if (OopWord == 'S') || (OopWord == "SOUTH") {
		dx = 0
		dy = 1
	} else if (OopWord == 'E') || (OopWord == "EAST") {
		dx = 1
		dy = 0
	} else if (OopWord == 'W') || (OopWord == "WEST") {
		dx = -1
		dy = 0
	} else if (OopWord == 'I') || (OopWord == "IDLE") {
		dx = 0
		dy = 0
	} else if OopWord == "SEEK" {
		CalcDirectionSeek(stat.X, stat.Y, dx, dy)
	} else if OopWord == "FLOW" {
		dx = stat.StepX
		dy = stat.StepY
	} else if OopWord == "RND" {
		CalcDirectionRnd(dx, dy)
	} else if OopWord == "RNDNS" {
		dx = 0
		dy = Random(2)*2 - 1
	} else if OopWord == "RNDNE" {
		dx = Random(2)
		if dx == 0 {
			dy = -1
		} else {
			dy = 0
		}
	} else if OopWord == "CW" {
		OopReadWord(statId, position)
		OopParseDirection = OopParseDirection(statId, position, dy, dx)
		dx = -dx
	} else if OopWord == "CCW" {
		OopReadWord(statId, position)
		OopParseDirection = OopParseDirection(statId, position, dy, dx)
		dy = -dy
	} else if OopWord == "RNDP" {
		OopReadWord(statId, position)
		OopParseDirection = OopParseDirection(statId, position, dy, dx)
		if Random(2) == 0 {
			dx = -dx
		} else {
			dy = -dy
		}
	} else if OopWord == "OPP" {
		OopReadWord(statId, position)
		OopParseDirection = OopParseDirection(statId, position, dx, dy)
		dx = -dx
		dy = -dy
	} else {
		dx = 0
		dy = 0
		OopParseDirection = false
	}

	return
}

func OopReadDirection(statId int16, position *int16, dx, dy *int16) {
	OopReadWord(statId, position)
	if !OopParseDirection(statId, position, dx, dy) {
		OopError(statId, "Bad direction")
	}
}

func OopFindString(statId int16, s string) (OopFindString int16) {
	var pos, wordPos, cmpPos int16
	stat := &Board.Stats[statId]
	pos = 0
	for pos <= stat.DataLen {
		wordPos = 1
		cmpPos = pos
		for {
			OopReadChar(statId, cmpPos)
			if UpCase(s[wordPos]) != UpCase(OopChar) {
				goto NoMatch
			}
			wordPos = wordPos + 1
			if wordPos > Length(s) {
				break
			}
		}
		OopReadChar(statId, cmpPos)
		OopChar = UpCase(OopChar)
		if ((OopChar >= 'A') && (OopChar <= 'Z')) || (OopChar == '_') {
		} else {
			OopFindString = pos
			exit()
		}
	NoMatch:
		pos = pos + 1

	}
	OopFindString = -1

	return
}

func OopIterateStat(statId int16, iStat *int16, lookup string) (OopIterateStat bool) {
	var (
		pos   int16
		found bool
	)
	iStat = iStat + 1
	found = false
	if lookup == "ALL" {
		if iStat <= Board.StatCount {
			found = true
		}
	} else if lookup == "OTHERS" {
		if iStat <= Board.StatCount {
			if iStat != statId {
				found = true
			} else {
				iStat = iStat + 1
				found = (iStat <= Board.StatCount)
			}
		}
	} else if lookup == "SELF" {
		if (statId > 0) && (iStat <= statId) {
			iStat = statId
			found = true
		}
	} else {
		for (iStat <= Board.StatCount) && !found {
			if Board.Stats[iStat].Data != nil {
				pos = 0
				OopReadChar(iStat, pos)
				if OopChar == '@' {
					OopReadWord(iStat, pos)
					if OopWord == lookup {
						found = true
					}
				}
			}
			if !found {
				iStat = iStat + 1
			}
		}
	}

	OopIterateStat = found
	return
}

func OopFindLabel(statId int16, sendLabel string, iStat, iDataPos *int16, labelPrefix string) (OopFindLabel bool) {
	var (
		targetSplitPos int16
		unk1           int16
		targetLookup   string
		objectMessage  string
		foundStat      bool
	)
	foundStat = false
	targetSplitPos = Pos(':', sendLabel)
	if targetSplitPos <= 0 {
		if iStat < statId {
			objectMessage = sendLabel
			iStat = statId
			targetSplitPos = 0
			foundStat = true
		}
	} else {
		targetLookup = Copy(sendLabel, 1, targetSplitPos-1)
		objectMessage = Copy(sendLabel, targetSplitPos+1, Length(sendLabel)-targetSplitPos)
	FindNextStat:
		foundStat = OopIterateStat(statId, iStat, targetLookup)

	}
	if foundStat {
		if objectMessage == "RESTART" {
			iDataPos = 0
		} else {
			iDataPos = OopFindString(iStat, labelPrefix+objectMessage)
			if (iDataPos < 0) && (targetSplitPos > 0) {
				goto FindNextStat
			}
		}
		foundStat = iDataPos >= 0
	}
	OopFindLabel = foundStat
	return
}

func WorldGetFlagPosition(name TString50) (WorldGetFlagPosition int16) {
	var i int16
	WorldGetFlagPosition = -1
	for i = 1; i <= 10; i++ {
		if World.Info.Flags[i+1] == name {
			WorldGetFlagPosition = i
		}
	}
	return
}

func WorldSetFlag(name TString50) {
	var i int16
	if WorldGetFlagPosition(name) < 0 {
		i = 1
		for (i < MAX_FLAG) && (Length(World.Info.Flags[i+1]) != 0) {
			i = i + 1
		}
		World.Info.Flags[i+1] = name
	}
}

func WorldClearFlag(name TString50) {
	var i int16
	if WorldGetFlagPosition(name) >= 0 {
		World.Info.Flags[WorldGetFlagPosition(name)+1] = ""
	}
}

func OopStringToWord(input TString50) (OopStringToWord TString50) {
	var (
		output string
		i      int16
	)
	output = ""
	for i = 1; i <= Length(input); i++ {
		if ((input[i] >= 'A') && (input[i] <= 'Z')) || ((input[i] >= '0') && (input[i] <= '9')) {
			output = output + input[i]
		} else if (input[i] >= 'a') && (input[i] <= 'z') {
			output = output + Chr(Ord(input[i])-0x20)
		}

	}
	OopStringToWord = output
	return
}

func OopParseTile(statId, position *int16, tile *TTile) (OopParseTile bool) {
	var i int16
	OopParseTile = false
	tile.Color = 0
	OopReadWord(statId, position)
	for i = 1; i <= 7; i++ {
		if OopWord == OopStringToWord(ColorNames[i+1]) {
			tile.Color = i + 0x08
			OopReadWord(statId, position)
			goto ColorFound
		}
	}
ColorFound:
	for i = 0; i <= MAX_ELEMENT; i++ {
		if OopWord == OopStringToWord(ElementDefs[i].Name) {
			OopParseTile = true
			tile.Element = i
			exit()
		}
	}

	return
}

func GetColorForTileMatch(tile *TTile) (GetColorForTileMatch byte) {
	if ElementDefs[tile.Element].Color < COLOR_SPECIAL_MIN {
		GetColorForTileMatch = ElementDefs[tile.Element].Color && 0x07
	} else if ElementDefs[tile.Element].Color == COLOR_WHITE_ON_CHOICE {
		GetColorForTileMatch = ((tile.Color >> 4) && 0x0F) + 8
	} else {
		GetColorForTileMatch = (tile.Color && 0x0F)
	}

	return
}

func FindTileOnBoard(x, y *int16, tile TTile) (FindTileOnBoard bool) {
	FindTileOnBoard = false
	for true {
		x = x + 1
		if x > BOARD_WIDTH {
			x = 1
			y = y + 1
			if y > BOARD_HEIGHT {
				exit()
			}
		}
		if Board.Tiles[x][y].Element == tile.Element {
			if (tile.Color == 0) || (GetColorForTileMatch(Board.Tiles[x][y]) == tile.Color) {
				FindTileOnBoard = true
				exit()
			}
		}
	}
	return
}

func OopPlaceTile(x, y int16, tile *TTile) {
	var color byte
	if Board.Tiles[x][y].Element != 4 {
		color = tile.Color
		if ElementDefs[tile.Element].Color < COLOR_SPECIAL_MIN {
			color = ElementDefs[tile.Element].Color
		} else {
			if color == 0 {
				color = Board.Tiles[x][y].Color
			}
			if color == 0 {
				color = 0x0F
			}
			if ElementDefs[tile.Element].Color == COLOR_WHITE_ON_CHOICE {
				color = ((color - 8) * 0x10) + 0x0F
			}
		}
		if Board.Tiles[x][y].Element == tile.Element {
			Board.Tiles[x][y].Color = color
		} else {
			BoardDamageTile(x, y)
			if ElementDefs[tile.Element].Cycle >= 0 {
				AddStat(x, y, tile.Element, color, ElementDefs[tile.Element].Cycle, StatTemplateDefault)
			} else {
				Board.Tiles[x][y].Element = tile.Element
				Board.Tiles[x][y].Color = color
			}
		}
		BoardDrawTile(x, y)
	}
}

func OopCheckCondition(statId int16, position *int16) (OopCheckCondition bool) {
	var (
		deltaX, deltaY int16
		tile           TTile
		ix, iy         int16
	)
	stat := &Board.Stats[statId]
	if OopWord == "NOT" {
		OopReadWord(statId, position)
		OopCheckCondition = !OopCheckCondition(statId, position)
	} else if OopWord == "ALLIGNED" {
		OopCheckCondition = (stat.X == Board.Stats[0].X) || (stat.Y == Board.Stats[0].Y)
	} else if OopWord == "CONTACT" {
		OopCheckCondition = (Sqr(stat.X-Board.Stats[0].X) + Sqr(stat.Y-Board.Stats[0].Y)) == 1
	} else if OopWord == "BLOCKED" {
		OopReadDirection(statId, position, deltaX, deltaY)
		OopCheckCondition = !ElementDefs[Board.Tiles[stat.X+deltaX][stat.Y+deltaY].Element].Walkable
	} else if OopWord == "ENERGIZED" {
		OopCheckCondition = World.Info.EnergizerTicks > 0
	} else if OopWord == "ANY" {
		if !OopParseTile(statId, position, tile) {
			OopError(statId, "Bad object kind")
		}
		ix = 0
		iy = 1
		OopCheckCondition = FindTileOnBoard(ix, iy, tile)
	} else {
		OopCheckCondition = WorldGetFlagPosition(OopWord) >= 0
	}

	return
}

func OopReadLineToEnd(statId int16, position *int16) (OopReadLineToEnd string) {
	var s string
	s = ""
	OopReadChar(statId, position)
	for (OopChar != '\x00') && (OopChar != '\r') {
		s = s + OopChar
		OopReadChar(statId, position)
	}
	OopReadLineToEnd = s
	return
}

func OopSend(statId int16, sendLabel string, ignoreLock bool) (OopSend bool) {
	var (
		iDataPos, iStat int16
		ignoreSelfLock  bool
	)
	if statId < 0 {
		statId = -statId
		ignoreSelfLock = true
	} else {
		ignoreSelfLock = false
	}
	OopSend = false
	iStat = 0
	for OopFindLabel(statId, sendLabel, iStat, iDataPos, "\r:") {
		if ((Board.Stats[iStat].P2 == 0) || (ignoreLock)) || ((statId == iStat) && !ignoreSelfLock) {
			if iStat == statId {
				OopSend = true
			}
			Board.Stats[iStat].DataPos = iDataPos
		}
	}
	return
}

func OopExecute(statId int16, position *int16, name TString50) {
	var (
		textWindow        TTextWindowState
		textLine          string
		deltaX, deltaY    int16
		ix, iy            int16
		stopRunning       bool
		replaceStat       bool
		endOfProgram      bool
		replaceTile       TTile
		namePosition      int16
		lastPosition      int16
		repeatInsNextTick bool
		lineFinished      bool
		labelPtr          uintptr
		labelDataPos      int16
		labelStatId       int16
		counterPtr        *int16
		counterSubtract   bool
		bindStatId        int16
		insCount          int16
		argTile           TTile
		argTile2          TTile
	)
	stat := &Board.Stats[statId]
StartParsing:
	TextWindowInitState(textWindow)

	textWindow.Selectable = false
	stopRunning = false
	repeatInsNextTick = false
	replaceStat = false
	endOfProgram = false
	insCount = 0
	for {
	ReadInstruction:
		lineFinished = true

		lastPosition = position
		OopReadChar(statId, position)
		for OopChar == ':' {
			for {
				OopReadChar(statId, position)
				if (OopChar == '\x00') || (OopChar == '\r') {
					break
				}
			}
			OopReadChar(statId, position)
		}
		if OopChar == '\'' {
			OopSkipLine(statId, position)
		} else if OopChar == '@' {
			OopSkipLine(statId, position)
		} else if (OopChar == '/') || (OopChar == '?') {
			if OopChar == '/' {
				repeatInsNextTick = true
			}
			OopReadWord(statId, position)
			if OopParseDirection(statId, position, deltaX, deltaY) {
				if (deltaX != 0) || (deltaY != 0) {
					if !ElementDefs[Board.Tiles[stat.X+deltaX][stat.Y+deltaY].Element].Walkable {
						ElementPushablePush(stat.X+deltaX, stat.Y+deltaY, deltaX, deltaY)
					}
					if ElementDefs[Board.Tiles[stat.X+deltaX][stat.Y+deltaY].Element].Walkable {
						MoveStat(statId, stat.X+deltaX, stat.Y+deltaY)
						repeatInsNextTick = false
					}
				} else {
					repeatInsNextTick = false
				}
				OopReadChar(statId, position)
				if OopChar != '\r' {
					Dec(position)
				}
				stopRunning = true
			} else {
				OopError(statId, "Bad direction")
			}
		} else if OopChar == '#' {
		ReadCommand:
			OopReadWord(statId, position)

			if OopWord == "THEN" {
				OopReadWord(statId, position)
			}
			if Length(OopWord) == 0 {
				goto ReadInstruction
			}
			Inc(insCount)
			if Length(OopWord) != 0 {
				if OopWord == "GO" {
					OopReadDirection(statId, position, deltaX, deltaY)
					if !ElementDefs[Board.Tiles[stat.X+deltaX][stat.Y+deltaY].Element].Walkable {
						ElementPushablePush(stat.X+deltaX, stat.Y+deltaY, deltaX, deltaY)
					}
					if ElementDefs[Board.Tiles[stat.X+deltaX][stat.Y+deltaY].Element].Walkable {
						MoveStat(statId, stat.X+deltaX, stat.Y+deltaY)
					} else {
						repeatInsNextTick = true
					}
					stopRunning = true
				} else if OopWord == "TRY" {
					OopReadDirection(statId, position, deltaX, deltaY)
					if !ElementDefs[Board.Tiles[stat.X+deltaX][stat.Y+deltaY].Element].Walkable {
						ElementPushablePush(stat.X+deltaX, stat.Y+deltaY, deltaX, deltaY)
					}
					if ElementDefs[Board.Tiles[stat.X+deltaX][stat.Y+deltaY].Element].Walkable {
						MoveStat(statId, stat.X+deltaX, stat.Y+deltaY)
						stopRunning = true
					} else {
						goto ReadCommand
					}
				} else if OopWord == "WALK" {
					OopReadDirection(statId, position, deltaX, deltaY)
					stat.StepX = deltaX
					stat.StepY = deltaY
				} else if OopWord == "SET" {
					OopReadWord(statId, position)
					WorldSetFlag(OopWord)
				} else if OopWord == "CLEAR" {
					OopReadWord(statId, position)
					WorldClearFlag(OopWord)
				} else if OopWord == "IF" {
					OopReadWord(statId, position)
					if OopCheckCondition(statId, position) {
						goto ReadCommand
					}
				} else if OopWord == "SHOOT" {
					OopReadDirection(statId, position, deltaX, deltaY)
					if BoardShoot(E_BULLET, stat.X, stat.Y, deltaX, deltaY, SHOT_SOURCE_ENEMY) {
						SoundQueue(2, "0\x01&\x01")
					}
					stopRunning = true
				} else if OopWord == "THROWSTAR" {
					OopReadDirection(statId, position, deltaX, deltaY)
					if BoardShoot(E_STAR, stat.X, stat.Y, deltaX, deltaY, SHOT_SOURCE_ENEMY) {
					}
					stopRunning = true
				} else if (OopWord == "GIVE") || (OopWord == "TAKE") {
					if OopWord == "TAKE" {
						counterSubtract = true
					} else {
						counterSubtract = false
					}
					OopReadWord(statId, position)
					if OopWord == "HEALTH" {
						counterPtr = *World.Info.Health
					} else if OopWord == "AMMO" {
						counterPtr = *World.Info.Ammo
					} else if OopWord == "GEMS" {
						counterPtr = *World.Info.Gems
					} else if OopWord == "TORCHES" {
						counterPtr = *World.Info.Torches
					} else if OopWord == "SCORE" {
						counterPtr = *World.Info.Score
					} else if OopWord == "TIME" {
						counterPtr = *World.Info.BoardTimeSec
					} else {
						counterPtr = nil
					}

					if counterPtr != nil {
						OopReadValue(statId, position)
						if OopValue > 0 {
							if counterSubtract {
								OopValue = -OopValue
							}
							if (counterPtr + OopValue) >= 0 {
								counterPtr = counterPtr + OopValue
							} else {
								goto ReadCommand
							}
						}
					}
					GameUpdateSidebar()
				} else if OopWord == "END" {
					position = -1
					OopChar = '\x00'
				} else if OopWord == "ENDGAME" {
					World.Info.Health = 0
				} else if OopWord == "IDLE" {
					stopRunning = true
				} else if OopWord == "RESTART" {
					position = 0
					lineFinished = false
				} else if OopWord == "ZAP" {
					OopReadWord(statId, position)
					labelStatId = 0
					for OopFindLabel(statId, OopWord, labelStatId, labelDataPos, "\r:") {
						labelPtr = Board.Stats[labelStatId].Data
						AdvancePointer(labelPtr, labelDataPos+1)
						labelPtr = '\''
					}
				} else if OopWord == "RESTORE" {
					OopReadWord(statId, position)
					labelStatId = 0
					for OopFindLabel(statId, OopWord, labelStatId, labelDataPos, "\r'") {
						for {
							labelPtr = Board.Stats[labelStatId].Data
							AdvancePointer(labelPtr, labelDataPos+1)
							labelPtr = ':'
							labelDataPos = OopFindString(labelStatId, "\r'"+OopWord+'\r')
							if labelDataPos <= 0 {
								break
							}
						}
					}
				} else if OopWord == "LOCK" {
					stat.P2 = 1
				} else if OopWord == "UNLOCK" {
					stat.P2 = 0
				} else if OopWord == "SEND" {
					OopReadWord(statId, position)
					if OopSend(statId, OopWord, false) {
						lineFinished = false
					}
				} else if OopWord == "BECOME" {
					if OopParseTile(statId, position, argTile) {
						replaceStat = true
						replaceTile.Element = argTile.Element
						replaceTile.Color = argTile.Color
					} else {
						OopError(statId, "Bad #BECOME")
					}
				} else if OopWord == "PUT" {
					OopReadDirection(statId, position, deltaX, deltaY)
					if (deltaX == 0) && (deltaY == 0) {
						OopError(statId, "Bad #PUT")
					} else if !OopParseTile(statId, position, argTile) {
						OopError(statId, "Bad #PUT")
					} else if ((stat.X + deltaX) > 0) && ((stat.X + deltaX) <= BOARD_WIDTH) && ((stat.Y + deltaY) > 0) && ((stat.Y + deltaY) < BOARD_HEIGHT) {
						if !ElementDefs[Board.Tiles[stat.X+deltaX][stat.Y+deltaY].Element].Walkable {
							ElementPushablePush(stat.X+deltaX, stat.Y+deltaY, deltaX, deltaY)
						}
						OopPlaceTile(stat.X+deltaX, stat.Y+deltaY, argTile)
					}

				} else if OopWord == "CHANGE" {
					if !OopParseTile(statId, position, argTile) {
						OopError(statId, "Bad #CHANGE")
					}
					if !OopParseTile(statId, position, argTile2) {
						OopError(statId, "Bad #CHANGE")
					}
					ix = 0
					iy = 1
					if (argTile2.Color == 0) && (ElementDefs[argTile2.Element].Color < COLOR_SPECIAL_MIN) {
						argTile2.Color = ElementDefs[argTile2.Element].Color
					}
					for FindTileOnBoard(ix, iy, argTile) {
						OopPlaceTile(ix, iy, argTile2)
					}
				} else if OopWord == "PLAY" {
					textLine = SoundParse(OopReadLineToEnd(statId, position))
					if Length(textLine) != 0 {
						SoundQueue(-1, textLine)
					}
					lineFinished = false
				} else if OopWord == "CYCLE" {
					OopReadValue(statId, position)
					if OopValue > 0 {
						stat.Cycle = OopValue
					}
				} else if OopWord == "CHAR" {
					OopReadValue(statId, position)
					if (OopValue > 0) && (OopValue <= 255) {
						stat.P1 = OopValue
						BoardDrawTile(stat.X, stat.Y)
					}
				} else if OopWord == "DIE" {
					replaceStat = true
					replaceTile.Element = E_EMPTY
					replaceTile.Color = 0x0F
				} else if OopWord == "BIND" {
					OopReadWord(statId, position)
					bindStatId = 0
					if OopIterateStat(statId, bindStatId, OopWord) {
						FreeMem(stat.Data, stat.DataLen)
						stat.Data = Board.Stats[bindStatId].Data
						stat.DataLen = Board.Stats[bindStatId].DataLen
						position = 0
					}
				} else {
					textLine = OopWord
					if OopSend(statId, OopWord, false) {
						lineFinished = false
					} else {
						if Pos(':', textLine) <= 0 {
							OopError(statId, "Bad command "+textLine)
						}
					}
				}

			}
			if lineFinished {
				OopSkipLine(statId, position)
			}
		} else if OopChar == '\r' {
			if textWindow.LineCount > 0 {
				TextWindowAppend(textWindow, "")
			}
		} else if OopChar == '\x00' {
			endOfProgram = true
		} else {
			textLine = OopChar + OopReadLineToEnd(statId, position)
			TextWindowAppend(textWindow, textLine)
		}

		if endOfProgram || stopRunning || repeatInsNextTick || replaceStat || (insCount > 32) {
			break
		}
	}
	if repeatInsNextTick {
		position = lastPosition
	}
	if OopChar == '\x00' {
		position = -1
	}
	if textWindow.LineCount > 1 {
		namePosition = 0
		OopReadChar(statId, namePosition)
		if OopChar == '@' {
			name = OopReadLineToEnd(statId, namePosition)
		}
		if Length(name) == 0 {
			name = "Interaction"
		}
		textWindow.Title = name
		TextWindowDrawOpen(textWindow)
		TextWindowSelect(textWindow, true, false)
		TextWindowDrawClose(textWindow)
		TextWindowFree(textWindow)
		if Length(textWindow.Hyperlink) != 0 {
			if OopSend(statId, textWindow.Hyperlink, false) {
				goto StartParsing
			}
		}
	} else if textWindow.LineCount == 1 {
		DisplayMessage(200, textWindow.Lines[2])
		TextWindowFree(textWindow)
	}

	if replaceStat {
		ix = stat.X
		iy = stat.Y
		DamageStat(statId)
		OopPlaceTile(ix, iy, replaceTile)
	}

}

func init() {
}
