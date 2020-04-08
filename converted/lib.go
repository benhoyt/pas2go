package main

import (
	"encoding/binary"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

// String functions

func Ord(c byte) byte {
	return c
}

func Chr(i byte) byte {
	return i
}

func Length(s string) int16 {
	return int16(len(s))
}

func UpCase(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - ('a' - 'A')
	}
	return b
}

func Copy(s string, index, count int16) string {
	return s[index-1 : index-1+count]
}

func Pos(b byte, s string) int16 {
	return int16(strings.IndexByte(s, b) + 1)
}

// NOTE: in Turbo Pascal Delete() is a procedure that modifies the string in-place
func Delete(s string, index, count int16) string {
	return s[:index-1] + s[index-1+count:]
}

// Misc functions

type TVideoLine string

func VideoWriteText(x, y, color byte, text TVideoLine) {
	// TODO
}

func Delay(milliseconds int16) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

func Random(end int16) int16 {
	return int16(rand.Intn(int(end)))
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

// File functions

type File struct {
	name string
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

func Assign(f *File, name string) {
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
	err := f.file.Close()
	setIOResult(err)
}

func Erase(f *File) {
	err := os.Remove(f.name)
	setIOResult(err)
}

func Seek(f *File, offset int16) {
	_, err := f.file.Seek(int64(offset), io.SeekStart)
	setIOResult(err)
}
