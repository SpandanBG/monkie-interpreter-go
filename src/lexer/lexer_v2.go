package lexer

import (
	"bufio"
	"bytes"
	"io"
	"unicode"

	"sudocoding.xyz/interpreter_in_go/src/token"
)

type Lexer_V2 struct {
	input *bufio.Reader
	ch    rune
}

func New_V2(reader io.Reader) *Lexer_V2 {
	l_v2 := &Lexer_V2{input: bufio.NewReader(reader)}
	l_v2.readChar_v2()
	return l_v2
}

func (l_v2 *Lexer_V2) NextToken_V2() token.Token {
	var tok token.Token

	l_v2.skipWhitespace()

	switch l_v2.ch {
	case rune('='):
		if l_v2.peekChar_v2() == rune('=') {
			ch := l_v2.ch
			l_v2.readChar_v2()
			literal := string([]rune{ch, l_v2.ch})

			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l_v2.ch)}
		}
	case rune('+'):
		tok = token.Token{Type: token.PLUS, Literal: string(l_v2.ch)}
	case rune('('):
		tok = token.Token{Type: token.LPAREN, Literal: string(l_v2.ch)}
	case rune(')'):
		tok = token.Token{Type: token.RPAREN, Literal: string(l_v2.ch)}
	case rune('{'):
		tok = token.Token{Type: token.LBRACE, Literal: string(l_v2.ch)}
	case rune('}'):
		tok = token.Token{Type: token.RBRACE, Literal: string(l_v2.ch)}
	case rune('['):
		tok = token.Token{Type: token.LBRACKET, Literal: string(l_v2.ch)}
	case rune(']'):
		tok = token.Token{Type: token.RBRACKET, Literal: string(l_v2.ch)}
	case rune(','):
		tok = token.Token{Type: token.COMMA, Literal: string(l_v2.ch)}
	case rune(':'):
		tok = token.Token{Type: token.COLON, Literal: string(l_v2.ch)}
	case rune(';'):
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l_v2.ch)}
	case rune('-'):
		tok = token.Token{Type: token.MINUS, Literal: string(l_v2.ch)}
	case rune('*'):
		tok = token.Token{Type: token.ASTERISK, Literal: string(l_v2.ch)}
	case rune('/'):
		tok = token.Token{Type: token.SLASH, Literal: string(l_v2.ch)}
	case rune('!'):
		if l_v2.peekChar_v2() == rune('=') {
			ch := l_v2.ch
			l_v2.readChar_v2()
			literal := string([]rune{ch, l_v2.ch})

			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l_v2.ch)}
		}
	case rune('>'):
		if l_v2.peekChar_v2() == rune('=') {
			ch := l_v2.ch
			l_v2.readChar_v2()
			literal := string([]rune{ch, l_v2.ch})

			tok = token.Token{Type: token.GTE, Literal: literal}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l_v2.ch)}
		}
	case rune('<'):
		if l_v2.peekChar_v2() == rune('=') {
			ch := l_v2.ch
			l_v2.readChar_v2()
			literal := string([]rune{ch, l_v2.ch})

			tok = token.Token{Type: token.LTE, Literal: literal}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l_v2.ch)}
		}
	case rune('"'):
		if str, ok := l_v2.readStr_v2(); ok {
			tok = token.Token{Type: token.STR, Literal: str}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l_v2.ch)}
		}
	case rune(0):
		tok = token.Token{Type: token.EOF, Literal: string(l_v2.ch)}
	default:
		if isLetter_v2(l_v2.ch) {
			tok.Literal = l_v2.readGroup_v2(isLetter_v2)
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit_v2(l_v2.ch) {
			tok.Literal = l_v2.readGroup_v2(isDigit_v2)
			tok.Type = token.INT
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l_v2.ch)}
		}
	}

	l_v2.readChar_v2()
	return tok
}

func (l_v2 *Lexer_V2) peekChar_v2() rune {
	ch, _, err := l_v2.input.ReadRune()
	if err != nil {
		ch = rune(0)
	} else {
		l_v2.input.UnreadRune()
	}

	return ch
}

func (l_v2 *Lexer_V2) readChar_v2() {
	ch, _, err := l_v2.input.ReadRune()
	if err != nil {
		ch = rune(0)
	}

	l_v2.ch = ch
}

func (l_v2 *Lexer_V2) readStr_v2() (string, bool) {
	var out bytes.Buffer
	l_v2.readChar_v2()

	for l_v2.ch != rune('"') {
		if l_v2.ch == rune(0) {
			return "", false
		}
		out.WriteRune(l_v2.ch)
		l_v2.readChar_v2()
	}

	return out.String(), true
}

func (l_v2 *Lexer_V2) readGroup_v2(filterFn func(ch rune) bool) string {
	var idBuffer bytes.Buffer

	for filterFn(l_v2.ch) {
		idBuffer.WriteRune(l_v2.ch)
		l_v2.readChar_v2()
	}

	return idBuffer.String()
}

func (l_v2 *Lexer_V2) skipWhitespace() {
	for unicode.IsSpace(l_v2.ch) {
		l_v2.readChar_v2()
	}
}

func isLetter_v2(ch rune) bool {
	return unicode.IsLetter(ch) || ch == rune('_')
}

func isDigit_v2(ch rune) bool {
	return unicode.IsDigit(ch)
}
