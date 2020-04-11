package main

import (
	"encoding/binary"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
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

func Val(s string) (n int16, code int16) {
	nInt, err := strconv.Atoi(s)
	if err != nil {
		return int16(nInt), 1 // TODO: code should be 1-based position of first non-number char
	}
	return nInt, 0
}

func Str(n int16, s string) {
	// TODO
}

func StrWidth(n, width int16, s string) {
	// TODO
}

// NOTE: in Turbo Pascal Delete() is a procedure that modifies the string in-place
func Delete(s string, index, count int16) string {
	return s[:index-1] + s[index-1+count:]
}

// Misc functions

var Time int16 // TODO

func VideoWriteText(x, y, color byte, text string) {
	// TODO
}

var VideoMonochrome bool

func Delay(milliseconds int16) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

func Random(end int16) int16 {
	return int16(rand.Intn(int(end)))
}

func Sqr(n int16) int16 {
	return n * n
}

func Sound(x int16) {
	// TODO
}

func NoSound() {
	// TODO
}

func Abs(n int16) int16 {
	if n < 0 {
		return -n
	}
	return n
}

func Ln(x float64) float64 {
	return math.Log(x)
}

func Exp(x float64) float64 {
	return math.Exp(x)
}

func Trunc(x float64) int16 {
	return int16(x)
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

func Reset(f *File, _ ...interface{}) {
	file, err := os.Open(f.name)
	f.file = file
	setIOResult(err)
}

func Rewrite(f *File, _ ...interface{}) {
	file, err := os.Create(f.name)
	f.file = file
	setIOResult(err)
}

func Read(f *File, data interface{}) {
	err := binary.Read(f.file, binary.LittleEndian, data)
	setIOResult(err)
}

func BlockRead(f *File, buf interface{}, count int16) {
	// TODO
}

func BlockWrite(f *File, buf interface{}, count int16) {
	// TODO
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

// Memory functions

var MemAvail = 32767

func New(p interface{}) {
	// TODO
}

func FreeMem(p interface{}, size int16) {
	// TODO
}

func GetMem(p interface{}, size int16) {
	// TODO
}

func Move(src, dest interface{}, size int16) {
	// TODO
}

func SizeOf(val interface{}) int16 {
	// TODO
	return 0
}

func FillChar(dest interface{}, count int16, c byte) {
	// TODO
}
