package lexer

import (
	"bufio"
	"io"

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

  switch l_v2.ch {
  case rune('='):
    tok = token.Token{Type: token.ASSIGN, Literal: string(l_v2.ch)}
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
	case rune(','):
		tok = token.Token{Type: token.COMMA, Literal: string(l_v2.ch)}
	case rune(';'):
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l_v2.ch)}
  case rune(0):
    tok = token.Token{Type: token.EOF, Literal: string(l_v2.ch)}
  }

  l_v2.readChar_v2()
  return tok
}

func (l_v2 *Lexer_V2) readChar_v2() {
	ch, _, err := l_v2.input.ReadRune()
	if err != nil {
		ch = rune(0)
	}

	l_v2.ch = ch
}
