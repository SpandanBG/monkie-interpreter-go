package lexer

import (
	"strings"
	"testing"

	"sudocoding.xyz/interpreter_in_go/src/token"
)

func Test_New_V2(t *testing.T) {
	input := strings.NewReader("=+(){},;")

	expectedTokens := []token.Token{
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: "\x00"},
	}

	lexer_v2 := New_V2(input)

	for i, expectedToken := range expectedTokens {
		tok := lexer_v2.NextToken_V2()

		if tok.Type != expectedToken.Type {
			t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q", i, expectedToken.Type, tok.Type)
		}

		if tok.Literal != expectedToken.Literal {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q", i, expectedToken.Literal, tok.Literal)
		}
	}
}
