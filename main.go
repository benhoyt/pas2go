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
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintf(os.Stderr, "usage: pas2go [lex | parse] [file.pas]\n")
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
		parse(src)
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

func parse(src []byte) {
	file, err := Parse(src)
	if err != nil {
		errMsg := fmt.Sprintf("%s", err)
		if err, ok := err.(*ParseError); ok {
			showSourceLine(src, err.Position, len(errMsg))
		}
		fmt.Fprintf(os.Stderr, "%s\n", errMsg)
		os.Exit(1)
	}
	fmt.Print(file)
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
