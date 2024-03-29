package lexer

import (
	"testing"

	"sudocoding.xyz/interpreter_in_go/src/token"
)

func TestBasicNextToken(t *testing.T) {
	input := "=+(){},;-!*/><"

	expectedTokens := []token.Token{
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.MINUS, Literal: "-"},
		{Type: token.BANG, Literal: "!"},
		{Type: token.ASTERISK, Literal: "*"},
		{Type: token.SLASH, Literal: "/"},
		{Type: token.GT, Literal: ">"},
		{Type: token.LT, Literal: "<"},
		{Type: token.EOF, Literal: ""},
	}

	lexer := New(input)

	for i, expectedToken := range expectedTokens {
		tok := lexer.NextToken()

		if tok.Type != expectedToken.Type {
			t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q", i, expectedToken.Type, tok.Type)
		}

		if tok.Literal != expectedToken.Literal {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q", i, expectedToken.Literal, tok.Literal)
		}
	}
}

func TestNextToken(t *testing.T) {
	input := `let five = 5;
  let ten=10;

  let add = fn(x, y){
    x + y;
  }

  let result = add(five, ten);

  if (5 > 10) {
    return true;
  } else {
    return false;
  }

  5 == 10;
  5 != 10;
  5 >= 10;
  5 <= 10;

  "abcd";
  ""
  "asdf;
  `

	expectedTokens := []token.Token{
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "five"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "ten"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "add"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.FUNCTION, Literal: "fn"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENT, Literal: "x"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.IDENT, Literal: "x"},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "result"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.IDENT, Literal: "add"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENT, Literal: "five"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENT, Literal: "ten"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.IF, Literal: "if"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.INT, Literal: "5"},
		{Type: token.GT, Literal: ">"},
		{Type: token.INT, Literal: "10"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RETURN, Literal: "return"},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.ELSE, Literal: "else"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RETURN, Literal: "return"},
		{Type: token.FALSE, Literal: "false"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.INT, Literal: "5"},
		{Type: token.EQ, Literal: "=="},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.INT, Literal: "5"},
		{Type: token.NOT_EQ, Literal: "!="},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.INT, Literal: "5"},
		{Type: token.GTE, Literal: ">="},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.INT, Literal: "5"},
		{Type: token.LTE, Literal: "<="},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.STR, Literal: "abcd"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.STR, Literal: ""},
		{Type: token.ILLEGAL, Literal: "\x00"},
		{Type: token.EOF, Literal: ""},
	}

	lexer := New(input)

	for i, expectedToken := range expectedTokens {
		tok := lexer.NextToken()

		if tok.Type != expectedToken.Type {
			t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q", i, expectedToken.Type, tok.Type)
		}

		if tok.Literal != expectedToken.Literal {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q", i, expectedToken.Literal, tok.Literal)
		}
	}
}
