// Turbo Pascal lexer
//
// The lexer turns a string of Pascal source code into a stream of
// tokens for parsing.
//
// To tokenize some source, create a new lexer with NewLexer(src) and
// then call Scan() until the token type is EOF or ILLEGAL.

package main

import (
	"fmt"
)

// Lexer tokenizes a byte string of source code. Use NewLexer to
// actually create a lexer, and Scan() to get tokens.
type Lexer struct {
	src     []byte
	offset  int
	ch      byte
	pos     Position
	nextPos Position
}

// Position stores the source line and column where a token starts.
type Position struct {
	// Line number of the token (starts at 1).
	Line int
	// Column on the line (starts at 1). Note that this is the byte
	// offset into the line, not rune offset.
	Column int
}

// NewLexer creates a new lexer that will tokenize the given source
// code.
func NewLexer(src []byte) *Lexer {
	l := &Lexer{src: src}
	l.nextPos.Line = 1
	l.nextPos.Column = 1
	l.next()
	return l
}

// Scan scans the next token and returns its position (line/column),
// token value (one of the uppercased token constants), and the
// string value of the token. For most tokens, the token value is
// empty. For IDENT, NUMBER, and STRING tokens, it's the token's
// value. For an ILLEGAL token, it's the error message.
func (l *Lexer) Scan() (Position, Token, string) {
	// Skip whitespace and comments
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' || l.ch == '{' {
		if l.ch == '{' {
			l.next()
			for l.ch != '}' && l.ch != 0 {
				l.next()
			}
			l.next()
		}
		l.next()
	}
	if l.ch == 0 {
		// l.next() reached end of input
		return l.pos, EOF, ""
	}

	pos := l.pos
	tok := ILLEGAL
	val := ""

	ch := l.ch
	l.next()

	// Names: keywords and functions
	if isNameStart(ch) {
		start := l.offset - 2
		for isNameStart(l.ch) || (l.ch >= '0' && l.ch <= '9') || l.ch == '_' {
			l.next()
		}
		name := string(l.src[start : l.offset-1])
		tok := KeywordToken(name)
		if tok == ILLEGAL {
			tok = IDENT
			val = name
		}
		return pos, tok, val
	}

	switch ch {
	case '@':
		tok = AT
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		start := l.offset - 2
		for l.ch >= '0' && l.ch <= '9' {
			l.next()
		}
		if l.ch == '.' {
			l.next()
		}
		for l.ch >= '0' && l.ch <= '9' {
			l.next()
		}
		if l.ch == 'e' || l.ch == 'E' {
			l.next()
			gotSign := false
			if l.ch == '+' || l.ch == '-' {
				gotSign = true
				l.next()
			}
			gotDigit := false
			for l.ch >= '0' && l.ch <= '9' {
				l.next()
				gotDigit = true
			}
			// "1e" is allowed, but not "1e+"
			if gotSign && !gotDigit {
				return l.pos, ILLEGAL, "expected digits"
			}
		}
		tok = NUM
		val = string(l.src[start : l.offset-1])
	case '$':
		l.next()
		start := l.offset - 2
		for isHexDigit(l.ch) {
			l.next()
		}
		tok = HEX
		val = string(l.src[start : l.offset-1])
	case ':':
		tok = l.choice('=', COLON, ASSIGN)
	case '=':
		tok = EQUALS
	case '<':
		switch l.ch {
		case '=':
			l.next()
			tok = LTE
		case '>':
			l.next()
			tok = NOT_EQUALS
		default:
			tok = LESS
		}
	case '>':
		tok = l.choice('=', GREATER, GTE)
	case '\'', '#':
		chars := make([]byte, 0, 32) // most won't require heap allocation
		for {
			if ch == '\'' {
				for {
					c := l.ch
					if c == 0 {
						return l.pos, ILLEGAL, "didn't find end quote in string"
					}
					if c == '\r' || c == '\n' {
						return l.pos, ILLEGAL, "can't have newline in string"
					}
					if c == '\'' {
						l.next()
						if l.ch != '\'' {
							break
						}
					}
					chars = append(chars, c)
					l.next()
				}
			} else { // '#'
				num := 0
				for l.ch >= '0' && l.ch <= '9' {
					num = num*10 + int(l.ch-'0')
					if num > 255 {
						return l.pos, ILLEGAL, "#char greater than 255"
					}
					l.next()
				}
				chars = append(chars, byte(num))
			}
			if l.ch != '\'' && l.ch != '#' {
				break
			}
			ch = l.ch
			l.next()
		}
		tok = STR
		val = string(chars)
	case '(':
		tok = LPAREN
	case ')':
		tok = RPAREN
	case ',':
		tok = COMMA
	case ';':
		tok = SEMICOLON
	case '+':
		tok = PLUS
	case '-':
		tok = MINUS
	case '*':
		tok = STAR
	case '/':
		tok = SLASH
	case '[':
		tok = LBRACKET
	case ']':
		tok = RBRACKET
	case '^':
		tok = POINTER
	case '.':
		tok = l.choice('.', DOT, DOT_DOT)
	default:
		tok = ILLEGAL
		val = fmt.Sprintf("unexpected char %q", ch)
	}
	return pos, tok, val
}

// Load the next character into l.ch (or 0 on end of input) and update
// line and column position.
func (l *Lexer) next() {
	l.pos = l.nextPos
	if l.offset >= len(l.src) {
		// For last character, move offset 1 past the end as it
		// simplifies offset calculations in IDENT and NUMBER
		if l.ch != 0 {
			l.ch = 0
			l.offset++
		}
		return
	}
	ch := l.src[l.offset]
	if ch == '\n' {
		l.nextPos.Line++
		l.nextPos.Column = 1
	} else {
		l.nextPos.Column++
	}
	l.ch = ch
	l.offset++
}

func isNameStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (l *Lexer) choice(ch byte, one, two Token) Token {
	if l.ch == ch {
		l.next()
		return two
	}
	return one
}
