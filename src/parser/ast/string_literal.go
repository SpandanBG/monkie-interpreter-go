package ast

import (
	"fmt"

	"sudocoding.xyz/interpreter_in_go/src/token"
)

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", sl.Token.Literal)
}
