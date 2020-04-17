package main // unit: Elements

// interface uses: GameVars

// implementation uses: Crt, Video, Sounds, Input, TxtWind, Editor, Oop, Game

const (
	TransporterNSChars string = "^~^-v_v-"
	TransporterEWChars string = "(<(\xb3)>)\xb3"
	StarAnimChars      string = "\xb3/\xc4\\"
)

func ElementDefaultTick(statId int16) {
}

func ElementDefaultTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
}

func ElementDefaultDraw(x, y int16, ch *byte) {
	*ch = Ord('?')
}

func ElementMessageTimerTick(statId int16) {
	stat := &Board.Stats[statId]
	switch stat.X {
	case 0:
		VideoWriteText((60-Length(Board.Info.Message))/2, 24, 9+(stat.P2%7), " "+Board.Info.Message+" ")
		stat.P2--
		if stat.P2 <= 0 {
			RemoveStat(statId)
			CurrentStatTicked--
			BoardDrawBorder()
			Board.Info.Message = ""
		}
	}
}

func ElementDamagingTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	BoardAttack(sourceStatId, x, y)
}

func ElementLionTick(statId int16) {
	var deltaX, deltaY int16
	stat := &Board.Stats[statId]
	if int16(stat.P1) < Random(10) {
		CalcDirectionRnd(&deltaX, &deltaY)
	} else {
		CalcDirectionSeek(int16(stat.X), int16(stat.Y), &deltaX, &deltaY)
	}
	if ElementDefs[Board.Tiles[int16(stat.X)+deltaX][int16(stat.Y)+deltaY].Element].Walkable {
		MoveStat(statId, int16(stat.X)+deltaX, int16(stat.Y)+deltaY)
	} else if Board.Tiles[int16(stat.X)+deltaX][int16(stat.Y)+deltaY].Element == E_PLAYER {
		BoardAttack(statId, int16(stat.X)+deltaX, int16(stat.Y)+deltaY)
	}

}

func ElementTigerTick(statId int16) {
	var (
		shot    bool
		element byte
	)
	stat := &Board.Stats[statId]
	element = E_BULLET
	if stat.P2 >= 0x80 {
		element = E_STAR
	}
	if (Random(10) * 3) <= int16(stat.P2%0x80) {
		if Difference(int16(stat.X), int16(Board.Stats[0].X)) <= 2 {
			shot = BoardShoot(element, int16(stat.X), int16(stat.Y), 0, Signum(int16(Board.Stats[0].Y-stat.Y)), SHOT_SOURCE_ENEMY)
		} else {
			shot = false
		}
		if !shot {
			if Difference(int16(stat.Y), int16(Board.Stats[0].Y)) <= 2 {
				shot = BoardShoot(element, int16(stat.X), int16(stat.Y), Signum(int16(Board.Stats[0].X-stat.X)), 0, SHOT_SOURCE_ENEMY)
			}
		}
	}
	ElementLionTick(statId)
}

func ElementRuffianTick(statId int16) {
	stat := &Board.Stats[statId]
	if (stat.StepX == 0) && (stat.StepY == 0) {
		if int16(stat.P2+8) <= Random(17) {
			if int16(stat.P1) >= Random(9) {
				CalcDirectionSeek(int16(stat.X), int16(stat.Y), &stat.StepX, &stat.StepY)
			} else {
				CalcDirectionRnd(&stat.StepX, &stat.StepY)
			}
		}
	} else {
		if ((stat.Y == Board.Stats[0].Y) || (stat.X == Board.Stats[0].X)) && (Random(9) <= int16(stat.P1)) {
			CalcDirectionSeek(int16(stat.X), int16(stat.Y), &stat.StepX, &stat.StepY)
		}
		tile := &Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY]
		if tile.Element == E_PLAYER {
			BoardAttack(statId, int16(stat.X)+stat.StepX, int16(stat.Y)+stat.StepY)
		} else if ElementDefs[tile.Element].Walkable {
			MoveStat(statId, int16(stat.X)+stat.StepX, int16(stat.Y)+stat.StepY)
			if int16(stat.P2+8) <= Random(17) {
				stat.StepX = 0
				stat.StepY = 0
			}
		} else {
			stat.StepX = 0
			stat.StepY = 0
		}

	}
}

func ElementBearTick(statId int16) {
	var deltaX, deltaY int16
	stat := &Board.Stats[statId]
	if stat.X != Board.Stats[0].X {
		if Difference(int16(stat.Y), int16(Board.Stats[0].Y)) <= int16(8-stat.P1) {
			deltaX = Signum(int16(Board.Stats[0].X - stat.X))
			deltaY = 0
			goto Movement
		}
	}
	if Difference(int16(stat.X), int16(Board.Stats[0].X)) <= int16(8-stat.P1) {
		deltaY = Signum(int16(Board.Stats[0].Y - stat.Y))
		deltaX = 0
	} else {
		deltaX = 0
		deltaY = 0
	}
Movement:
	tile := &Board.Tiles[int16(stat.X)+deltaX][int16(stat.Y)+deltaY]
	if ElementDefs[tile.Element].Walkable {
		MoveStat(statId, int16(stat.X)+deltaX, int16(stat.Y)+deltaY)
	} else if (tile.Element == E_PLAYER) || (tile.Element == E_BREAKABLE) {
		BoardAttack(statId, int16(stat.X)+deltaX, int16(stat.Y)+deltaY)
	}

}

func ElementCentipedeHeadTick(statId int16) {
	var (
		ix, iy int16
		tx, ty int16
		tmp    int16
	)
	stat := &Board.Stats[statId]
	if (stat.X == Board.Stats[0].X) && (Random(10) < int16(stat.P1)) {
		stat.StepY = Signum(int16(Board.Stats[0].Y - stat.Y))
		stat.StepX = 0
	} else if (stat.Y == Board.Stats[0].Y) && (Random(10) < int16(stat.P1)) {
		stat.StepX = Signum(int16(Board.Stats[0].X - stat.X))
		stat.StepY = 0
	} else if ((Random(10) * 4) < int16(stat.P2)) || ((stat.StepX == 0) && (stat.StepY == 0)) {
		CalcDirectionRnd(&stat.StepX, &stat.StepY)
	}

	if !ElementDefs[Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY].Element].Walkable && (Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY].Element != E_PLAYER) {
		ix = stat.StepX
		iy = stat.StepY
		tmp = ((Random(2) * 2) - 1) * stat.StepY
		stat.StepY = ((Random(2) * 2) - 1) * stat.StepX
		stat.StepX = tmp
		if !ElementDefs[Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY].Element].Walkable && (Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY].Element != E_PLAYER) {
			stat.StepX = -stat.StepX
			stat.StepY = -stat.StepY
			if !ElementDefs[Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY].Element].Walkable && (Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY].Element != E_PLAYER) {
				if ElementDefs[Board.Tiles[int16(stat.X)-ix][int16(stat.Y)-iy].Element].Walkable || (Board.Tiles[int16(stat.X)-ix][int16(stat.Y)-iy].Element == E_PLAYER) {
					stat.StepX = -ix
					stat.StepY = -iy
				} else {
					stat.StepX = 0
					stat.StepY = 0
				}
			}
		}
	}
	if (stat.StepX == 0) && (stat.StepY == 0) {
		Board.Tiles[stat.X][stat.Y].Element = E_CENTIPEDE_SEGMENT
		stat.Leader = -1
		for Board.Stats[statId].Follower > 0 {
			tmp = Board.Stats[statId].Follower
			Board.Stats[statId].Follower = Board.Stats[statId].Leader
			Board.Stats[statId].Leader = tmp
			statId = tmp
		}
		Board.Stats[statId].Follower = Board.Stats[statId].Leader
		Board.Tiles[Board.Stats[statId].X][Board.Stats[statId].Y].Element = E_CENTIPEDE_HEAD
	} else if Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY].Element == E_PLAYER {
		if stat.Follower != -1 {
			Board.Tiles[Board.Stats[stat.Follower].X][Board.Stats[stat.Follower].Y].Element = E_CENTIPEDE_HEAD
			Board.Stats[stat.Follower].StepX = stat.StepX
			Board.Stats[stat.Follower].StepY = stat.StepY
			BoardDrawTile(int16(Board.Stats[stat.Follower].X), int16(Board.Stats[stat.Follower].Y))
		}
		BoardAttack(statId, int16(stat.X)+stat.StepX, int16(stat.Y)+stat.StepY)
	} else {
		MoveStat(statId, int16(stat.X)+stat.StepX, int16(stat.Y)+stat.StepY)
		tx = int16(stat.X) - stat.StepX
		ty = int16(stat.Y) - stat.StepY
		ix = stat.StepX
		iy = stat.StepY
		for {
			stat2 := &Board.Stats[statId]
			tx = int16(stat2.X) - stat2.StepX
			ty = int16(stat2.Y) - stat2.StepY
			ix = stat2.StepX
			iy = stat2.StepY
			if stat2.Follower < 0 {
				if (Board.Tiles[tx-ix][ty-iy].Element == E_CENTIPEDE_SEGMENT) && (Board.Stats[GetStatIdAt(tx-ix, ty-iy)].Leader < 0) {
					stat2.Follower = GetStatIdAt(tx-ix, ty-iy)
				} else if (Board.Tiles[tx-iy][ty-ix].Element == E_CENTIPEDE_SEGMENT) && (Board.Stats[GetStatIdAt(tx-iy, ty-ix)].Leader < 0) {
					stat2.Follower = GetStatIdAt(tx-iy, ty-ix)
				} else if (Board.Tiles[tx+iy][ty+ix].Element == E_CENTIPEDE_SEGMENT) && (Board.Stats[GetStatIdAt(tx+iy, ty+ix)].Leader < 0) {
					stat2.Follower = GetStatIdAt(tx+iy, ty+ix)
				}

			}
			if stat2.Follower > 0 {
				Board.Stats[stat2.Follower].Leader = statId
				Board.Stats[stat2.Follower].P1 = stat2.P1
				Board.Stats[stat2.Follower].P2 = stat2.P2
				Board.Stats[stat2.Follower].StepX = tx - int16(Board.Stats[stat2.Follower].X)
				Board.Stats[stat2.Follower].StepY = ty - int16(Board.Stats[stat2.Follower].Y)
				MoveStat(stat2.Follower, tx, ty)
			}
			statId = stat2.Follower
			if statId == -1 {
				break
			}
		}
	}

}

func ElementCentipedeSegmentTick(statId int16) {
	stat := &Board.Stats[statId]
	if stat.Leader < 0 {
		if stat.Leader < -1 {
			Board.Tiles[stat.X][stat.Y].Element = E_CENTIPEDE_HEAD
		} else {
			stat.Leader--
		}
	}
}

func ElementBulletTick(statId int16) {
	var (
		ix, iy   int16
		iStat    int16
		iElem    byte
		firstTry bool
	)
	stat := &Board.Stats[statId]
	firstTry = true
TryMove:
	ix = int16(stat.X) + stat.StepX

	iy = int16(stat.Y) + stat.StepY
	iElem = Board.Tiles[ix][iy].Element
	if ElementDefs[iElem].Walkable || (iElem == E_WATER) {
		MoveStat(statId, ix, iy)
		return
	}
	if (iElem == E_RICOCHET) && firstTry {
		stat.StepX = -stat.StepX
		stat.StepY = -stat.StepY
		SoundQueue(1, "\xf9\x01")
		firstTry = false
		goto TryMove
		return
	}
	if (iElem == E_BREAKABLE) || (ElementDefs[iElem].Destructible && ((iElem == E_PLAYER) || (stat.P1 == 0))) {
		if ElementDefs[iElem].ScoreValue != 0 {
			World.Info.Score += ElementDefs[iElem].ScoreValue
			GameUpdateSidebar()
		}
		BoardAttack(statId, ix, iy)
		return
	}
	if (Board.Tiles[int16(stat.X)+stat.StepY][int16(stat.Y)+stat.StepX].Element == E_RICOCHET) && firstTry {
		ix = stat.StepX
		stat.StepX = -stat.StepY
		stat.StepY = -ix
		SoundQueue(1, "\xf9\x01")
		firstTry = false
		goto TryMove
		return
	}
	if (Board.Tiles[int16(stat.X)-stat.StepY][int16(stat.Y)-stat.StepX].Element == E_RICOCHET) && firstTry {
		ix = stat.StepX
		stat.StepX = stat.StepY
		stat.StepY = ix
		SoundQueue(1, "\xf9\x01")
		firstTry = false
		goto TryMove
		return
	}
	RemoveStat(statId)
	CurrentStatTicked--
	if (iElem == E_OBJECT) || (iElem == E_SCROLL) {
		iStat = GetStatIdAt(ix, iy)
		if OopSend(-iStat, "SHOT", false) {
		}
	}
}

func ElementSpinningGunDraw(x, y int16, ch *byte) {
	switch CurrentTick % 8 {
	case 0, 1:
		*ch = 24
	case 2, 3:
		*ch = 26
	case 4, 5:
		*ch = 25
	default:
		*ch = 27
	}
}

func ElementLineDraw(x, y int16, ch *byte) {
	var i, v, shift int16
	v = 1
	shift = 1
	for i = 0; i <= 3; i++ {
		switch Board.Tiles[x+NeighborDeltaX[i]][y+NeighborDeltaY[i]].Element {
		case E_LINE, E_BOARD_EDGE:
			v += shift
		}
		shift = shift << 1
	}
	*ch = Ord(LineChars[v-1])
}

func ElementSpinningGunTick(statId int16) {
	var (
		shot           bool
		deltaX, deltaY int16
		element        byte
	)
	stat := &Board.Stats[statId]
	BoardDrawTile(int16(stat.X), int16(stat.Y))
	element = E_BULLET
	if stat.P2 >= 0x80 {
		element = E_STAR
	}
	if Random(9) < int16(stat.P2%0x80) {
		if Random(9) <= int16(stat.P1) {
			if Difference(int16(stat.X), int16(Board.Stats[0].X)) <= 2 {
				shot = BoardShoot(element, int16(stat.X), int16(stat.Y), 0, Signum(int16(Board.Stats[0].Y-stat.Y)), SHOT_SOURCE_ENEMY)
			} else {
				shot = false
			}
			if !shot {
				if Difference(int16(stat.Y), int16(Board.Stats[0].Y)) <= 2 {
					shot = BoardShoot(element, int16(stat.X), int16(stat.Y), Signum(int16(Board.Stats[0].X-stat.X)), 0, SHOT_SOURCE_ENEMY)
				}
			}
		} else {
			CalcDirectionRnd(&deltaX, &deltaY)
			shot = BoardShoot(element, int16(stat.X), int16(stat.Y), deltaX, deltaY, SHOT_SOURCE_ENEMY)
		}
	}
}

func ElementConveyorTick(x, y int16, direction int16) {
	var (
		i          int16
		iStat      int16
		ix, iy     int16
		canMove    bool
		tiles      [8]TTile
		iMin, iMax int16
		tmpTile    TTile
	)
	if direction == 1 {
		iMin = 0
		iMax = 8
	} else {
		iMin = 7
		iMax = -1
	}
	canMove = true
	i = iMin
	for {
		tiles[i] = Board.Tiles[x+DiagonalDeltaX[i]][y+DiagonalDeltaY[i]]
		tile := &tiles[i]
		if tile.Element == E_EMPTY {
			canMove = true
		} else if !ElementDefs[tile.Element].Pushable {
			canMove = false
		}

		i += direction
		if i == iMax {
			break
		}
	}
	i = iMin
	for {
		tile2 := &tiles[i]
		if canMove {
			if ElementDefs[tile2.Element].Pushable {
				ix = x + DiagonalDeltaX[(i-direction+8)%8]
				iy = y + DiagonalDeltaY[(i-direction+8)%8]
				if ElementDefs[tile2.Element].Cycle > -1 {
					tmpTile = Board.Tiles[x+DiagonalDeltaX[i]][y+DiagonalDeltaY[i]]
					iStat = GetStatIdAt(x+DiagonalDeltaX[i], y+DiagonalDeltaY[i])
					Board.Tiles[x+DiagonalDeltaX[i]][y+DiagonalDeltaY[i]] = tiles[i]
					Board.Tiles[ix][iy].Element = E_EMPTY
					MoveStat(iStat, ix, iy)
					Board.Tiles[x+DiagonalDeltaX[i]][y+DiagonalDeltaY[i]] = tmpTile
				} else {
					Board.Tiles[ix][iy] = tiles[i]
					BoardDrawTile(ix, iy)
				}
				if !ElementDefs[tiles[(i+direction+8)%8].Element].Pushable {
					Board.Tiles[x+DiagonalDeltaX[i]][y+DiagonalDeltaY[i]].Element = E_EMPTY
					BoardDrawTile(x+DiagonalDeltaX[i], y+DiagonalDeltaY[i])
				}
			} else {
				canMove = false
			}
		} else if tile2.Element == E_EMPTY {
			canMove = true
		} else if !ElementDefs[tile2.Element].Pushable {
			canMove = false
		}

		i += direction
		if i == iMax {
			break
		}
	}
}

func ElementConveyorCWDraw(x, y int16, ch *byte) {
	switch (CurrentTick / ElementDefs[E_CONVEYOR_CW].Cycle) % 4 {
	case 0:
		*ch = 179
	case 1:
		*ch = 47
	case 2:
		*ch = 196
	default:
		*ch = 92
	}
}

func ElementConveyorCWTick(statId int16) {
	stat := &Board.Stats[statId]
	BoardDrawTile(int16(stat.X), int16(stat.Y))
	ElementConveyorTick(int16(stat.X), int16(stat.Y), 1)
}

func ElementConveyorCCWDraw(x, y int16, ch *byte) {
	switch (CurrentTick / ElementDefs[E_CONVEYOR_CCW].Cycle) % 4 {
	case 3:
		*ch = 179
	case 2:
		*ch = 47
	case 1:
		*ch = 196
	default:
		*ch = 92
	}
}

func ElementConveyorCCWTick(statId int16) {
	stat := &Board.Stats[statId]
	BoardDrawTile(int16(stat.X), int16(stat.Y))
	ElementConveyorTick(int16(stat.X), int16(stat.Y), -1)
}

func ElementBombDraw(x, y int16, ch *byte) {
	stat := &Board.Stats[GetStatIdAt(x, y)]
	if stat.P1 <= 1 {
		*ch = 11
	} else {
		*ch = 48 + stat.P1
	}
}

func ElementBombTick(statId int16) {
	var oldX, oldY int16
	stat := &Board.Stats[statId]
	if stat.P1 > 0 {
		stat.P1--
		BoardDrawTile(int16(stat.X), int16(stat.Y))
		if stat.P1 == 1 {
			SoundQueue(1, "`\x01P\x01@\x010\x01 \x01\x10\x01")
			DrawPlayerSurroundings(int16(stat.X), int16(stat.Y), 1)
		} else if stat.P1 == 0 {
			oldX = int16(stat.X)
			oldY = int16(stat.Y)
			RemoveStat(statId)
			DrawPlayerSurroundings(oldX, oldY, 2)
		} else {
			if (stat.P1 % 2) == 0 {
				SoundQueue(1, "\xf8\x01")
			} else {
				SoundQueue(1, "\xf5\x01")
			}
		}

	}
}

func ElementBombTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	stat := &Board.Stats[GetStatIdAt(x, y)]
	if stat.P1 == 0 {
		stat.P1 = 9
		BoardDrawTile(int16(stat.X), int16(stat.Y))
		DisplayMessage(200, "Bomb activated!")
		SoundQueue(4, "0\x015\x01@\x01E\x01P\x01")
	} else {
		ElementPushablePush(int16(stat.X), int16(stat.Y), *deltaX, *deltaY)
	}
}

func ElementTransporterMove(x, y, deltaX, deltaY int16) {
	var (
		ix, iy       int16
		newX, newY   int16
		iStat        int16
		finishSearch bool
		isValidDest  bool
	)
	stat := &Board.Stats[GetStatIdAt(x+deltaX, y+deltaY)]
	if (deltaX == stat.StepX) && (deltaY == stat.StepY) {
		ix = int16(stat.X)
		iy = int16(stat.Y)
		newX = -1
		finishSearch = false
		isValidDest = true
		for {
			ix += deltaX
			iy += deltaY
			tile := &Board.Tiles[ix][iy]
			if tile.Element == E_BOARD_EDGE {
				finishSearch = true
			} else if isValidDest {
				isValidDest = false
				if !ElementDefs[tile.Element].Walkable {
					ElementPushablePush(ix, iy, deltaX, deltaY)
				}
				if ElementDefs[tile.Element].Walkable {
					finishSearch = true
					newX = ix
					newY = iy
				} else {
					newX = -1
				}
			}

			if tile.Element == E_TRANSPORTER {
				iStat = GetStatIdAt(ix, iy)
				if (Board.Stats[iStat].StepX == -deltaX) && (Board.Stats[iStat].StepY == -deltaY) {
					isValidDest = true
				}
			}
			if finishSearch {
				break
			}
		}
		if newX != -1 {
			ElementMove(int16(stat.X)-deltaX, int16(stat.Y)-deltaY, newX, newY)
			SoundQueue(3, "0\x01B\x014\x01F\x018\x01J\x01@\x01R\x01")
		}
	}
}

func ElementTransporterTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	ElementTransporterMove(x-*deltaX, y-*deltaY, *deltaX, *deltaY)
	*deltaX = 0
	*deltaY = 0
}

func ElementTransporterTick(statId int16) {
	stat := &Board.Stats[statId]
	BoardDrawTile(int16(stat.X), int16(stat.Y))
}

func ElementTransporterDraw(x, y int16, ch *byte) {
	stat := &Board.Stats[GetStatIdAt(x, y)]
	if stat.StepX == 0 {
		*ch = Ord(TransporterNSChars[stat.StepY*2+3+(CurrentTick/stat.Cycle)%4])
	} else {
		*ch = Ord(TransporterEWChars[stat.StepX*2+3+(CurrentTick/stat.Cycle)%4])
	}
}

func ElementStarDraw(x, y int16, ch *byte) {
	*ch = Ord(StarAnimChars[(CurrentTick%4)+1])
	Board.Tiles[x][y].Color++
	if Board.Tiles[x][y].Color > 15 {
		Board.Tiles[x][y].Color = 9
	}
}

func ElementStarTick(statId int16) {
	stat := &Board.Stats[statId]
	stat.P2--
	if stat.P2 <= 0 {
		RemoveStat(statId)
	} else if (stat.P2 % 2) == 0 {
		CalcDirectionSeek(int16(stat.X), int16(stat.Y), &stat.StepX, &stat.StepY)
		tile := &Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY]
		if (tile.Element == E_PLAYER) || (tile.Element == E_BREAKABLE) {
			BoardAttack(statId, int16(stat.X)+stat.StepX, int16(stat.Y)+stat.StepY)
		} else {
			if !ElementDefs[tile.Element].Walkable {
				ElementPushablePush(int16(stat.X)+stat.StepX, int16(stat.Y)+stat.StepY, stat.StepX, stat.StepY)
			}
			if ElementDefs[tile.Element].Walkable || (tile.Element == E_WATER) {
				MoveStat(statId, int16(stat.X)+stat.StepX, int16(stat.Y)+stat.StepY)
			}
		}
	} else {
		BoardDrawTile(int16(stat.X), int16(stat.Y))
	}

}

func ElementEnergizerTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	SoundQueue(9, " \x03#\x03$\x03%\x035\x03%\x03#\x03 \x03"+"0\x03#\x03$\x03%\x035\x03%\x03#\x03 \x03"+"0\x03#\x03$\x03%\x035\x03%\x03#\x03 \x03"+"0\x03#\x03$\x03%\x035\x03%\x03#\x03 \x03"+"0\x03#\x03$\x03%\x035\x03%\x03#\x03 \x03"+"0\x03#\x03$\x03%\x035\x03%\x03#\x03 \x03"+"0\x03#\x03$\x03%\x035\x03%\x03#\x03 \x03")
	Board.Tiles[x][y].Element = E_EMPTY
	BoardDrawTile(x, y)
	World.Info.EnergizerTicks = 75
	GameUpdateSidebar()
	if MessageEnergizerNotShown {
		DisplayMessage(200, "Energizer - You are invincible")
		MessageEnergizerNotShown = false
	}
	if OopSend(0, "ALL:ENERGIZE", false) {
	}
}

func ElementSlimeTick(statId int16) {
	var (
		dir, color, changedTiles int16
		startX, startY           int16
	)
	stat := &Board.Stats[statId]
	if stat.P1 < stat.P2 {
		stat.P1++
	} else {
		color = int16(Board.Tiles[stat.X][stat.Y].Color)
		stat.P1 = 0
		startX = int16(stat.X)
		startY = int16(stat.Y)
		changedTiles = 0
		for dir = 0; dir <= 3; dir++ {
			if ElementDefs[Board.Tiles[startX+NeighborDeltaX[dir]][startY+NeighborDeltaY[dir]].Element].Walkable {
				if changedTiles == 0 {
					MoveStat(statId, startX+NeighborDeltaX[dir], startY+NeighborDeltaY[dir])
					Board.Tiles[startX][startY].Color = byte(color)
					Board.Tiles[startX][startY].Element = E_BREAKABLE
					BoardDrawTile(startX, startY)
				} else {
					AddStat(startX+NeighborDeltaX[dir], startY+NeighborDeltaY[dir], E_SLIME, color, ElementDefs[E_SLIME].Cycle, StatTemplateDefault)
					Board.Stats[Board.StatCount].P2 = stat.P2
				}
				changedTiles++
			}
		}
		if changedTiles == 0 {
			RemoveStat(statId)
			Board.Tiles[startX][startY].Element = E_BREAKABLE
			Board.Tiles[startX][startY].Color = byte(color)
			BoardDrawTile(startX, startY)
		}
	}
}

func ElementSlimeTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	var color int16
	color = int16(Board.Tiles[x][y].Color)
	DamageStat(GetStatIdAt(x, y))
	Board.Tiles[x][y].Element = E_BREAKABLE
	Board.Tiles[x][y].Color = byte(color)
	BoardDrawTile(x, y)
	SoundQueue(2, " \x01#\x01")
}

func ElementSharkTick(statId int16) {
	var deltaX, deltaY int16
	stat := &Board.Stats[statId]
	if int16(stat.P1) < Random(10) {
		CalcDirectionRnd(&deltaX, &deltaY)
	} else {
		CalcDirectionSeek(int16(stat.X), int16(stat.Y), &deltaX, &deltaY)
	}
	if Board.Tiles[int16(stat.X)+deltaX][int16(stat.Y)+deltaY].Element == E_WATER {
		MoveStat(statId, int16(stat.X)+deltaX, int16(stat.Y)+deltaY)
	} else if Board.Tiles[int16(stat.X)+deltaX][int16(stat.Y)+deltaY].Element == E_PLAYER {
		BoardAttack(statId, int16(stat.X)+deltaX, int16(stat.Y)+deltaY)
	}

}

func ElementBlinkWallDraw(x, y int16, ch *byte) {
	*ch = 206
}

func ElementBlinkWallTick(statId int16) {
	var (
		ix, iy       int16
		hitBoundary  bool
		playerStatId int16
		el           int16
	)
	stat := &Board.Stats[statId]
	if stat.P3 == 0 {
		stat.P3 = stat.P1 + 1
	}
	if stat.P3 == 1 {
		ix = int16(stat.X) + stat.StepX
		iy = int16(stat.Y) + stat.StepY
		if stat.StepX != 0 {
			el = E_BLINK_RAY_EW
		} else {
			el = E_BLINK_RAY_NS
		}
		for (int16(Board.Tiles[ix][iy].Element) == el) && (Board.Tiles[ix][iy].Color == Board.Tiles[stat.X][stat.Y].Color) {
			Board.Tiles[ix][iy].Element = E_EMPTY
			BoardDrawTile(ix, iy)
			ix += stat.StepX
			iy += stat.StepY
			stat.P3 = (stat.P2)*2 + 1
		}
		if ((int16(stat.X) + stat.StepX) == ix) && ((int16(stat.Y) + stat.StepY) == iy) {
			hitBoundary = false
			for {
				if (Board.Tiles[ix][iy].Element != E_EMPTY) && (ElementDefs[Board.Tiles[ix][iy].Element].Destructible) {
					BoardDamageTile(ix, iy)
				}
				if Board.Tiles[ix][iy].Element == E_PLAYER {
					playerStatId = GetStatIdAt(ix, iy)
					if stat.StepX != 0 {
						if Board.Tiles[ix][iy-1].Element == E_EMPTY {
							MoveStat(playerStatId, ix, iy-1)
						} else if Board.Tiles[ix][iy+1].Element == E_EMPTY {
							MoveStat(playerStatId, ix, iy+1)
						}

					} else {
						if Board.Tiles[ix+1][iy].Element == E_EMPTY {
							MoveStat(playerStatId, ix+1, iy)
						} else if Board.Tiles[ix-1][iy].Element == E_EMPTY {
							MoveStat(playerStatId, ix+1, iy)
						}

					}
					if Board.Tiles[ix][iy].Element == E_PLAYER {
						for World.Info.Health > 0 {
							DamageStat(playerStatId)
						}
						hitBoundary = true
					}
				}
				if Board.Tiles[ix][iy].Element == E_EMPTY {
					Board.Tiles[ix][iy].Element = byte(el)
					Board.Tiles[ix][iy].Color = Board.Tiles[stat.X][stat.Y].Color
					BoardDrawTile(ix, iy)
				} else {
					hitBoundary = true
				}
				ix += stat.StepX
				iy += stat.StepY
				if hitBoundary {
					break
				}
			}
			stat.P3 = (stat.P2 * 2) + 1
		}
	} else {
		stat.P3--
	}
}

func ElementMove(oldX, oldY, newX, newY int16) {
	var statId int16
	statId = GetStatIdAt(oldX, oldY)
	if statId >= 0 {
		MoveStat(statId, newX, newY)
	} else {
		Board.Tiles[newX][newY] = Board.Tiles[oldX][oldY]
		BoardDrawTile(newX, newY)
		Board.Tiles[oldX][oldY].Element = E_EMPTY
		BoardDrawTile(oldX, oldY)
	}
}

func ElementPushablePush(x, y int16, deltaX, deltaY int16) {
	var unk1 int16
	tile := &Board.Tiles[x][y]
	if ((tile.Element == E_SLIDER_NS) && (deltaX == 0)) || ((tile.Element == E_SLIDER_EW) && (deltaY == 0)) || ElementDefs[tile.Element].Pushable {
		if Board.Tiles[x+deltaX][y+deltaY].Element == E_TRANSPORTER {
			ElementTransporterMove(x, y, deltaX, deltaY)
		} else if Board.Tiles[x+deltaX][y+deltaY].Element != E_EMPTY {
			ElementPushablePush(x+deltaX, y+deltaY, deltaX, deltaY)
		}

		if !ElementDefs[Board.Tiles[x+deltaX][y+deltaY].Element].Walkable && ElementDefs[Board.Tiles[x+deltaX][y+deltaY].Element].Destructible && (Board.Tiles[x+deltaX][y+deltaY].Element != E_PLAYER) {
			BoardDamageTile(x+deltaX, y+deltaY)
		}
		if ElementDefs[Board.Tiles[x+deltaX][y+deltaY].Element].Walkable {
			ElementMove(x, y, x+deltaX, y+deltaY)
		}
	}
}

func ElementDuplicatorDraw(x, y int16, ch *byte) {
	stat := &Board.Stats[GetStatIdAt(x, y)]
	switch stat.P1 {
	case 1:
		*ch = 250
	case 2:
		*ch = 249
	case 3:
		*ch = 248
	case 4:
		*ch = 111
	case 5:
		*ch = 79
	default:
		*ch = 250
	}
}

func ElementObjectTick(statId int16) {
	var retVal bool
	stat := &Board.Stats[statId]
	if stat.DataPos >= 0 {
		OopExecute(statId, &stat.DataPos, "Interaction")
	}
	if (stat.StepX != 0) || (stat.StepY != 0) {
		if ElementDefs[Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY].Element].Walkable {
			MoveStat(statId, int16(stat.X)+stat.StepX, int16(stat.Y)+stat.StepY)
		} else {
			retVal = OopSend(-statId, "THUD", false)
		}
	}
}

func ElementObjectDraw(x, y int16, ch *byte) {
	*ch = Board.Stats[GetStatIdAt(x, y)].P1
}

func ElementObjectTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	var (
		statId int16
		retVal bool
	)
	statId = GetStatIdAt(x, y)
	retVal = OopSend(-statId, "TOUCH", false)
}

func ElementDuplicatorTick(statId int16) {
	var sourceStatId int16
	stat := &Board.Stats[statId]
	if stat.P1 <= 4 {
		stat.P1++
		BoardDrawTile(int16(stat.X), int16(stat.Y))
	} else {
		stat.P1 = 0
		if Board.Tiles[int16(stat.X)-stat.StepX][int16(stat.Y)-stat.StepY].Element == E_PLAYER {
			ElementDefs[Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY].Element].TouchProc(int16(stat.X)+stat.StepX, int16(stat.Y)+stat.StepY, 0, &InputDeltaX, &InputDeltaY)
		} else {
			if Board.Tiles[int16(stat.X)-stat.StepX][int16(stat.Y)-stat.StepY].Element != E_EMPTY {
				ElementPushablePush(int16(stat.X)-stat.StepX, int16(stat.Y)-stat.StepY, -stat.StepX, -stat.StepY)
			}
			if Board.Tiles[int16(stat.X)-stat.StepX][int16(stat.Y)-stat.StepY].Element == E_EMPTY {
				sourceStatId = GetStatIdAt(int16(stat.X)+stat.StepX, int16(stat.Y)+stat.StepY)
				if sourceStatId > 0 {
					if Board.StatCount < 174 {
						AddStat(int16(stat.X)-stat.StepX, int16(stat.Y)-stat.StepY, Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY].Element, int16(Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY].Color), Board.Stats[sourceStatId].Cycle, Board.Stats[sourceStatId])
						BoardDrawTile(int16(stat.X)-stat.StepX, int16(stat.Y)-stat.StepY)
					}
				} else if sourceStatId != 0 {
					Board.Tiles[int16(stat.X)-stat.StepX][int16(stat.Y)-stat.StepY] = Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY]
					BoardDrawTile(int16(stat.X)-stat.StepX, int16(stat.Y)-stat.StepY)
				}

				SoundQueue(3, "0\x022\x024\x025\x027\x02")
			} else {
				SoundQueue(3, "\x18\x01\x16\x01")
			}
		}
		stat.P1 = 0
		BoardDrawTile(int16(stat.X), int16(stat.Y))
	}
	stat.Cycle = int16((9 - stat.P2) * 3)
}

func ElementScrollTick(statId int16) {
	stat := &Board.Stats[statId]
	Board.Tiles[stat.X][stat.Y].Color++
	if Board.Tiles[stat.X][stat.Y].Color > 0x0F {
		Board.Tiles[stat.X][stat.Y].Color = 0x09
	}
	BoardDrawTile(int16(stat.X), int16(stat.Y))
}

func ElementScrollTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	var (
		textWindow TTextWindowState
		statId     int16
		unk1       int16
	)
	statId = GetStatIdAt(x, y)
	stat := &Board.Stats[statId]
	textWindow.Selectable = false
	textWindow.LinePos = 1
	SoundQueue(2, SoundParse("c-c+d-d+e-e+f-f+g-g"))
	stat.DataPos = 0
	OopExecute(statId, &stat.DataPos, "Scroll")
	RemoveStat(GetStatIdAt(x, y))
}

func ElementKeyTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	var key int16
	key = int16(Board.Tiles[x][y].Color % 8)
	if World.Info.Keys[key-1] {
		DisplayMessage(200, "You already have a "+ColorNames[key-1]+" key!")
		SoundQueue(2, "0\x02 \x02")
	} else {
		World.Info.Keys[key-1] = true
		Board.Tiles[x][y].Element = E_EMPTY
		GameUpdateSidebar()
		DisplayMessage(200, "You now have the "+ColorNames[key-1]+" key.")
		SoundQueue(2, "@\x01D\x01G\x01@\x01D\x01G\x01@\x01D\x01G\x01P\x02")
	}
}

func ElementAmmoTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	World.Info.Ammo += 5
	Board.Tiles[x][y].Element = E_EMPTY
	GameUpdateSidebar()
	SoundQueue(2, "0\x011\x012\x01")
	if MessageAmmoNotShown {
		MessageAmmoNotShown = false
		DisplayMessage(200, "Ammunition - 5 shots per container.")
	}
}

func ElementGemTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	World.Info.Gems++
	World.Info.Health++
	World.Info.Score += 10
	Board.Tiles[x][y].Element = E_EMPTY
	GameUpdateSidebar()
	SoundQueue(2, "@\x017\x014\x010\x01")
	if MessageGemNotShown {
		MessageGemNotShown = false
		DisplayMessage(200, "Gems give you Health!")
	}
}

func ElementPassageTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	BoardPassageTeleport(x, y)
	*deltaX = 0
	*deltaY = 0
}

func ElementDoorTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	var key int16
	key = int16((Board.Tiles[x][y].Color / 16) % 8)
	if World.Info.Keys[key-1] {
		Board.Tiles[x][y].Element = E_EMPTY
		BoardDrawTile(x, y)
		World.Info.Keys[key-1] = false
		GameUpdateSidebar()
		DisplayMessage(200, "The "+ColorNames[key-1]+" door is now open.")
		SoundQueue(3, "0\x017\x01;\x010\x017\x01;\x01@\x04")
	} else {
		DisplayMessage(200, "The "+ColorNames[key-1]+" door is locked!")
		SoundQueue(3, "\x17\x01\x10\x01")
	}
}

func ElementPushableTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	ElementPushablePush(x, y, *deltaX, *deltaY)
	SoundQueue(2, "\x15\x01")
}

func ElementPusherDraw(x, y int16, ch *byte) {
	stat := &Board.Stats[GetStatIdAt(x, y)]
	if stat.StepX == 1 {
		*ch = 16
	} else if stat.StepX == -1 {
		*ch = 17
	} else if stat.StepY == -1 {
		*ch = 30
	} else {
		*ch = 31
	}

}

func ElementPusherTick(statId int16) {
	var i, startX, startY int16
	stat := &Board.Stats[statId]
	startX = int16(stat.X)
	startY = int16(stat.Y)
	if !ElementDefs[Board.Tiles[int16(stat.X)+stat.StepX][int16(stat.Y)+stat.StepY].Element].Walkable {
		ElementPushablePush(int16(stat.X)+stat.StepX, int16(stat.Y)+stat.StepY, stat.StepX, stat.StepY)
	}
	statId = GetStatIdAt(startX, startY)
	stat2 := &Board.Stats[statId]
	if ElementDefs[Board.Tiles[int16(stat2.X)+stat2.StepX][int16(stat2.Y)+stat2.StepY].Element].Walkable {
		MoveStat(statId, int16(stat2.X)+stat2.StepX, int16(stat2.Y)+stat2.StepY)
		SoundQueue(2, "\x15\x01")
		if Board.Tiles[int16(stat2.X)-(stat2.StepX*2)][int16(stat2.Y)-(stat2.StepY*2)].Element == E_PUSHER {
			i = GetStatIdAt(int16(stat2.X)-(stat2.StepX*2), int16(stat2.Y)-(stat2.StepY*2))
			if (Board.Stats[i].StepX == stat2.StepX) && (Board.Stats[i].StepY == stat2.StepY) {
				ElementDefs[E_PUSHER].TickProc(i)
			}
		}
	}
}

func ElementTorchTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	World.Info.Torches++
	Board.Tiles[x][y].Element = E_EMPTY
	BoardDrawTile(x, y)
	GameUpdateSidebar()
	if MessageTorchNotShown {
		DisplayMessage(200, "Torch - used for lighting in the underground.")
	}
	MessageTorchNotShown = false
	SoundQueue(3, "0\x019\x014\x02")
}

func ElementInvisibleTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	tile := &Board.Tiles[x][y]
	tile.Element = E_NORMAL
	BoardDrawTile(x, y)
	SoundQueue(3, "\x12\x01\x10\x01")
	DisplayMessage(100, "You are blocked by an invisible wall.")
}

func ElementForestTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	Board.Tiles[x][y].Element = E_EMPTY
	BoardDrawTile(x, y)
	SoundQueue(3, "9\x01")
	if MessageForestNotShown {
		DisplayMessage(200, "A path is cleared through the forest.")
	}
	MessageForestNotShown = false
}

func ElementFakeTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	if MessageFakeNotShown {
		DisplayMessage(150, "A fake wall - secret passage!")
	}
	MessageFakeNotShown = false
}

func ElementBoardEdgeTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	var (
		neighborId     int16
		boardId        int16
		entryX, entryY int16
	)
	entryX = int16(Board.Stats[0].X)
	entryY = int16(Board.Stats[0].Y)
	if *deltaY == -1 {
		neighborId = 0
		entryY = BOARD_HEIGHT
	} else if *deltaY == 1 {
		neighborId = 1
		entryY = 1
	} else if *deltaX == -1 {
		neighborId = 2
		entryX = BOARD_WIDTH
	} else {
		neighborId = 3
		entryX = 1
	}

	if Board.Info.NeighborBoards[neighborId] != 0 {
		boardId = World.Info.CurrentBoard
		BoardChange(int16(Board.Info.NeighborBoards[neighborId]))
		if Board.Tiles[entryX][entryY].Element != E_PLAYER {
			ElementDefs[Board.Tiles[entryX][entryY].Element].TouchProc(entryX, entryY, sourceStatId, &InputDeltaX, &InputDeltaY)
		}
		if ElementDefs[Board.Tiles[entryX][entryY].Element].Walkable || (Board.Tiles[entryX][entryY].Element == E_PLAYER) {
			if Board.Tiles[entryX][entryY].Element != E_PLAYER {
				MoveStat(0, entryX, entryY)
			}
			TransitionDrawBoardChange()
			*deltaX = 0
			*deltaY = 0
			BoardEnter()
		} else {
			BoardChange(boardId)
		}
	}
}

func ElementWaterTouch(x, y int16, sourceStatId int16, deltaX, deltaY *int16) {
	SoundQueue(3, "@\x01P\x01")
	DisplayMessage(100, "Your way is blocked by water.")
}

func DrawPlayerSurroundings(x, y int16, bombPhase int16) {
	var (
		ix, iy int16
		istat  int16
		result bool
	)
	for ix = (x - TORCH_DX) - 1; ix <= ((x + TORCH_DX) + 1); ix++ {
		if (ix >= 1) && (ix <= BOARD_WIDTH) {
			for iy = (y - TORCH_DY) - 1; iy <= ((y + TORCH_DY) + 1); iy++ {
				if (iy >= 1) && (iy <= BOARD_HEIGHT) {
					tile := &Board.Tiles[ix][iy]
					if (bombPhase > 0) && ((Sqr(ix-x) + Sqr(iy-y)*2) < TORCH_DIST_SQR) {
						if bombPhase == 1 {
							if Length(ElementDefs[tile.Element].ParamTextName) != 0 {
								istat = GetStatIdAt(ix, iy)
								if istat > 0 {
									result = OopSend(-istat, "BOMBED", false)
								}
							}
							if ElementDefs[tile.Element].Destructible || (tile.Element == E_STAR) {
								BoardDamageTile(ix, iy)
							}
							if (tile.Element == E_EMPTY) || (tile.Element == E_BREAKABLE) {
								tile.Element = E_BREAKABLE
								tile.Color = byte(0x09 + Random(7))
								BoardDrawTile(ix, iy)
							}
						} else {
							if tile.Element == E_BREAKABLE {
								tile.Element = E_EMPTY
							}
						}
					}
					BoardDrawTile(ix, iy)
				}
			}
		}
	}
}

func GamePromptEndPlay() {
	if World.Info.Health <= 0 {
		GamePlayExitRequested = true
		BoardDrawBorder()
	} else {
		GamePlayExitRequested = SidebarPromptYesNo("End this game? ", true)
		if InputKeyPressed == '\x1b' {
			GamePlayExitRequested = false
		}
	}
	InputKeyPressed = '\x00'
}

func ElementPlayerTick(statId int16) {
	var (
		unk1, unk2, unk3 int16
		i                int16
		bulletCount      int16
	)
	stat := &Board.Stats[statId]
	if World.Info.EnergizerTicks > 0 {
		if ElementDefs[E_PLAYER].Character == '\x02' {
			ElementDefs[E_PLAYER].Character = '\x01'
		} else {
			ElementDefs[E_PLAYER].Character = '\x02'
		}
		if (CurrentTick % 2) != 0 {
			Board.Tiles[stat.X][stat.Y].Color = 0x0F
		} else {
			Board.Tiles[stat.X][stat.Y].Color = byte((((CurrentTick % 7) + 1) * 16) + 0x0F)
		}
		BoardDrawTile(int16(stat.X), int16(stat.Y))
	} else if (Board.Tiles[stat.X][stat.Y].Color != 0x1F) || (ElementDefs[E_PLAYER].Character != '\x02') {
		Board.Tiles[stat.X][stat.Y].Color = 0x1F
		ElementDefs[E_PLAYER].Character = '\x02'
		BoardDrawTile(int16(stat.X), int16(stat.Y))
	}

	if World.Info.Health <= 0 {
		InputDeltaX = 0
		InputDeltaY = 0
		InputShiftPressed = false
		if GetStatIdAt(0, 0) == -1 {
			DisplayMessage(32000, " Game over  -  Press ESCAPE")
		}
		TickTimeDuration = 0
		SoundBlockQueueing = true
	}
	if InputShiftPressed || (InputKeyPressed == ' ') {
		if InputShiftPressed && ((InputDeltaX != 0) || (InputDeltaY != 0)) {
			PlayerDirX = InputDeltaX
			PlayerDirY = InputDeltaY
		}
		if (PlayerDirX != 0) || (PlayerDirY != 0) {
			if Board.Info.MaxShots == 0 {
				if MessageNoShootingNotShown {
					DisplayMessage(200, "Can't shoot in this place!")
				}
				MessageNoShootingNotShown = false
			} else if World.Info.Ammo == 0 {
				if MessageOutOfAmmoNotShown {
					DisplayMessage(200, "You don't have any ammo!")
				}
				MessageOutOfAmmoNotShown = false
			} else {
				bulletCount = 0
				for i = 0; i <= Board.StatCount; i++ {
					if (Board.Tiles[Board.Stats[i].X][Board.Stats[i].Y].Element == E_BULLET) && (Board.Stats[i].P1 == 0) {
						bulletCount++
					}
				}
				if bulletCount < int16(Board.Info.MaxShots) {
					if BoardShoot(E_BULLET, int16(stat.X), int16(stat.Y), PlayerDirX, PlayerDirY, SHOT_SOURCE_PLAYER) {
						World.Info.Ammo--
						GameUpdateSidebar()
						SoundQueue(2, "@\x010\x01 \x01")
						InputDeltaX = 0
						InputDeltaY = 0
					}
				}
			}

		}
	} else if (InputDeltaX != 0) || (InputDeltaY != 0) {
		PlayerDirX = InputDeltaX
		PlayerDirY = InputDeltaY
		ElementDefs[Board.Tiles[int16(stat.X)+InputDeltaX][int16(stat.Y)+InputDeltaY].Element].TouchProc(int16(stat.X)+InputDeltaX, int16(stat.Y)+InputDeltaY, 0, &InputDeltaX, &InputDeltaY)
		if (InputDeltaX != 0) || (InputDeltaY != 0) {
			if SoundEnabled && !SoundIsPlaying {
				Sound(110)
			}
			if ElementDefs[Board.Tiles[int16(stat.X)+InputDeltaX][int16(stat.Y)+InputDeltaY].Element].Walkable {
				if SoundEnabled && !SoundIsPlaying {
					NoSound()
				}
				MoveStat(0, int16(stat.X)+InputDeltaX, int16(stat.Y)+InputDeltaY)
			} else if SoundEnabled && !SoundIsPlaying {
				NoSound()
			}

		}
	}

	switch UpCase(InputKeyPressed) {
	case 'T':
		if World.Info.TorchTicks <= 0 {
			if World.Info.Torches > 0 {
				if Board.Info.IsDark {
					World.Info.Torches--
					World.Info.TorchTicks = TORCH_DURATION
					DrawPlayerSurroundings(int16(stat.X), int16(stat.Y), 0)
					GameUpdateSidebar()
				} else {
					if MessageRoomNotDarkNotShown {
						DisplayMessage(200, "Don't need torch - room is not dark!")
						MessageRoomNotDarkNotShown = false
					}
				}
			} else {
				if MessageOutOfTorchesNotShown {
					DisplayMessage(200, "You don't have any torches!")
					MessageOutOfTorchesNotShown = false
				}
			}
		}
	case '\x1b', 'Q':
		GamePromptEndPlay()
	case 'S':
		GameWorldSave("Save game:", &SavedGameFileName, ".SAV")
	case 'P':
		if World.Info.Health > 0 {
			GamePaused = true
		}
	case 'B':
		SoundEnabled = !SoundEnabled
		SoundClearQueue()
		GameUpdateSidebar()
		InputKeyPressed = ' '
	case 'H':
		TextWindowDisplayFile("GAME.HLP", "Playing ZZT")
	case 'F':
		TextWindowDisplayFile("ORDER.HLP", "Order form")
	case '?':
		GameDebugPrompt()
		InputKeyPressed = '\x00'
	}
	if World.Info.TorchTicks > 0 {
		World.Info.TorchTicks--
		if World.Info.TorchTicks <= 0 {
			DrawPlayerSurroundings(int16(stat.X), int16(stat.Y), 0)
			SoundQueue(3, "0\x01 \x01\x10\x01")
		}
		if (World.Info.TorchTicks % 40) == 0 {
			GameUpdateSidebar()
		}
	}
	if World.Info.EnergizerTicks > 0 {
		World.Info.EnergizerTicks--
		if World.Info.EnergizerTicks == 10 {
			SoundQueue(9, " \x03\x1a\x03\x17\x03\x16\x03\x15\x03\x13\x03\x10\x03")
		} else if World.Info.EnergizerTicks <= 0 {
			Board.Tiles[stat.X][stat.Y].Color = ElementDefs[E_PLAYER].Color
			BoardDrawTile(int16(stat.X), int16(stat.Y))
		}

	}
	if (Board.Info.TimeLimitSec > 0) && (World.Info.Health > 0) {
		if SoundHasTimeElapsed(&World.Info.BoardTimeHsec, 100) {
			World.Info.BoardTimeSec++
			if (Board.Info.TimeLimitSec - 10) == World.Info.BoardTimeSec {
				DisplayMessage(200, "Running out of time!")
				SoundQueue(3, "@\x06E\x06@\x065\x06@\x06E\x06@\n")
			} else if World.Info.BoardTimeSec > Board.Info.TimeLimitSec {
				DamageStat(0)
			}

			GameUpdateSidebar()
		}
	}
}

func ElementMonitorTick(statId int16) {
	if UpCase(InputKeyPressed) == '\x1b' || UpCase(InputKeyPressed) == 'A' || UpCase(InputKeyPressed) == 'E' || UpCase(InputKeyPressed) == 'H' || UpCase(InputKeyPressed) == 'N' || UpCase(InputKeyPressed) == 'P' || UpCase(InputKeyPressed) == 'Q' || UpCase(InputKeyPressed) == 'R' || UpCase(InputKeyPressed) == 'S' || UpCase(InputKeyPressed) == 'W' || UpCase(InputKeyPressed) == '|' {
		GamePlayExitRequested = true
	}
}

func ResetMessageNotShownFlags() {
	MessageAmmoNotShown = true
	MessageOutOfAmmoNotShown = true
	MessageNoShootingNotShown = true
	MessageTorchNotShown = true
	MessageOutOfTorchesNotShown = true
	MessageRoomNotDarkNotShown = true
	MessageHintTorchNotShown = true
	MessageForestNotShown = true
	MessageFakeNotShown = true
	MessageGemNotShown = true
	MessageEnergizerNotShown = true
}

func InitElementDefs() {
	var i int16
	for i = 0; i <= MAX_ELEMENT; i++ {
		def := &ElementDefs[i]
		def.Character = ' '
		def.Color = COLOR_CHOICE_ON_BLACK
		def.Destructible = false
		def.Pushable = false
		def.VisibleInDark = false
		def.PlaceableOnTop = false
		def.Walkable = false
		def.HasDrawProc = false
		def.Cycle = -1
		def.TickProc = ElementDefaultTick
		def.DrawProc = ElementDefaultDraw
		def.TouchProc = ElementDefaultTouch
		def.EditorCategory = 0
		def.EditorShortcut = '\x00'
		def.Name = ""
		def.CategoryName = ""
		def.Param1Name = ""
		def.Param2Name = ""
		def.ParamBulletTypeName = ""
		def.ParamBoardName = ""
		def.ParamDirName = ""
		def.ParamTextName = ""
		def.ScoreValue = 0
	}
	ElementDefs[0].Character = ' '
	ElementDefs[0].Color = 0x70
	ElementDefs[0].Pushable = true
	ElementDefs[0].Walkable = true
	ElementDefs[0].Name = "Empty"
	ElementDefs[3].Character = ' '
	ElementDefs[3].Color = 0x07
	ElementDefs[3].Cycle = 1
	ElementDefs[3].TickProc = ElementMonitorTick
	ElementDefs[3].Name = "Monitor"
	ElementDefs[19].Character = '\xb0'
	ElementDefs[19].Color = 0xF9
	ElementDefs[19].PlaceableOnTop = true
	ElementDefs[19].EditorCategory = CATEGORY_TERRAIN
	ElementDefs[19].TouchProc = ElementWaterTouch
	ElementDefs[19].EditorShortcut = 'W'
	ElementDefs[19].Name = "Water"
	ElementDefs[19].CategoryName = "Terrains:"
	ElementDefs[20].Character = '\xb0'
	ElementDefs[20].Color = 0x20
	ElementDefs[20].Walkable = false
	ElementDefs[20].TouchProc = ElementForestTouch
	ElementDefs[20].EditorCategory = CATEGORY_TERRAIN
	ElementDefs[20].EditorShortcut = 'F'
	ElementDefs[20].Name = "Forest"
	ElementDefs[4].Character = '\x02'
	ElementDefs[4].Color = 0x1F
	ElementDefs[4].Destructible = true
	ElementDefs[4].Pushable = true
	ElementDefs[4].VisibleInDark = true
	ElementDefs[4].Cycle = 1
	ElementDefs[4].TickProc = ElementPlayerTick
	ElementDefs[4].EditorCategory = CATEGORY_ITEM
	ElementDefs[4].EditorShortcut = 'Z'
	ElementDefs[4].Name = "Player"
	ElementDefs[4].CategoryName = "Items:"
	ElementDefs[41].Character = '\xea'
	ElementDefs[41].Color = 0x0C
	ElementDefs[41].Destructible = true
	ElementDefs[41].Pushable = true
	ElementDefs[41].Cycle = 2
	ElementDefs[41].TickProc = ElementLionTick
	ElementDefs[41].TouchProc = ElementDamagingTouch
	ElementDefs[41].EditorCategory = CATEGORY_CREATURE
	ElementDefs[41].EditorShortcut = 'L'
	ElementDefs[41].Name = "Lion"
	ElementDefs[41].CategoryName = "Beasts:"
	ElementDefs[41].Param1Name = "Intelligence?"
	ElementDefs[41].ScoreValue = 1
	ElementDefs[42].Character = '\xe3'
	ElementDefs[42].Color = 0x0B
	ElementDefs[42].Destructible = true
	ElementDefs[42].Pushable = true
	ElementDefs[42].Cycle = 2
	ElementDefs[42].TickProc = ElementTigerTick
	ElementDefs[42].TouchProc = ElementDamagingTouch
	ElementDefs[42].EditorCategory = CATEGORY_CREATURE
	ElementDefs[42].EditorShortcut = 'T'
	ElementDefs[42].Name = "Tiger"
	ElementDefs[42].Param1Name = "Intelligence?"
	ElementDefs[42].Param2Name = "Firing rate?"
	ElementDefs[42].ParamBulletTypeName = "Firing type?"
	ElementDefs[42].ScoreValue = 2
	ElementDefs[44].Character = '\xe9'
	ElementDefs[44].Destructible = true
	ElementDefs[44].Cycle = 2
	ElementDefs[44].TickProc = ElementCentipedeHeadTick
	ElementDefs[44].TouchProc = ElementDamagingTouch
	ElementDefs[44].EditorCategory = CATEGORY_CREATURE
	ElementDefs[44].EditorShortcut = 'H'
	ElementDefs[44].Name = "Head"
	ElementDefs[44].CategoryName = "Centipedes"
	ElementDefs[44].Param1Name = "Intelligence?"
	ElementDefs[44].Param2Name = "Deviance?"
	ElementDefs[44].ScoreValue = 1
	ElementDefs[45].Character = 'O'
	ElementDefs[45].Destructible = true
	ElementDefs[45].Cycle = 2
	ElementDefs[45].TickProc = ElementCentipedeSegmentTick
	ElementDefs[45].TouchProc = ElementDamagingTouch
	ElementDefs[45].EditorCategory = CATEGORY_CREATURE
	ElementDefs[45].EditorShortcut = 'S'
	ElementDefs[45].Name = "Segment"
	ElementDefs[45].ScoreValue = 3
	ElementDefs[18].Character = '\xf8'
	ElementDefs[18].Color = 0x0F
	ElementDefs[18].Destructible = true
	ElementDefs[18].Cycle = 1
	ElementDefs[18].TickProc = ElementBulletTick
	ElementDefs[18].TouchProc = ElementDamagingTouch
	ElementDefs[18].Name = "Bullet"
	ElementDefs[15].Character = 'S'
	ElementDefs[15].Color = 0x0F
	ElementDefs[15].Destructible = false
	ElementDefs[15].Cycle = 1
	ElementDefs[15].TickProc = ElementStarTick
	ElementDefs[15].TouchProc = ElementDamagingTouch
	ElementDefs[15].HasDrawProc = true
	ElementDefs[15].DrawProc = ElementStarDraw
	ElementDefs[15].Name = "Star"
	ElementDefs[8].Character = '\x0c'
	ElementDefs[8].Pushable = true
	ElementDefs[8].TouchProc = ElementKeyTouch
	ElementDefs[8].EditorCategory = CATEGORY_ITEM
	ElementDefs[8].EditorShortcut = 'K'
	ElementDefs[8].Name = "Key"
	ElementDefs[5].Character = '\x84'
	ElementDefs[5].Color = 0x03
	ElementDefs[5].Pushable = true
	ElementDefs[5].TouchProc = ElementAmmoTouch
	ElementDefs[5].EditorCategory = CATEGORY_ITEM
	ElementDefs[5].EditorShortcut = 'A'
	ElementDefs[5].Name = "Ammo"
	ElementDefs[7].Character = '\x04'
	ElementDefs[7].Pushable = true
	ElementDefs[7].TouchProc = ElementGemTouch
	ElementDefs[7].Destructible = true
	ElementDefs[7].EditorCategory = CATEGORY_ITEM
	ElementDefs[7].EditorShortcut = 'G'
	ElementDefs[7].Name = "Gem"
	ElementDefs[11].Character = '\xf0'
	ElementDefs[11].Color = COLOR_WHITE_ON_CHOICE
	ElementDefs[11].Cycle = 0
	ElementDefs[11].VisibleInDark = true
	ElementDefs[11].TouchProc = ElementPassageTouch
	ElementDefs[11].EditorCategory = CATEGORY_ITEM
	ElementDefs[11].EditorShortcut = 'P'
	ElementDefs[11].Name = "Passage"
	ElementDefs[11].ParamBoardName = "Room thru passage?"
	ElementDefs[9].Character = '\n'
	ElementDefs[9].Color = COLOR_WHITE_ON_CHOICE
	ElementDefs[9].TouchProc = ElementDoorTouch
	ElementDefs[9].EditorCategory = CATEGORY_ITEM
	ElementDefs[9].EditorShortcut = 'D'
	ElementDefs[9].Name = "Door"
	ElementDefs[10].Character = '\xe8'
	ElementDefs[10].Color = 0x0F
	ElementDefs[10].TouchProc = ElementScrollTouch
	ElementDefs[10].TickProc = ElementScrollTick
	ElementDefs[10].Pushable = true
	ElementDefs[10].Cycle = 1
	ElementDefs[10].EditorCategory = CATEGORY_ITEM
	ElementDefs[10].EditorShortcut = 'S'
	ElementDefs[10].Name = "Scroll"
	ElementDefs[10].ParamTextName = "Edit text of scroll"
	ElementDefs[12].Character = '\xfa'
	ElementDefs[12].Color = 0x0F
	ElementDefs[12].Cycle = 2
	ElementDefs[12].TickProc = ElementDuplicatorTick
	ElementDefs[12].HasDrawProc = true
	ElementDefs[12].DrawProc = ElementDuplicatorDraw
	ElementDefs[12].EditorCategory = CATEGORY_ITEM
	ElementDefs[12].EditorShortcut = 'U'
	ElementDefs[12].Name = "Duplicator"
	ElementDefs[12].ParamDirName = "Source direction?"
	ElementDefs[12].Param2Name = "Duplication rate?;SF"
	ElementDefs[6].Character = '\x9d'
	ElementDefs[6].Color = 0x06
	ElementDefs[6].VisibleInDark = true
	ElementDefs[6].TouchProc = ElementTorchTouch
	ElementDefs[6].EditorCategory = CATEGORY_ITEM
	ElementDefs[6].EditorShortcut = 'T'
	ElementDefs[6].Name = "Torch"
	ElementDefs[39].Character = '\x18'
	ElementDefs[39].Cycle = 2
	ElementDefs[39].TickProc = ElementSpinningGunTick
	ElementDefs[39].HasDrawProc = true
	ElementDefs[39].DrawProc = ElementSpinningGunDraw
	ElementDefs[39].EditorCategory = CATEGORY_CREATURE
	ElementDefs[39].EditorShortcut = 'G'
	ElementDefs[39].Name = "Spinning gun"
	ElementDefs[39].Param1Name = "Intelligence?"
	ElementDefs[39].Param2Name = "Firing rate?"
	ElementDefs[39].ParamBulletTypeName = "Firing type?"
	ElementDefs[35].Character = '\x05'
	ElementDefs[35].Color = 0x0D
	ElementDefs[35].Destructible = true
	ElementDefs[35].Pushable = true
	ElementDefs[35].Cycle = 1
	ElementDefs[35].TickProc = ElementRuffianTick
	ElementDefs[35].TouchProc = ElementDamagingTouch
	ElementDefs[35].EditorCategory = CATEGORY_CREATURE
	ElementDefs[35].EditorShortcut = 'R'
	ElementDefs[35].Name = "Ruffian"
	ElementDefs[35].Param1Name = "Intelligence?"
	ElementDefs[35].Param2Name = "Resting time?"
	ElementDefs[35].ScoreValue = 2
	ElementDefs[34].Character = '\x99'
	ElementDefs[34].Color = 0x06
	ElementDefs[34].Destructible = true
	ElementDefs[34].Pushable = true
	ElementDefs[34].Cycle = 3
	ElementDefs[34].TickProc = ElementBearTick
	ElementDefs[34].TouchProc = ElementDamagingTouch
	ElementDefs[34].EditorCategory = CATEGORY_CREATURE
	ElementDefs[34].EditorShortcut = 'B'
	ElementDefs[34].Name = "Bear"
	ElementDefs[34].CategoryName = "Creatures:"
	ElementDefs[34].Param1Name = "Sensitivity?"
	ElementDefs[34].ScoreValue = 1
	ElementDefs[37].Character = '*'
	ElementDefs[37].Color = COLOR_CHOICE_ON_BLACK
	ElementDefs[37].Destructible = false
	ElementDefs[37].Cycle = 3
	ElementDefs[37].TickProc = ElementSlimeTick
	ElementDefs[37].TouchProc = ElementSlimeTouch
	ElementDefs[37].EditorCategory = CATEGORY_CREATURE
	ElementDefs[37].EditorShortcut = 'V'
	ElementDefs[37].Name = "Slime"
	ElementDefs[37].Param2Name = "Movement speed?;FS"
	ElementDefs[38].Character = '^'
	ElementDefs[38].Color = 0x07
	ElementDefs[38].Destructible = false
	ElementDefs[38].Cycle = 3
	ElementDefs[38].TickProc = ElementSharkTick
	ElementDefs[38].EditorCategory = CATEGORY_CREATURE
	ElementDefs[38].EditorShortcut = 'Y'
	ElementDefs[38].Name = "Shark"
	ElementDefs[38].Param1Name = "Intelligence?"
	ElementDefs[16].Character = '/'
	ElementDefs[16].Cycle = 3
	ElementDefs[16].HasDrawProc = true
	ElementDefs[16].TickProc = ElementConveyorCWTick
	ElementDefs[16].DrawProc = ElementConveyorCWDraw
	ElementDefs[16].EditorCategory = CATEGORY_ITEM
	ElementDefs[16].EditorShortcut = '1'
	ElementDefs[16].Name = "Clockwise"
	ElementDefs[16].CategoryName = "Conveyors:"
	ElementDefs[17].Character = '\\'
	ElementDefs[17].Cycle = 2
	ElementDefs[17].HasDrawProc = true
	ElementDefs[17].DrawProc = ElementConveyorCCWDraw
	ElementDefs[17].TickProc = ElementConveyorCCWTick
	ElementDefs[17].EditorCategory = CATEGORY_ITEM
	ElementDefs[17].EditorShortcut = '2'
	ElementDefs[17].Name = "Counter"
	ElementDefs[21].Character = '\xdb'
	ElementDefs[21].EditorCategory = CATEGORY_TERRAIN
	ElementDefs[21].CategoryName = "Walls:"
	ElementDefs[21].EditorShortcut = 'S'
	ElementDefs[21].Name = "Solid"
	ElementDefs[22].Character = '\xb2'
	ElementDefs[22].EditorCategory = CATEGORY_TERRAIN
	ElementDefs[22].EditorShortcut = 'N'
	ElementDefs[22].Name = "Normal"
	ElementDefs[31].Character = '\xce'
	ElementDefs[31].HasDrawProc = true
	ElementDefs[31].DrawProc = ElementLineDraw
	ElementDefs[31].Name = "Line"
	ElementDefs[43].Character = '\xba'
	ElementDefs[33].Character = '\xcd'
	ElementDefs[32].Character = '*'
	ElementDefs[32].Color = 0x0A
	ElementDefs[32].EditorCategory = CATEGORY_TERRAIN
	ElementDefs[32].EditorShortcut = 'R'
	ElementDefs[32].Name = "Ricochet"
	ElementDefs[23].Character = '\xb1'
	ElementDefs[23].Destructible = false
	ElementDefs[23].EditorCategory = CATEGORY_TERRAIN
	ElementDefs[23].EditorShortcut = 'B'
	ElementDefs[23].Name = "Breakable"
	ElementDefs[24].Character = '\xfe'
	ElementDefs[24].Pushable = true
	ElementDefs[24].TouchProc = ElementPushableTouch
	ElementDefs[24].EditorCategory = CATEGORY_TERRAIN
	ElementDefs[24].EditorShortcut = 'O'
	ElementDefs[24].Name = "Boulder"
	ElementDefs[25].Character = '\x12'
	ElementDefs[25].TouchProc = ElementPushableTouch
	ElementDefs[25].EditorCategory = CATEGORY_TERRAIN
	ElementDefs[25].EditorShortcut = '1'
	ElementDefs[25].Name = "Slider (NS)"
	ElementDefs[26].Character = '\x1d'
	ElementDefs[26].TouchProc = ElementPushableTouch
	ElementDefs[26].EditorCategory = CATEGORY_TERRAIN
	ElementDefs[26].EditorShortcut = '2'
	ElementDefs[26].Name = "Slider (EW)"
	ElementDefs[30].Character = '\xc5'
	ElementDefs[30].TouchProc = ElementTransporterTouch
	ElementDefs[30].HasDrawProc = true
	ElementDefs[30].DrawProc = ElementTransporterDraw
	ElementDefs[30].Cycle = 2
	ElementDefs[30].TickProc = ElementTransporterTick
	ElementDefs[30].EditorCategory = CATEGORY_TERRAIN
	ElementDefs[30].EditorShortcut = 'T'
	ElementDefs[30].Name = "Transporter"
	ElementDefs[30].ParamDirName = "Direction?"
	ElementDefs[40].Character = '\x10'
	ElementDefs[40].Color = COLOR_CHOICE_ON_BLACK
	ElementDefs[40].HasDrawProc = true
	ElementDefs[40].DrawProc = ElementPusherDraw
	ElementDefs[40].Cycle = 4
	ElementDefs[40].TickProc = ElementPusherTick
	ElementDefs[40].EditorCategory = CATEGORY_CREATURE
	ElementDefs[40].EditorShortcut = 'P'
	ElementDefs[40].Name = "Pusher"
	ElementDefs[40].ParamDirName = "Push direction?"
	ElementDefs[13].Character = '\x0b'
	ElementDefs[13].HasDrawProc = true
	ElementDefs[13].DrawProc = ElementBombDraw
	ElementDefs[13].Pushable = true
	ElementDefs[13].Cycle = 6
	ElementDefs[13].TickProc = ElementBombTick
	ElementDefs[13].TouchProc = ElementBombTouch
	ElementDefs[13].EditorCategory = CATEGORY_ITEM
	ElementDefs[13].EditorShortcut = 'B'
	ElementDefs[13].Name = "Bomb"
	ElementDefs[14].Character = '\x7f'
	ElementDefs[14].Color = 0x05
	ElementDefs[14].TouchProc = ElementEnergizerTouch
	ElementDefs[14].EditorCategory = CATEGORY_ITEM
	ElementDefs[14].EditorShortcut = 'E'
	ElementDefs[14].Name = "Energizer"
	ElementDefs[29].Character = '\xce'
	ElementDefs[29].Cycle = 1
	ElementDefs[29].TickProc = ElementBlinkWallTick
	ElementDefs[29].HasDrawProc = true
	ElementDefs[29].DrawProc = ElementBlinkWallDraw
	ElementDefs[29].EditorCategory = CATEGORY_TERRAIN
	ElementDefs[29].EditorShortcut = 'L'
	ElementDefs[29].Name = "Blink wall"
	ElementDefs[29].Param1Name = "Starting time"
	ElementDefs[29].Param2Name = "Period"
	ElementDefs[29].ParamDirName = "Wall direction"
	ElementDefs[27].Character = '\xb2'
	ElementDefs[27].EditorCategory = CATEGORY_TERRAIN
	ElementDefs[27].PlaceableOnTop = true
	ElementDefs[27].Walkable = true
	ElementDefs[27].TouchProc = ElementFakeTouch
	ElementDefs[27].EditorShortcut = 'A'
	ElementDefs[27].Name = "Fake"
	ElementDefs[28].Character = ' '
	ElementDefs[28].EditorCategory = CATEGORY_TERRAIN
	ElementDefs[28].TouchProc = ElementInvisibleTouch
	ElementDefs[28].EditorShortcut = 'I'
	ElementDefs[28].Name = "Invisible"
	ElementDefs[36].Character = '\x02'
	ElementDefs[36].EditorCategory = CATEGORY_CREATURE
	ElementDefs[36].Cycle = 3
	ElementDefs[36].HasDrawProc = true
	ElementDefs[36].DrawProc = ElementObjectDraw
	ElementDefs[36].TickProc = ElementObjectTick
	ElementDefs[36].TouchProc = ElementObjectTouch
	ElementDefs[36].EditorShortcut = 'O'
	ElementDefs[36].Name = "Object"
	ElementDefs[36].Param1Name = "Character?"
	ElementDefs[36].ParamTextName = "Edit Program"
	ElementDefs[2].TickProc = ElementMessageTimerTick
	ElementDefs[1].TouchProc = ElementBoardEdgeTouch
	EditorPatternCount = 5
	EditorPatterns[0] = E_SOLID
	EditorPatterns[1] = E_NORMAL
	EditorPatterns[2] = E_BREAKABLE
	EditorPatterns[3] = E_EMPTY
	EditorPatterns[4] = E_LINE
}

func InitElementsEditor() {
	InitElementDefs()
	ElementDefs[28].Character = '\xb0'
	ElementDefs[28].Color = COLOR_CHOICE_ON_BLACK
	ForceDarknessOff = true
}

func InitElementsGame() {
	InitElementDefs()
	ForceDarknessOff = false
}

func InitEditorStatSettings() {
	var i int16
	PlayerDirX = 0
	PlayerDirY = 0
	for i = 0; i <= MAX_ELEMENT; i++ {
		setting := &World.EditorStatSettings[i]
		setting.P1 = 4
		setting.P2 = 4
		setting.P3 = 0
		setting.StepX = 0
		setting.StepY = -1
	}
	World.EditorStatSettings[E_OBJECT].P1 = 1
	World.EditorStatSettings[E_BEAR].P1 = 8
}
