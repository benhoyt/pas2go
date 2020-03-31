package main // unit: GameVars

const (
	MAX_STAT         = 150
	MAX_ELEMENT      = 53
	MAX_BOARD        = 100
	MAX_FLAG         = 10
	BOARD_WIDTH      = 60
	BOARD_HEIGHT     = 25
	HIGH_SCORE_COUNT = 30
	TORCH_DURATION   = 200
	TORCH_DX         = 8
	TORCH_DY         = 5
	TORCH_DIST_SQR   = 50
)

type (
	TString50 string
	TCoord    struct {
		X int16
		Y int16
	}
	TTile struct {
		Element byte
		Color   byte
	}
	TElementDrawProc  func(x, y int16, ch *byte)
	TElementTickProc  func(statId int16)
	TElementTouchProc func(x, y int16, sourceStatId int16, deltaX, deltaY *int16)
	TElementDef       struct {
		Character           byte
		Color               byte
		Destructible        bool
		Pushable            bool
		VisibleInDark       bool
		PlaceableOnTop      bool
		Walkable            bool
		HasDrawProc         bool
		DrawProc            TElementDrawProc
		Cycle               int16
		TickProc            TElementTickProc
		TouchProc           TElementTouchProc
		EditorCategory      int16
		EditorShortcut      byte
		Name                string
		CategoryName        string
		Param1Name          string
		Param2Name          string
		ParamBulletTypeName string
		ParamBoardName      string
		ParamDirName        string
		ParamTextName       string
		ScoreValue          int16
	}
	TStat struct {
		X, Y         byte
		StepX, StepY int16
		Cycle        int16
		P1, P2, P3   byte
		Follower     int16
		Leader       int16
		Under        TTile
		Data         *string
		DataPos      int16
		DataLen      int16
		unk1, unk2   uintptr
	}
	TRleTile struct {
		Count byte
		Tile  TTile
	}
	TBoardInfo struct {
		MaxShots          byte
		IsDark            bool
		NeighborBoards    [3 - 0 + 1]byte
		ReenterWhenZapped bool
		Message           string
		StartPlayerX      byte
		StartPlayerY      byte
		TimeLimitSec      int16
		unk1              [85 - 70 + 1]byte
	}
	TWorldInfo struct {
		Ammo           int16
		Gems           int16
		Keys           [7 - 1 + 1]bool
		Health         int16
		CurrentBoard   int16
		Torches        int16
		TorchTicks     int16
		EnergizerTicks int16
		unk1           int16
		Score          int16
		Name           string
		Flags          [MAX_FLAG - 1 + 1]string
		BoardTimeSec   int16
		BoardTimeHsec  int16
		IsSave         bool
		unkPad         [13 - 0 + 1]byte
	}
	TEditorStatSetting struct {
		P1, P2, P3   byte
		StepX, StepY int16
	}
	TBoard struct {
		Name      TString50
		Tiles     [BOARD_WIDTH + 1 - 0 + 1][BOARD_HEIGHT + 1 - 0 + 1]TTile
		StatCount int16
		Stats     [MAX_STAT + 1 - 0 + 1]TStat
		Info      TBoardInfo
	}
	TWorld struct {
		BoardCount         int16
		BoardData          [MAX_BOARD - 0 + 1]uintptr
		BoardLen           [MAX_BOARD - 0 + 1]int16
		Info               TWorldInfo
		EditorStatSettings [MAX_ELEMENT - 0 + 1]TEditorStatSetting
	}
	THighScoreEntry struct {
		Name  string
		Score int16
	}
	THighScoreList [HIGH_SCORE_COUNT - 1 + 1]THighScoreEntry
	TIoTmpBuf      [19999 - 0 + 1]byte
)

var (
	PlayerDirX                  int16
	PlayerDirY                  int16
	unkVar_0476                 int16
	unkVar_0478                 int16
	TransitionTable             [80*25 - 1 + 1]TCoord
	LoadedGameFileName          TString50
	SavedGameFileName           TString50
	SavedBoardFileName          TString50
	StartupWorldFileName        TString50
	Board                       TBoard
	World                       TWorld
	MessageAmmoNotShown         bool
	MessageOutOfAmmoNotShown    bool
	MessageNoShootingNotShown   bool
	MessageTorchNotShown        bool
	MessageOutOfTorchesNotShown bool
	MessageRoomNotDarkNotShown  bool
	MessageHintTorchNotShown    bool
	MessageForestNotShown       bool
	MessageFakeNotShown         bool
	MessageGemNotShown          bool
	MessageEnergizerNotShown    bool
	unkVar_4ABA                 [14 - 0 + 1]byte
	GameTitleExitRequested      bool
	GamePlayExitRequested       bool
	GameStateElement            int16
	ReturnBoardId               int16
	TransitionTableSize         int16
	TickSpeed                   byte
	IoTmpBuf                    *TIoTmpBuf
	ElementDefs                 [MAX_ELEMENT - 0 + 1]TElementDef
	EditorPatternCount          int16
	EditorPatterns              [10 - 1 + 1]byte
	TickTimeDuration            int16
	CurrentTick                 int16
	CurrentStatTicked           int16
	GamePaused                  bool
	TickTimeCounter             int16
	ForceDarknessOff            bool
	InitialTextAttr             byte
	OopChar                     byte
	OopWord                     string
	OopValue                    int16
	DebugEnabled                bool
	HighScoreList               THighScoreList
	ConfigRegistration          string
	ConfigWorldFile             TString50
	EditorEnabled               bool
	GameVersion                 TString50
	ParsingConfigFile           bool
	ResetConfig                 bool
	JustStarted                 bool
	WorldFileDescCount          int16
	WorldFileDescKeys           [10 - 1 + 1]TString50
	WorldFileDescValues         [10 - 1 + 1]TString50
)

const (
	E_EMPTY                = 0
	E_BOARD_EDGE           = 1
	E_MESSAGE_TIMER        = 2
	E_MONITOR              = 3
	E_PLAYER               = 4
	E_AMMO                 = 5
	E_TORCH                = 6
	E_GEM                  = 7
	E_KEY                  = 8
	E_DOOR                 = 9
	E_SCROLL               = 10
	E_PASSAGE              = 11
	E_DUPLICATOR           = 12
	E_BOMB                 = 13
	E_ENERGIZER            = 14
	E_STAR                 = 15
	E_CONVEYOR_CW          = 16
	E_CONVEYOR_CCW         = 17
	E_BULLET               = 18
	E_WATER                = 19
	E_FOREST               = 20
	E_SOLID                = 21
	E_NORMAL               = 22
	E_BREAKABLE            = 23
	E_BOULDER              = 24
	E_SLIDER_NS            = 25
	E_SLIDER_EW            = 26
	E_FAKE                 = 27
	E_INVISIBLE            = 28
	E_BLINK_WALL           = 29
	E_TRANSPORTER          = 30
	E_LINE                 = 31
	E_RICOCHET             = 32
	E_BLINK_RAY_EW         = 33
	E_BEAR                 = 34
	E_RUFFIAN              = 35
	E_OBJECT               = 36
	E_SLIME                = 37
	E_SHARK                = 38
	E_SPINNING_GUN         = 39
	E_PUSHER               = 40
	E_LION                 = 41
	E_TIGER                = 42
	E_BLINK_RAY_NS         = 43
	E_CENTIPEDE_HEAD       = 44
	E_CENTIPEDE_SEGMENT    = 45
	E_TEXT_BLUE            = 47
	E_TEXT_GREEN           = 48
	E_TEXT_CYAN            = 49
	E_TEXT_RED             = 50
	E_TEXT_PURPLE          = 51
	E_TEXT_YELLOW          = 52
	E_TEXT_WHITE           = 53
	E_TEXT_MIN             = E_TEXT_BLUE
	CATEGORY_ITEM          = 1
	CATEGORY_CREATURE      = 2
	CATEGORY_TERRAIN       = 3
	COLOR_SPECIAL_MIN      = 0xF0
	COLOR_CHOICE_ON_BLACK  = 0xFF
	COLOR_WHITE_ON_CHOICE  = 0xFE
	COLOR_CHOICE_ON_CHOICE = 0xFD
	SHOT_SOURCE_PLAYER     = 0
	SHOT_SOURCE_ENEMY      = 1
)

func init() {
}
