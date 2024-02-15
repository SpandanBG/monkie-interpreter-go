package lexer

import "sudocoding.xyz/interpreter_in_go/src/token"

// Lexer - Contains the lexer struct to parse input to tokens
type Lexer struct {
	input        string // input will be read as ASCII (byte) instead of Unicode (rune)
	position     int    // current position in input (points to current char)
	readPosition int    // current reading position in input (after current char)
	ch           byte   // current char under examination
}

// New - creates a new lexer for for input and starts off by reading the the
// first character
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// NextToken - reads the token in the current postion and returns it
// and moved the position to the beginning of the next token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string([]byte{ch, l.ch})

			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch)}
		}
	case '+':
		tok = token.Token{Type: token.PLUS, Literal: string(l.ch)}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch)}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch)}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch)}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch)}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch)}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch)}
	case '-':
		tok = token.Token{Type: token.MINUS, Literal: string(l.ch)}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch)}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch)}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string([]byte{ch, l.ch})

			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch)}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string([]byte{ch, l.ch})

			tok = token.Token{Type: token.GTE, Literal: literal}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch)}
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string([]byte{ch, l.ch})

			tok = token.Token{Type: token.LTE, Literal: literal}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch)}
		}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)

			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()

			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

// peekChar - returns the character at the `readPosition`. 0 if at the end of
// the input. Doesn't modify the lexer object
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

// readChar - updates `ch` to the character in `readPosition`.
// sets `ch` to `0` if `readPosition` is out of bound in input
// updates `positon` to `readPosition` and increments `readPosition` by `1`.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // `NUL` in ASCII
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

// readIdentifier - keeps reading char till non letter (based on `isLetter` func)
// is reached and returns the sequence of chars read.
func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// readNumber - keeps reading char till non number (base on `isDigit` func) is
// reached and returns the sequence of numbers read.
func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// skipWhitespace - skips all whitespace characters (\s \n \t \r) and moves the
// `position` to the next non-whitespace character.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// isLetter - returns true if a given char follows the following regex
// `[a-zA-Z]|_`
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit - returns truw if a given char follows the following regex
// `[0-9]`
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
