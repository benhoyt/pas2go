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

    x := (1 + 2) * (3 + 4);

    case ToUpper(3) of
        65: x := 1
    end;

    x := foo(1, 3.4, 0.01);
end.
