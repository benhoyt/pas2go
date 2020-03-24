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
	src, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %q\n", os.Args[1])
		os.Exit(1)
	}

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

	program, err := ParseProgram(src)
	if err != nil {
		errMsg := fmt.Sprintf("%s", err)
		if err, ok := err.(*ParseError); ok {
			showSourceLine(src, err.Position, len(errMsg))
		}
		fmt.Fprintf(os.Stderr, "%s\n", errMsg)
		os.Exit(1)
	}
	fmt.Printf("----------------------\n")
	fmt.Print(program)
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
