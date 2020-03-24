{ comment
 blah
}

program TestProg;

begin
    a := 1;
    b := 1234;
    if x + 12*3 = 34 then begin
        ClrScr;
    end else if 0 then begin
        DoNothing;
    end else begin FooBar(1, 'foo', -314); end;

    for x := 10 DOWNTO 0 do begin
        Writeln(1);
    end;

    @x := 1234;

    Thingy:
    xyz := 42;

    goto Thingy;

    begin
        inside := 42;
        ClrScr;
    end;

    while 42 do begin Clrscr; end;

    repeat
        this;
        that;
        otherThing;
    until 1;

    with FooBar do begin ClrScr; end;
end.
