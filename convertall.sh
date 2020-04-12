#!/usr/bin/env bash

go build
./pas2go convert orig/EDITOR.PAS orig/GAMEVARS.PAS orig/TXTWIND.PAS orig/SOUNDS.PAS orig/INPUT.PAS orig/ELEMENTS.PAS orig/OOP.PAS orig/GAME.PAS > converted/editor.go
./pas2go convert orig/ELEMENTS.PAS orig/GAMEVARS.PAS orig/SOUNDS.PAS orig/INPUT.PAS orig/TXTWIND.PAS orig/EDITOR.PAS orig/OOP.PAS orig/GAME.PAS > converted/elements.go
./pas2go convert orig/GAME.PAS orig/GAMEVARS.PAS orig/TXTWIND.PAS orig/SOUNDS.PAS orig/INPUT.PAS orig/ELEMENTS.PAS orig/EDITOR.PAS orig/OOP.PAS > converted/game.go
./pas2go convert orig/GAMEVARS.PAS > converted/gamevars.go
./pas2go convert orig/INPUT.PAS orig/KEYS.PAS orig/SOUNDS.PAS > converted/input.go
./pas2go convert orig/KEYS.PAS > converted/keys.go
./pas2go convert orig/OOP.PAS orig/GAMEVARS.PAS orig/SOUNDS.PAS orig/TXTWIND.PAS orig/GAME.PAS orig/ELEMENTS.PAS > converted/oop.go
./pas2go convert orig/SOUNDS.PAS > converted/sounds.go
./pas2go convert orig/TXTWIND.PAS orig/INPUT.PAS > converted/txtwind.go
./pas2go convert orig/ZZT.PAS orig/KEYS.PAS orig/SOUNDS.PAS orig/INPUT.PAS orig/TXTWIND.PAS orig/GAMEVARS.PAS orig/ELEMENTS.PAS orig/EDITOR.PAS orig/OOP.PAS orig/GAME.PAS > converted/zzt.go
go fmt converted/*.go
