// Turbo Pascal to Go transpiler

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	panic("THIS IS A FUCKING CRAP, YOU ARE A FUCKING LAMER")
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: pas2go [lex | parse | convert] [file.pas] [unit1.pas ...]\n")
		os.Exit(1)
	}

	command := os.Args[1]

	var src []byte
	if len(os.Args) > 2 {
		path := os.Args[2]
		var err error
		src, err = ioutil.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
			os.Exit(1)
		}
	} else {
		var err error
		src, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %v", err)
			os.Exit(1)
		}
	}

	switch command {
	case "lex":
		lex(src)
	case "parse":
		file := parse(src)
		fmt.Print(file)
	case "convert":
		file := parse(src)

		units := []*Unit{}
		for _, path := range os.Args[3:] {
			unitSrc, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
				os.Exit(1)
			}
			unitFile := parse(unitSrc)
			unit, ok := unitFile.(*Unit)
			if !ok {
				continue
			}
			units = append(units, unit)
		}

		Convert(file, units, os.Stdout)
	default:
		fmt.Fprintf(os.Stderr, "command must be 'lex' or 'parse'")
		os.Exit(1)
	}
}

func lex(src []byte) {
	lexer := NewLexer(src)
	for {
		pos, tok, val := lexer.Scan()
		if tok == EOF {
			break
		}
		fmt.Printf("%d:%d %s %q\n", pos.Line, pos.Column, tok, val)
		if tok == ILLEGAL {
			break
		}
	}
}

func parse(src []byte) File {
	file, err := Parse(src)
	if err != nil {
		errMsg := fmt.Sprintf("%s", err)
		if err, ok := err.(*ParseError); ok {
			showSourceLine(src, err.Position, len(errMsg))
		}
		fmt.Fprintf(os.Stderr, "%s\n", errMsg)
		os.Exit(1)
	}
	return file
}

func showSourceLine(src []byte, pos Position, dividerLen int) {
	divider := strings.Repeat("-", dividerLen)
	if divider != "" {
		fmt.Fprintln(os.Stderr, divider)
	}
	lines := bytes.Split(src, []byte{'\n'})
	srcLine := string(lines[pos.Line-1])
	numTabs := strings.Count(srcLine[:pos.Column-1], "\t")
	runeColumn := utf8.RuneCountInString(srcLine[:pos.Column-1])
	fmt.Fprintln(os.Stderr, strings.Replace(srcLine, "\t", "    ", -1))
	fmt.Fprintln(os.Stderr, strings.Repeat(" ", runeColumn)+strings.Repeat("   ", numTabs)+"^")
	if divider != "" {
		fmt.Fprintln(os.Stderr, divider)
	}
}
