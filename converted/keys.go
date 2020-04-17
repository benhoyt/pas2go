package main // unit: Keys

var (
	KeysRightShiftHeld bool
	KeysLeftShiftHeld  bool
	KeysShiftHeld      bool
	KeysCtrlHeld       bool
	KeysAltHeld        bool
	KeysNumLockHeld    bool
)

// implementation uses: Dos

func KeysUpdateModifiers() {
	var regs Registers
	regs.AH = 0x02
	Intr(0x16, regs)
	KeysRightShiftHeld = regs.AL%2 == 1
	KeysLeftShiftHeld = regs.AL/2%2 == 1
	KeysCtrlHeld = regs.AL/4%2 == 1
	KeysAltHeld = regs.AL/8%2 == 1
	KeysNumLockHeld = regs.AL/32%2 == 1
	KeysShiftHeld = KeysRightShiftHeld || KeysLeftShiftHeld
}
