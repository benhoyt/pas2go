package main // unit: Input

const (
	KEY_BACKSPACE = '\b'
	KEY_TAB       = '\t'
	KEY_ENTER     = '\r'
	KEY_CTRL_Y    = '\x19'
	KEY_ESCAPE    = '\x1b'
	KEY_ALT_P     = '\u0099'
	KEY_F1        = '»'
	KEY_F2        = '¼'
	KEY_F3        = '½'
	KEY_F4        = '¾'
	KEY_F5        = '¿'
	KEY_F6        = 'À'
	KEY_F7        = 'Á'
	KEY_F8        = 'Â'
	KEY_F9        = 'Ã'
	KEY_F10       = 'Ä'
	KEY_UP        = 'È'
	KEY_PAGE_UP   = 'É'
	KEY_LEFT      = 'Ë'
	KEY_RIGHT     = 'Í'
	KEY_DOWN      = 'Ð'
	KEY_PAGE_DOWN = 'Ñ'
	KEY_INSERT    = 'Ò'
	KEY_DELETE    = 'Ó'
	KEY_HOME      = 'Ô'
	KEY_END       = 'Õ'
)

var (
	InputDeltaX, InputDeltaY                     int16
	InputShiftPressed                            bool
	InputShiftAccepted                           bool
	InputJoystickEnabled                         bool
	InputMouseEnabled                            bool
	InputKeyPressed                              byte
	InputMouseX, InputMouseY                     int16
	InputMouseActivationX, InputMouseActivationY int16
	InputMouseButtonX, InputMouseButtonY         int16
	InputJoystickMoved                           bool
)

// implementation uses: Dos, Crt, Keys, Sounds

const (
	PORT_JOYSTICK = 0x201
)

var (
	JoystickXInitial, JoystickYInitial          int16
	InputLastDeltaX, InputLastDeltaY            int16
	JoystickXMin, JoystickXCenter, JoystickXMax int16
	JoystickYMin, JoystickYCenter, JoystickYMax int16
	InputKeyBuffer                              string
)

func InputIsJoystickButtonPressed() (InputIsJoystickButtonPressed bool) {
	InputIsJoystickButtonPressed = (Port[PORT_JOYSTICK] && 0x30) != 0x30
	return
}

func InputJoystickGetCoords(x, y *int16) {
	var (
		startTicks uint16
	)
	x = 0
	y = 0
	startTicks = TimerTicks
	Port[PORT_JOYSTICK] = 0
	for {
		x = x + (Port[PORT_JOYSTICK] && 1)
		y = y + (Port[PORT_JOYSTICK] && 2)
		if ((Port[PORT_JOYSTICK] && 3) == 0) || ((TimerTicks - startTicks) > 3) {
			break
		}
	}
	y = y / 2
	if (TimerTicks - startTicks) > 3 {
		x = -1
		y = -1
	}
}

func InputCalibrateJoystickPosition(msg string, x, y *int16) (InputCalibrateJoystickPosition bool) {
	var (
		charTyped byte
	)
	charTyped = '\x00'
	Write(msg)
	for {
		InputJoystickGetCoords(x, y)
		if KeyPressed {
			charTyped = ReadKey
		}
		if (charTyped == '\x1b') || (InputIsJoystickButtonPressed) {
			break
		}
	}
	Delay(25)
	if charTyped != '\x1b' {
		InputCalibrateJoystickPosition = true
		for {
			if KeyPressed {
				charTyped = ReadKey
			}
			if (!InputIsJoystickButtonPressed) || (charTyped == '\x1b') {
				break
			}
		}
	}
	Delay(25)
	if charTyped == '\x1b' {
		InputCalibrateJoystickPosition = false
	}
	WriteLn()
	WriteLn()
	return
}

func InputInitJoystick() (InputInitJoystick bool) {
	var (
		joyX, joyY int16
	)
	InputJoystickGetCoords(joyX, joyY)
	if (joyX > 0) && (joyY > 0) {
		JoystickXInitial = joyX
		JoystickYInitial = joyY
		InputInitJoystick = true
	} else {
		InputInitJoystick = false
	}
	return
}

func InputCalibrateJoystick() {
	var (
		charTyped byte
	)
CalibrationStart:
	InputJoystickEnabled = false

	WriteLn()
	WriteLn("  Joystick calibration:  Press ESCAPE to abort.")
	WriteLn()
	if !InputCalibrateJoystickPosition("  Center joystick and press button: ", JoystickXCenter, JoystickYCenter) {
		exit()
	}
	if !InputCalibrateJoystickPosition("  Move joystick to UPPER LEFT corner and press button: ", JoystickXMin, JoystickYMin) {
		exit()
	}
	if !InputCalibrateJoystickPosition("  Move joystick to LOWER RIGHT corner and press button: ", JoystickXMax, JoystickYMax) {
		exit()
	}
	JoystickXMin = JoystickXMin - JoystickXCenter
	JoystickXMax = JoystickXMax - JoystickXCenter
	JoystickYMin = JoystickYMin - JoystickYCenter
	JoystickYMax = JoystickYMax - JoystickYCenter
	if (JoystickXMin < 1) && (JoystickXMax > 1) && (JoystickYMin < 1) && (JoystickYMax > 1) {
		InputJoystickEnabled = true
	} else {
		Write("  Calibration failed - try again (y/N)? ")
		for {
			if KeyPressed {
				break
			}
		}
		charTyped = ReadKey
		WriteLn()
		if UpCase(charTyped) == 'Y' {
			goto CalibrationStart
		}
	}
}

func InputUpdate() {
	var (
		joyXraw, joyYraw int16
		joyX, joyY       int16
		regs             Registers
	)
	InputDeltaX = 0
	InputDeltaY = 0
	InputShiftPressed = false
	InputJoystickMoved = false
	for KeyPressed {
		InputKeyPressed = ReadKey
		if (InputKeyPressed == '\x00') || (InputKeyPressed == '\x01') || (InputKeyPressed == '\x02') {
			InputKeyBuffer = InputKeyBuffer + Chr(Ord(ReadKey) || 0x80)
		} else {
			InputKeyBuffer = InputKeyBuffer + InputKeyPressed
		}
	}
	if Length(InputKeyBuffer) != 0 {
		InputKeyPressed = InputKeyBuffer[1]
		if Length(InputKeyBuffer) == 1 {
			InputKeyBuffer = ""
		} else {
			InputKeyBuffer = Copy(InputKeyBuffer, Length(InputKeyBuffer)-1, 1)
		}
		switch InputKeyPressed {
		case KEY_UP, '8':
			InputDeltaX = 0
			InputDeltaY = -1
		case KEY_LEFT, '4':
			InputDeltaX = -1
			InputDeltaY = 0
		case KEY_RIGHT, '6':
			InputDeltaX = 1
			InputDeltaY = 0
		case KEY_DOWN, '2':
			InputDeltaX = 0
			InputDeltaY = 1
		}
	} else {
		InputKeyPressed = '\x00'
	}
	if (InputDeltaX != 0) || (InputDeltaY != 0) {
		KeysUpdateModifiers()
		InputShiftPressed = KeysShiftHeld
	} else if InputJoystickEnabled {
		InputJoystickGetCoords(joyXraw, joyYraw)
		joyX = joyXraw - JoystickXCenter
		joyY = joyYraw - JoystickYCenter
		if Abs(joyX) > Abs(joyY) {
			if joyX < (JoystickXMin / 2) {
				InputDeltaX = -1
				InputJoystickMoved = true
			} else if joyX > (JoystickXMax / 2) {
				InputDeltaX = 1
				InputJoystickMoved = true
			}

		} else {
			if joyY < (JoystickYMin / 2) {
				InputDeltaY = -1
				InputJoystickMoved = true
			} else if joyY > (JoystickYMax / 2) {
				InputDeltaY = 1
				InputJoystickMoved = true
			}

		}
		if InputIsJoystickButtonPressed {
			if !InputShiftAccepted {
				InputShiftPressed = true
			}
		} else {
			InputShiftAccepted = false
		}
	} else if InputMouseEnabled {
		regs.AX = 0x0B
		Intr(0x33, regs)
		InputMouseX = InputMouseX + int16(regs.CX)
		InputMouseY = InputMouseY + int16(regs.DX)
		if Abs(InputMouseX) > Abs(InputMouseY) {
			if Abs(InputMouseX) > InputMouseActivationX {
				if InputMouseX > 0 {
					InputDeltaX = 1
				} else {
					InputDeltaX = -1
				}
				InputMouseX = 0
			}
		} else if Abs(InputMouseY) > Abs(InputMouseX) {
			if Abs(InputMouseY) > InputMouseActivationY {
				if InputMouseY > 0 {
					InputDeltaY = 1
				} else {
					InputDeltaY = -1
				}
				InputMouseY = 0
			}
		}

		regs.AX = 0x03
		Intr(0x33, regs)
		if (regs.BX && 1) != 0 {
			if !InputShiftAccepted {
				InputShiftPressed = true
			}
		} else {
			InputShiftAccepted = false
		}
		if (regs.BX && 6) != 0 {
			if (InputDeltaX != 0) || (InputDeltaY != 0) {
				InputMouseButtonX = InputDeltaX
				InputMouseButtonY = InputDeltaY
			} else {
				InputDeltaX = InputMouseButtonX
				InputDeltaY = InputMouseButtonY
			}
		} else {
			InputMouseButtonX = 0
			InputMouseButtonY = 0
		}
	}

	if (InputDeltaX != 0) || (InputDeltaY != 0) {
		InputLastDeltaX = InputDeltaX
		InputLastDeltaY = InputDeltaY
	}
}

func InputInitMouse() (InputInitMouse bool) {
	var (
		regs Registers
	)
	regs.AX = 0
	Intr(0x33, regs)
	InputInitMouse = (regs.AX == 0)
	InputInitMouse = true
	return
}

func InputInitDevices() {
	InputJoystickEnabled = InputInitJoystick
	InputMouseEnabled = InputInitMouse
}

func InputConfigure() (InputConfigure bool) {
	var (
		charTyped byte
	)
	charTyped = ' '
	if InputJoystickEnabled || InputMouseEnabled {
		Writeln()
		Write("  Game controller:  K)eyboard")
		if InputJoystickEnabled {
			Write(",  J)oystick")
		}
		if InputMouseEnabled {
			Write(",  M)ouse")
		}
		Write("?  ")
		for {
			for {
				if KeyPressed {
					break
				}
			}
			charTyped = UpCase(ReadKey)
			if (charTyped == 'K') || (InputJoystickEnabled && (charTyped == 'J')) || (InputMouseEnabled && (charTyped == 'M')) || (charTyped == '\x1b') {
				break
			}
		}
		Writeln()
		InputJoystickEnabled = false
		InputMouseEnabled = false
		switch charTyped {
		case 'J':
			InputJoystickEnabled = true
			InputCalibrateJoystick()
		case 'M':
			InputMouseEnabled = true
		}
		Writeln()
	}
	InputConfigure = charTyped != '\x1b'
	return
}

func InputReadWaitKey() {
	for {
		InputUpdate()
		if InputKeyPressed != '\x00' {
			break
		}
	}
}

func init() {
	InputLastDeltaX = 0
	InputLastDeltaY = 0
	InputDeltaX = 0
	InputDeltaY = 0
	InputShiftPressed = false
	InputShiftAccepted = false
	InputMouseX = 0
	InputMouseY = 0
	InputMouseActivationX = 60
	InputMouseActivationY = 60
	InputMouseButtonX = 0
	InputMouseButtonY = 0
	InputKeyBuffer = ""
}
