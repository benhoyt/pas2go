package main

import (
	"io"
	"math/rand"
	"os"
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

// File routines

type File struct {
	name []byte
	file *os.File
}

var ioResult int16

func IOResult() int16 {
	return ioResult
}

func setIOResult(err error) {
	ioResult = 0
	if err != nil {
		ioResult = 2 // "File not found" (good enough for our purposes)
	}
}

func Assign(f *File, name []byte) {
	f.name = name
}

func Reset(f *File) {
	file, err := os.Open(f.name)
	f.file = file
	setIOResult(err)
}

func Rewrite(f *File) {
	file, err := os.Create(f.name)
	f.file = file
	setIOResult(err)
}

func Read(f *File, data interface{}) {
	err := binary.Read(f.file, binary.LittleEndian, data)
	setIOResult(err)
}

func Write(f *File, data interface{}) {
	err := binary.Write(f.file, binary.LittleEndian, data)
	setIOResult(err)
}

func Close(f *File) {
	err := f.Close()
	setIOResult(err)
}

func Erase(f *File) {
	err := os.Remove()
	setIOResult(err)
}

func Seek(f *File, offset int16) {
	_, err := f.file(int64(offset), io.SeekSet)
	setIOResult(err)
}
