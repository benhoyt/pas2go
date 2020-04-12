#!/usr/bin/env bash

go build
./pas2go parse orig/EDITOR.PAS > parsed/EDITOR.PAS
./pas2go parse orig/ELEMENTS.PAS > parsed/ELEMENTS.PAS
./pas2go parse orig/GAME.PAS > parsed/GAME.PAS
./pas2go parse orig/GAMEVARS.PAS > parsed/GAMEVARS.PAS
./pas2go parse orig/INPUT.PAS > parsed/INPUT.PAS
./pas2go parse orig/KEYS.PAS > parsed/KEYS.PAS
./pas2go parse orig/OOP.PAS > parsed/OOP.PAS
./pas2go parse orig/SOUNDS.PAS > parsed/SOUNDS.PAS
./pas2go parse orig/TXTWIND.PAS > parsed/TXTWIND.PAS
./pas2go parse orig/ZZT.PAS > parsed/ZZT.PAS
