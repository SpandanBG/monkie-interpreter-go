package ast

import (
	"bytes"

	"sudocoding.xyz/interpreter_in_go/src/token"
)

type Assignment struct {
	Token      token.Token
	Identifier *Identifier
	Value      Expression
}

func (a *Assignment) statementNode() {}

func (a *Assignment) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Assignment) String() string {
	var out bytes.Buffer

	out.WriteString(a.Identifier.String())
	out.WriteString(" = ")
	out.WriteString(a.Value.String())

	return out.String()
}
