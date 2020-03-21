// Turbo Pascal to Go transpiler

package main

import (
	"fmt"
	"io/ioutil"
	"os"
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
}
