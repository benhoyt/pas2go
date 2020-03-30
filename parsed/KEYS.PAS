unit Keys;

interface
    
    var
        KeysRightShiftHeld: boolean;
        KeysLeftShiftHeld: boolean;
        KeysShiftHeld: boolean;
        KeysCtrlHeld: boolean;
        KeysAltHeld: boolean;
        KeysNumLockHeld: boolean;
    procedure KeysUpdateModifiers;

implementation
uses Dos;

procedure KeysUpdateModifiers;
    var
        regs: Registers;
    begin
        regs.AH := $02;
        Intr($16, regs);
        KeysRightShiftHeld := (regs.AL mod 2) = 1;
        KeysLeftShiftHeld := ((regs.AL div 2) mod 2) = 1;
        KeysCtrlHeld := ((regs.AL div 4) mod 2) = 1;
        KeysAltHeld := ((regs.AL div 8) mod 2) = 1;
        KeysNumLockHeld := ((regs.AL div 32) mod 2) = 1;
        KeysShiftHeld := KeysRightShiftHeld or KeysLeftShiftHeld;
    end;

begin
    
end.
