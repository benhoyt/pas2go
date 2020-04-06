package main

import (
	"math/rand"
	"time"
)

type TVideoLine string

func VideoWriteText(x, y, color byte, text TVideoLine) {
	// TODO
}

func Ord(c byte) byte {
	return c
}

func Chr(i byte) byte {
	return i
}

func Length(b []byte) int16 {
	return len(b)
}

func Delay(milliseconds int16) {
	time.Sleep(milliseconds * time.Millisecond)
}

func Random(end int16) int16 {
	return int16(rand.Intn(end))
}

func Sqr(n int16) int16 {
	return n * n
}

func Str(n int16, s []byte) {
	// TODO
}

func StrWidth(n, width int16, s []byte) {
	// TODO
}
