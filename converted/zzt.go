package main

// uses: Crt, Dos, Video, Keys, Sounds, Input, TxtWind, GameVars, Elements, Editor, Oop, Game

func ParseArguments() {
	var (
		i    int16
		pArg string
	)
	for i = 1; i <= ParamCount; i++ {
		pArg = ParamStr(i)
		if pArg[1] == '/' {
			switch UpCase(pArg[2]) {
			case 'T':
				SoundTimeCheckCounter = 0
				UseSystemTimeForElapsed = false
			case 'R':
				ResetConfig = true
			}
		} else {
			StartupWorldFileName = pArg
			if (Length(StartupWorldFileName) > 4) && (StartupWorldFileName[Length(StartupWorldFileName)-3] == '.') {
				StartupWorldFileName = Copy(StartupWorldFileName, 1, Length(StartupWorldFileName)-4)
			}
		}
	}
}

func GameConfigure() {
	var (
		unk1                          int16
		joystickEnabled, mouseEnabled bool
		cfgFile                       text
	)
	ParsingConfigFile = true
	EditorEnabled = true
	ConfigRegistration = ""
	ConfigWorldFile = ""
	GameVersion = "3.2"
	Assign(cfgFile, "zzt.cfg")
	Reset(cfgFile)
	if IOResult == 0 {
		Readln(cfgFile, ConfigWorldFile)
		Readln(cfgFile, ConfigRegistration)
	}
	if ConfigWorldFile[1] == '*' {
		EditorEnabled = false
		ConfigWorldFile = Copy(ConfigWorldFile, 2, Length(ConfigWorldFile)-1)
	}
	if Length(ConfigWorldFile) != 0 {
		StartupWorldFileName = ConfigWorldFile
	}
	InputInitDevices()
	joystickEnabled = InputJoystickEnabled
	mouseEnabled = InputMouseEnabled
	ParsingConfigFile = false
	Window(1, 1, 80, 25)
	TextBackground(Black)
	ClrScr()
	TextColor(White)
	TextColor(White)
	Writeln()
	Writeln("                                 <=-  ZZT  -=>")
	TextColor(Yellow)
	if Length(ConfigRegistration) == 0 {
		Writeln("                             Shareware version 3.2")
	} else {
		Writeln("                                  Version  3.2")
	}
	Writeln("                            Created by Tim Sweeney")
	GotoXY(1, 7)
	TextColor(Blue)
	Write("================================================================================")
	GotoXY(1, 24)
	Write("================================================================================")
	TextColor(White)
	GotoXY(30, 7)
	Write(" Game Configuration ")
	GotoXY(1, 25)
	Write(" Copyright (c) 1991 Epic MegaGames                         Press ... to abort")
	TextColor(Black)
	TextBackground(LightGray)
	GotoXY(66, 25)
	Write("ESC")
	Window(1, 8, 80, 23)
	TextColor(Yellow)
	TextBackground(Black)
	ClrScr()
	TextColor(Yellow)
	if !InputConfigure {
		GameTitleExitRequested = true
	} else {
		TextColor(LightGreen)
		if !VideoConfigure {
			GameTitleExitRequested = true
		}
	}
	Window(1, 1, 80, 25)
}

func main() {
	WorldFileDescCount = 7
	WorldFileDescKeys[1+1] = "TOWN"
	WorldFileDescValues[1+1] = "TOWN       The Town of ZZT"
	WorldFileDescKeys[2+1] = "DEMO"
	WorldFileDescValues[2+1] = "DEMO       Demo of the ZZT World Editor"
	WorldFileDescKeys[3+1] = "CAVES"
	WorldFileDescValues[3+1] = "CAVES      The Caves of ZZT"
	WorldFileDescKeys[4+1] = "DUNGEONS"
	WorldFileDescValues[4+1] = "DUNGEONS   The Dungeons of ZZT"
	WorldFileDescKeys[5+1] = "CITY"
	WorldFileDescValues[5+1] = "CITY       Underground City of ZZT"
	WorldFileDescKeys[6+1] = "BEST"
	WorldFileDescValues[6+1] = "BEST       The Best of ZZT"
	WorldFileDescKeys[7+1] = "TOUR"
	WorldFileDescValues[7+1] = "TOUR       Guided Tour ZZT's Other Worlds"
	Randomize()
	SetCBreak(false)
	InitialTextAttr = TextAttr
	StartupWorldFileName = "TOWN"
	ResourceDataFileName = "ZZT.DAT"
	ResetConfig = false
	GameTitleExitRequested = false
	GameConfigure()
	ParseArguments()
	if !GameTitleExitRequested {
		VideoInstall(80, Blue)
		OrderPrintId = *GameVersion
		TextWindowInit(5, 3, 50, 18)
		New(IoTmpBuf)
		VideoHideCursor()
		ClrScr()
		TickSpeed = 4
		DebugEnabled = false
		SavedGameFileName = "SAVED"
		SavedBoardFileName = "TEMP"
		GenerateTransitionTable()
		WorldCreate()
		GameTitleLoop()
		Dispose(IoTmpBuf)
	}
	SoundUninstall()
	SoundClearQueue()
	VideoUninstall()
	Port[PORT_CGA_PALETTE] = 0
	TextAttr = InitialTextAttr
	ClrScr()
	if Length(ConfigRegistration) == 0 {
		GamePrintRegisterMessage()
	} else {
		Writeln()
		Writeln("  Registered version -- Thank you for playing ZZT.")
		Writeln()
	}
	VideoShowCursor()
}
