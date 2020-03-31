package main // unit: Sounds

type (
	TDrumData struct {
		Len  int16
		Data [255 - 1 + 1]word
	}
)

var (
	SoundEnabled            bool
	SoundBlockQueueing      bool
	SoundCurrentPriority    int16
	SoundFreqTable          [255 - 1 + 1]word
	SoundDurationMultiplier byte
	SoundDurationCounter    byte
	SoundBuffer             string
	SoundNewVector          uintptr
	SoundOldVector          uintptr
	SoundBufferPos          int16
	SoundIsPlaying          bool
	SoundTimeCheckCounter   int16
	UseSystemTimeForElapsed bool
	TimerTicks              word
	SoundTimeCheckHsec      int16
	SoundDrumTable          [9 - 0 + 1]TDrumData
)

// implementation uses: Crt, Dos

func SoundQueue(priority int16, pattern string) {
	if !SoundBlockQueueing && (!SoundIsPlaying || (((priority >= SoundCurrentPriority) && (SoundCurrentPriority != -1)) || (priority == -1))) {
		if (priority >= 0) || !SoundIsPlaying {
			SoundCurrentPriority = priority
			SoundBuffer = pattern
			SoundBufferPos = 1
			SoundDurationCounter = 1
		} else {
			SoundBuffer = Copy(SoundBuffer, SoundBufferPos, Length(SoundBuffer)-SoundBufferPos+1)
			SoundBufferPos = 1
			if (Length(SoundBuffer) + Length(pattern)) < 255 {
				SoundBuffer = SoundBuffer + pattern
			}
		}
		SoundIsPlaying = true
	}
}

func SoundClearQueue() {
	SoundBuffer = ""
	SoundIsPlaying = false
	NoSound()
}

func SoundInitFreqTable() {
	var (
		octave, note                    int16
		freqC1, noteStep, noteBase, ln2 float64
	)
	freqC1 = 32.0
	ln2 = Ln(2.0)
	noteStep = Exp(ln2 / 12.0)
	for octave := 1; octave <= 15; octave++ {
		noteBase = Exp(octave*ln2) * freqC1
		for note := 0; note <= 11; note++ {
			SoundFreqTable[octave*16+note] = Trunc(noteBase)
			noteBase = noteBase * noteStep
		}
	}
}

func SoundInitDrumTable() {
	var (
		i int16
	)
	SoundDrumTable[0].Len = 1
	SoundDrumTable[0].Data[1] = 3200
	for i := 1; i <= 9; i++ {
		SoundDrumTable[i].Len = 14
	}
	for i := 1; i <= 14; i++ {
		SoundDrumTable[1].Data[i] = i*100 + 1000
	}
	for i := 1; i <= 16; i++ {
		SoundDrumTable[2].Data[i] = (i%2)*1600 + 1600 + (i%4)*1600
	}
	for i := 1; i <= 14; i++ {
		SoundDrumTable[4].Data[i] = Random(5000) + 500
	}
	for i := 1; i <= 8; i++ {
		SoundDrumTable[5].Data[i*2-1] = 1600
		SoundDrumTable[5].Data[i*2] = Random(1600) + 800
	}
	for i := 1; i <= 14; i++ {
		SoundDrumTable[6].Data[i] = ((i % 2) * 880) + 880 + ((i % 3) * 440)
	}
	for i := 1; i <= 14; i++ {
		SoundDrumTable[7].Data[i] = 700 - (i * 12)
	}
	for i := 1; i <= 14; i++ {
		SoundDrumTable[8].Data[i] = (i*20 + 1200) - Random(i*40)
	}
	for i := 1; i <= 14; i++ {
		SoundDrumTable[9].Data[i] = Random(440) + 220
	}
}

func SoundPlayDrum(drum *TDrumData) {
	var (
		i int16
	)
	for i := 1; i <= drum.Len; i++ {
		Sound(drum.Data[i])
		Delay(1)
	}
	NoSound()
}

func SoundCheckTimeIntr() {
	var (
		hour, minute, sec, hSec word
	)
	GetTime(hour, minute, sec, hSec)
	if (SoundTimeCheckHsec != 0) && (int16(hSec) != SoundTimeCheckHsec) {
		SoundTimeCheckCounter = 0
		UseSystemTimeForElapsed = true
	}
	SoundTimeCheckHsec = int16(hSec)
}

func SoundHasTimeElapsed(counter *int16, duration int16) (SoundHasTimeElapsed bool) {
	var (
		hour, minute, sec, hSec word
		hSecsDiff               word
		hSecsTotal              int16
	)
	if (SoundTimeCheckCounter > 0) && ((SoundTimeCheckCounter % 2) == 1) {
		SoundTimeCheckCounter = SoundTimeCheckCounter - 1
		SoundCheckTimeIntr()
	}
	if UseSystemTimeForElapsed {
		GetTime(hour, minute, sec, hSec)
		hSecsTotal = sec*100 + hSec
		hSecsDiff = Word((hSecsTotal-counter)+6000) % 6000
	} else {
		hSecsTotal = TimerTicks * 6
		hSecsDiff = hSecsTotal - counter
	}
	if hSecsDiff >= duration {
		SoundHasTimeElapsed = true
		counter = hSecsTotal
	} else {
		SoundHasTimeElapsed = false
	}
	return
}

func SoundTimerHandler() {
	Inc(TimerTicks)
	if (SoundTimeCheckCounter > 0) && ((SoundTimeCheckCounter % 2) == 0) {
		SoundTimeCheckCounter = SoundTimeCheckCounter - 1
	}
	if !SoundEnabled {
		SoundIsPlaying = false
		NoSound()
	} else if SoundIsPlaying {
		Dec(SoundDurationCounter)
		if SoundDurationCounter <= 0 {
			NoSound()
			if SoundBufferPos >= Length(SoundBuffer) {
				NoSound()
				SoundIsPlaying = false
			} else {
				if SoundBuffer[SoundBufferPos] == '\x00' {
					NoSound()
				} else if SoundBuffer[SoundBufferPos] < 'ð' {
					Sound(SoundFreqTable[Ord(SoundBuffer[SoundBufferPos])])
				} else {
					SoundPlayDrum(SoundDrumTable[Ord(SoundBuffer[SoundBufferPos])-240])
				}

				Inc(SoundBufferPos)
				SoundDurationCounter = SoundDurationMultiplier * Ord(SoundBuffer[SoundBufferPos])
				Inc(SoundBufferPos)
			}
		}
	}

}

func SoundUninstall() {
	SetIntVec(0x1C, SoundOldVector)
}

func SoundParse(input string) (SoundParse string) {
	var (
		noteOctave   int16
		noteDuration int16
		output       string
		noteTone     int16
	)
	AdvanceInput := func() {
		input = Copy(input, 2, Length(input)-1)
	}

	output = ""
	noteOctave = 3
	noteDuration = 1
	for Length(input) != 0 {
		noteTone = -1
		switch UpCase(input[1]) {
		case 'T':
			noteDuration = 1
			AdvanceInput()
		case 'S':
			noteDuration = 2
			AdvanceInput()
		case 'I':
			noteDuration = 4
			AdvanceInput()
		case 'Q':
			noteDuration = 8
			AdvanceInput()
		case 'H':
			noteDuration = 16
			AdvanceInput()
		case 'W':
			noteDuration = 32
			AdvanceInput()
		case '.':
			noteDuration = (noteDuration * 3) / 2
			AdvanceInput()
		case '3':
			noteDuration = noteDuration / 3
			AdvanceInput()
		case '+':
			if noteOctave < 6 {
				noteOctave = noteOctave + 1
			}
			AdvanceInput()
		case '-':
			if noteOctave > 1 {
				noteOctave = noteOctave - 1
			}
			AdvanceInput()
		case 'A', 'B', 'C', 'D', 'E', 'F', 'G':
			switch UpCase(input[1]) {
			case 'C':
				noteTone = 0
				AdvanceInput()
			case 'D':
				noteTone = 2
				AdvanceInput()
			case 'E':
				noteTone = 4
				AdvanceInput()
			case 'F':
				noteTone = 5
				AdvanceInput()
			case 'G':
				noteTone = 7
				AdvanceInput()
			case 'A':
				noteTone = 9
				AdvanceInput()
			case 'B':
				noteTone = 11
				AdvanceInput()
			}
			switch UpCase(input[1]) {
			case '!':
				noteTone = noteTone - 1
				AdvanceInput()
			case '#':
				noteTone = noteTone + 1
				AdvanceInput()
			}
			output = output + Chr(noteOctave*0x10+noteTone) + Chr(noteDuration)
		case 'X':
			output = output + '\x00' + Chr(noteDuration)
			AdvanceInput()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			output = output + Chr(Ord(input[1])+0xF0-Ord('0')) + Chr(noteDuration)
			AdvanceInput()
		default:
			AdvanceInput()
		}
	}
	SoundParse = output
	return
}

func init() {
	SoundInitFreqTable()
	SoundInitDrumTable()
	SoundTimeCheckCounter = 36
	UseSystemTimeForElapsed = false
	TimerTicks = 0
	SoundTimeCheckHsec = 0
	SoundEnabled = true
	SoundBlockQueueing = false
	SoundClearQueue()
	SoundDurationMultiplier = 1
	SoundIsPlaying = false
	TimerTicks = 0
	SoundNewVector = *SoundTimerHandler
	GetIntVec(0x1C, SoundOldVector)
	SetIntVec(0x1C, SoundNewVector)
}
