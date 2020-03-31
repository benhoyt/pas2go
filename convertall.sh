#!/usr/bin/env bash

go build
./pas2go convert orig/ELEMENTS.PAS > converted/elements.go
./pas2go convert orig/GAME.PAS > converted/game.go
./pas2go convert orig/GAMEVARS.PAS > converted/gamevars.go
./pas2go convert orig/INPUT.PAS > converted/input.go
./pas2go convert orig/KEYS.PAS > converted/keys.go
./pas2go convert orig/OOP.PAS > converted/oop.go
./pas2go convert orig/EDITOR.PAS > converted/editor.go
./pas2go convert orig/SOUNDS.PAS > converted/sounds.go
./pas2go convert orig/TXTWIND.PAS > converted/txtwind.go
./pas2go convert orig/ZZT.PAS > converted/zzt.go
go fmt converted/*.go
